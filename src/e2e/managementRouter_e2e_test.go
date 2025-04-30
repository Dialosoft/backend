//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/app/database"
	"github.com/Dialosoft/src/domain/models"
)

type ManagementRouterTestSuite struct {
	BaseE2eTestSuite
	server  *TestServer
	helpers TestHelpers
}

// SetupSuite is called before any Management Router tests in the suite are run.
func (suite *ManagementRouterTestSuite) SetupSuite() {
	suite.BaseE2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
	suite.helpers = TestHelpers{
		Suite:  &suite.Suite,
		Server: suite.server,
	}
}

// TearDownSuite is called after all Management Router tests in the suite have been run.
func (suite *ManagementRouterTestSuite) TearDownSuite() {
	// Shutdown the server after a set of tests
	err := suite.server.Shutdown()
	suite.NoError(err)

	suite.BaseE2eTestSuite.TearDownSuite()
}

func (suite *ManagementRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)

	// Verify if the roles were created
	var count int64
	err = suite.db.Model(&models.RoleEntity{}).Count(&count).Error
	suite.NoError(err)
	suite.Greater(count, int64(0), "No roles were created")
}

func (suite *ManagementRouterTestSuite) TearDownTest() {
	// truncate Tables maintaining the structure and restrictions
	tables := []string{
		"comment_votes",
		"posts_likes",
		"comments",
		"posts",
		"forums",
		"categories",
		"tokens",
		"users",
		"role_permissions",
		"roles",
	}

	for _, table := range tables {
		err := suite.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error
		suite.NoError(err)
	}

	// Clean Redis at the end of each test
	err := suite.rdClient.FlushAll(suite.testCtx).Err()
	suite.NoError(err)
}

// TestSuccessfulRoleChangeByAdmin tests the successful role change by an admin
func (suite *ManagementRouterTestSuite) TestSuccessfulRoleChangeByAdmin() {
	// Get admin credentials
	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

	// Now register a regular user to change their role
	regularUserName := "usertest"
	regularUserEmail := "user@example.com"
	regularUserPassword := "userPassword123"
	userID, _, userRefreshToken := suite.helpers.RegisterAndLoginUser(regularUserName, regularUserEmail, regularUserPassword)

	// Get moderator role ID
	moderatorRoleID := suite.helpers.GetRoleIDByType(suite.db, "moderator")

	changeRolePayload := request.ChangeUserRole{
		UserID: userID,
		RoleID: moderatorRoleID,
	}

	// Set custom options with admin tokens for authorization
	options := &RequestOptions{
		Token:        adminToken,
		RefreshToken: adminRefreshToken,
	}

	// Make change-user-role request as admin
	resp, err := suite.helpers.MakeRequestWithOptions(
		http.MethodPost,
		"/management/change-user-role",
		changeRolePayload,
		http.StatusOK,
		options,
	)
	suite.NoError(err)
	defer resp.Body.Close()

	// Verify role was changed in database
	var updatedUser models.UserEntity
	result := suite.db.Preload("Role").First(&updatedUser, "id = ?", userID)
	suite.NoError(result.Error)
	suite.Equal(moderatorRoleID, updatedUser.RoleID.String())
	suite.Equal("moderator", updatedUser.Role.RoleType)
	suite.True(updatedUser.Role.ModRole, "User should have moderator role privileges")

	// Verify user's refresh token was invalidated by checking directly in Redis
	isBlacklisted, err := suite.rdClient.Exists(suite.testCtx, fmt.Sprintf("blacklist:%s", userRefreshToken)).Result()
	suite.NoError(err)
	suite.GreaterOrEqual(isBlacklisted, int64(1), "Refresh token should be blacklisted in Redis")

	// Also verify the token fails when used for authentication
	refreshPayload := map[string]string{
		"refreshToken": userRefreshToken,
	}

	resp, err = suite.helpers.MakeRequest(http.MethodPost, "/auth/refresh-token", refreshPayload, http.StatusUnauthorized)
	suite.NoError(err)
	defer resp.Body.Close()
}

