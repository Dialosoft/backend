//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/app/database"
	"github.com/Dialosoft/src/domain/models"
)

// RoleRouterTestSuite is a suite of end-to-end tests for role routes
type RoleRouterTestSuite struct {
	BaseE2eTestSuite
	server  *TestServer
	helpers TestHelpers
}

// SetupSuite is called before any Role Router tests in the suite are run.
func (suite *RoleRouterTestSuite) SetupSuite() {
	suite.BaseE2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
	suite.helpers = TestHelpers{
		Suite:  &suite.Suite,
		Server: suite.server,
	}
}

// TearDownSuite is called after all Role Router tests in the suite have been run.
func (suite *RoleRouterTestSuite) TearDownSuite() {
	// Shutdown the server after all tests
	err := suite.server.Shutdown()
	suite.NoError(err)

	suite.BaseE2eTestSuite.TearDownSuite()
}

func (suite *RoleRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)

	// Verify if the roles were created
	var count int64
	err = suite.db.Model(&models.RoleEntity{}).Count(&count).Error
	suite.NoError(err)
	suite.Greater(count, int64(0), "No roles were created")
}

func (suite *RoleRouterTestSuite) TearDownTest() {
	// Truncate tables maintaining the structure and restrictions
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

// TestGetRolePermissionsByRoleID tests retrieving role permissions by role ID
func (suite *RoleRouterTestSuite) TestGetRolePermissionsByRoleID() {
	t := suite.T()

	t.Run("should successfully get role permissions by role ID", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Get a role ID from the database
		var role models.RoleEntity
		err := suite.db.Where("role_type = ?", "administrator").First(&role).Error
		suite.NoError(err)
		suite.NotEmpty(role.ID)

		// Get role permissions using the API
		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Make the API request to get role permissions
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodGet,
			fmt.Sprintf("/roles/protected/get-role-permissions-by-id/%s", role.ID),
			nil,
			http.StatusOK,
			options,
		)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var getResponse response.StandardResponse
		err = json.Unmarshal(body, &getResponse)
		suite.NoError(err)

		// Verify that permissions are returned correctly
		permissionsData, ok := getResponse.Data.(map[string]interface{})
		suite.True(ok)

		// Check for the presence of permission fields
		_, hasManageCategories := permissionsData["CanManageCategories"]
		_, hasManageForums := permissionsData["CanManageForums"]
		_, hasManageRoles := permissionsData["CanManageRoles"]
		_, hasManageUsers := permissionsData["CanManageUsers"]

		suite.True(hasManageCategories, "Missing canManageCategories field")
		suite.True(hasManageForums, "Missing canManageForums field")
		suite.True(hasManageRoles, "Missing canManageRoles field")
		suite.True(hasManageUsers, "Missing canManageUsers field")

		// For administrator, these permissions should be true
		suite.Equal(true, hasManageCategories)
		suite.Equal(true, hasManageForums)
		suite.Equal(true, hasManageRoles)
		suite.Equal(true, hasManageUsers)
	})

	t.Run("should return 404 when role ID does not exist", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Generate a non-existent UUID
		nonExistentID := uuid.New().String()

		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Make API request with non-existent ID
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodGet,
			fmt.Sprintf("/roles/protected/get-role-permissions-by-id/%s", nonExistentID),
			nil,
			http.StatusNotFound,
			options,
		)
		suite.NoError(err)
		defer resp.Body.Close()
	})

	t.Run("should return 400 when role ID is invalid", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Make API request with invalid ID
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodGet,
			"/roles/protected/get-role-permissions-by-id/invalid-uuid",
			nil,
			http.StatusBadRequest,
			options,
		)
		suite.NoError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)
		suite.Equal("ID provided is not a valid UUID type", errorResponse.ErrorMessage)
	})
}

