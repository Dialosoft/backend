//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/http/request"
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

// TestChangeUserRole tests the user role change endpoint
func (suite *ManagementRouterTestSuite) TestChangeUserRole() {
	// Get admin credentials
	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

	// Now register a regular user to change their role
	regularUserName := "usertest"
	regularUserEmail := "user@example.com"
	regularUserPassword := "userPassword123"
	userID, _, userRefreshToken := suite.helpers.RegisterAndLoginUser(regularUserName, regularUserEmail, regularUserPassword)

	// 1. Test: Successful role change (user -> moderator) by admin
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

// Entry point for the Management Router test suite
func TestManagementRouterSuite(t *testing.T) {
	suite.Run(t, new(ManagementRouterTestSuite))
}