// TestChangeRoleWithoutAdminPermissions tests that a non-admin user cannot change roles
func (suite *ManagementRouterTestSuite) TestChangeRoleWithoutAdminPermissions() {
	// Create admin first to get moderator role ID (but we won't use the admin for requests)
	_, _, _ = suite.helpers.CreateAdminUserAndLogin(suite.db)
	moderatorRoleID := suite.helpers.GetRoleIDByType(suite.db, "moderator")

	// Register first regular user who will attempt to change roles
	anotherUserName := "usertest2"
	anotherUserEmail := "another@example.com"
	anotherUserPassword := "anotherPassword123"
	_, anotherUserToken, anotherUserRefreshToken := suite.helpers.RegisterAndLoginUser(anotherUserName, anotherUserEmail, anotherUserPassword)

	// Register second regular user who will be the target of the role change
	targetUserName := "targetuser"
	targetUserEmail := "target@example.com"
	targetUserPassword := "targetPassword123"
	targetUserID, _, _ := suite.helpers.RegisterAndLoginUser(targetUserName, targetUserEmail, targetUserPassword)

	// Attempt to change role as a regular user without privileges
	changeRolePayload := request.ChangeUserRole{
		UserID: targetUserID,
		RoleID: moderatorRoleID,
	}

	// Set options with regular user token (non-admin)
	options := &RequestOptions{
		Token:        anotherUserToken,
		RefreshToken: anotherUserRefreshToken,
	}

	// Make role change request as regular user (should fail)
	resp, err := suite.helpers.MakeRequestWithOptions(
		http.MethodPost,
		"/management/change-user-role",
		changeRolePayload,
		http.StatusForbidden, // Expect forbidden status
		options,
	)
	suite.NoError(err)
	defer resp.Body.Close()

	// Verify that the target user's role was not changed
	var targetUser models.UserEntity
	result := suite.db.Preload("Role").First(&targetUser, "id = ?", targetUserID)
	suite.NoError(result.Error)
	suite.Equal("user", targetUser.Role.RoleType, "Target user role should not have changed")
}

// TestInvalidUserIDForRoleChange tests that invalid user ID is properly rejected
func (suite *ManagementRouterTestSuite) TestInvalidUserIDForRoleChange() {
	// Get admin credentials
	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)
	moderatorRoleID := suite.helpers.GetRoleIDByType(suite.db, "moderator")

	// Set options with admin tokens
	options := &RequestOptions{
		Token:        adminToken,
		RefreshToken: adminRefreshToken,
	}

	// Invalid user ID test case
	invalidChangeRolePayload := request.ChangeUserRole{
		UserID: "invalid-uuid",
		RoleID: moderatorRoleID,
	}

	resp, err := suite.helpers.MakeRequestWithOptions(
		http.MethodPost,
		"/management/change-user-role",
		invalidChangeRolePayload,
		http.StatusBadRequest,
		options,
	)
	suite.NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	// Verify error response body
	var errorResponse response.StandardError
	err = json.Unmarshal(body, &errorResponse)
	suite.NoError(err)

	suite.Equal("ID provided is not a valid UUID type", errorResponse.ErrorMessage)
}

// TestMissingRoleIDForRoleChange tests that missing required role ID field is properly rejected
func (suite *ManagementRouterTestSuite) TestMissingRoleIDForRoleChange() {
	// Get admin credentials
	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

	// Register a user to use as valid ID in test
	regularUserName := "usertest"
	regularUserEmail := "user@example.com"
	regularUserPassword := "userPassword123"
	userID, _, _ := suite.helpers.RegisterAndLoginUser(regularUserName, regularUserEmail, regularUserPassword)

	// Set options with admin tokens
	options := &RequestOptions{
		Token:        adminToken,
		RefreshToken: adminRefreshToken,
	}

	// Missing required field (RoleID) test case
	missingFieldsPayload := request.ChangeUserRole{
		UserID: userID,
		// RoleID intentionally omitted
	}

	resp, err := suite.helpers.MakeRequestWithOptions(
		http.MethodPost,
		"/management/change-user-role",
		missingFieldsPayload,
		http.StatusBadRequest,
		options,
	)
	suite.NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	// Verify error response body
	var errorResponse response.StandardError
	err = json.Unmarshal(body, &errorResponse)
	suite.NoError(err)

	suite.Equal("One of the parameters or arguments is empty", errorResponse.ErrorMessage)
}

// TestNonExistentRoleIDForRoleChange tests that non-existent role ID is properly handled
func (suite *ManagementRouterTestSuite) TestNonExistentRoleIDForRoleChange() {
	// Get admin credentials
	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

	// Register a user to use as valid ID in test
	regularUserName := "usertest"
	regularUserEmail := "user@example.com"
	regularUserPassword := "userPassword123"
	userID, _, _ := suite.helpers.RegisterAndLoginUser(regularUserName, regularUserEmail, regularUserPassword)

	// Set options with admin tokens
	options := &RequestOptions{
		Token:        adminToken,
		RefreshToken: adminRefreshToken,
	}

	// Non-existent role ID test case
	nonExistentRolePayload := request.ChangeUserRole{
		UserID: userID,
		RoleID: "00000000-0000-0000-0000-000000000000", // Non-existent UUID
	}

	resp, err := suite.helpers.MakeRequestWithOptions(
		http.MethodPost,
		"/management/change-user-role",
		nonExistentRolePayload,
		http.StatusInternalServerError, // Expect internal error when role doesn't exist
		options,
	)
	suite.NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	// Verify error response body
	var errorResponse response.StandardError
	err = json.Unmarshal(body, &errorResponse)
	suite.NoError(err)

	suite.Equal("record not found", errorResponse.Details)
}

// Entry point for the Management Router test suite
func TestManagementRouterSuite(t *testing.T) {
	suite.Run(t, new(ManagementRouterTestSuite))
}