// TestSetRolePermissionsByRoleID tests updating role permissions by role ID
func (suite *RoleRouterTestSuite) TestSetRolePermissionsByRoleID() {
	t := suite.T()

	t.Run("should successfully update role permissions", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Get a role ID that's not administrator (to avoid breaking admin tests)
		var role models.RoleEntity
		err := suite.db.Where("role_type = ?", "user").First(&role).Error
		suite.NoError(err)
		suite.NotEmpty(role.ID)

		// Prepare updated permissions
		canManageCategories := true
		canManageForums := true
		canManageRoles := false
		canManageUsers := false

		updatePermissionsPayload := request.NewRolePermissions{
			CanManageCategories: &canManageCategories,
			CanManageForums:     &canManageForums,
			CanManageRoles:      &canManageRoles,
			CanManageUsers:      &canManageUsers,
		}

		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Update the permissions using the API
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			fmt.Sprintf("/roles/protected/set-role-permissions-by-id/%s", role.ID),
			updatePermissionsPayload,
			http.StatusOK,
			options,
		)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var updateResponse response.StandardResponse
		err = json.Unmarshal(body, &updateResponse)
		suite.NoError(err)
		suite.Equal("UPDATED", updateResponse.Message)

		// Verify that the permissions were updated by getting them
		getResp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodGet,
			fmt.Sprintf("/roles/protected/get-role-permissions-by-id/%s", role.ID),
			nil,
			http.StatusOK,
			options,
		)
		suite.NoError(err)

		getBody, err := io.ReadAll(getResp.Body)
		suite.NoError(err)
		defer getResp.Body.Close()

		var getResponse response.StandardResponse
		err = json.Unmarshal(getBody, &getResponse)
		suite.NoError(err)

		// Verify permissions were updated correctly
		permissionsData, ok := getResponse.Data.(map[string]interface{})
		suite.True(ok)

		suite.Equal(canManageCategories, permissionsData["CanManageCategories"])
		suite.Equal(canManageForums, permissionsData["CanManageForums"])
		suite.Equal(canManageRoles, permissionsData["CanManageRoles"])
		suite.Equal(canManageUsers, permissionsData["CanManageUsers"])
	})

	t.Run("should return 403 when user doesn't have permission", func(t *testing.T) {
		// Register a regular user with no admin permissions
		_, userToken, userRefreshToken := suite.helpers.RegisterAndLoginUser("testuser", "test@example.com", "password123")

		// Get a role ID
		var role models.RoleEntity
		err := suite.db.Where("role_type = ?", "user").First(&role).Error
		suite.NoError(err)

		// Prepare updated permissions
		canManageCategories := true
		canManageForums := true

		updatePermissionsPayload := request.NewRolePermissions{
			CanManageCategories: &canManageCategories,
			CanManageForums:     &canManageForums,
		}

		options := &RequestOptions{
			Token:        userToken,
			RefreshToken: userRefreshToken,
		}

		// Try to update permissions as regular user
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			fmt.Sprintf("/roles/protected/set-role-permissions-by-id/%s", role.ID),
			updatePermissionsPayload,
			http.StatusForbidden,
			options,
		)
		suite.NoError(err)
		defer resp.Body.Close()
	})

	t.Run("should return 404 when role ID does not exist", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Generate a non-existent UUID
		nonExistentID := uuid.New().String()

		// Prepare updated permissions
		canManageCategories := true
		canManageForums := true

		updatePermissionsPayload := request.NewRolePermissions{
			CanManageCategories: &canManageCategories,
			CanManageForums:     &canManageForums,
		}

		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Make API request with non-existent ID
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			fmt.Sprintf("/roles/protected/set-role-permissions-by-id/%s", nonExistentID),
			updatePermissionsPayload,
			http.StatusNotFound,
			options,
		)
		suite.NoError(err)
		defer resp.Body.Close()
	})

	t.Run("should return 400 when role ID is invalid", func(t *testing.T) {
		// Create an admin user for authentication
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Prepare updated permissions
		canManageCategories := true
		canManageForums := true

		updatePermissionsPayload := request.NewRolePermissions{
			CanManageCategories: &canManageCategories,
			CanManageForums:     &canManageForums,
		}

		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Make API request with invalid ID
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			"/roles/protected/set-role-permissions-by-id/invalid-uuid",
			updatePermissionsPayload,
			http.StatusBadRequest,
			options,
		)
		suite.NoError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)
		suite.Equal("ID provided is not a valid UUID type", errorResponse.ErrorMessage)
	})
}

// Entry point for the Role Router test suite
func TestRoleRouterSuite(t *testing.T) {
	suite.Run(t, new(RoleRouterTestSuite))
}
