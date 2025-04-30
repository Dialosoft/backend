//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/app/database"
	"github.com/Dialosoft/src/domain/models"
	// "github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// ForumRouterTestSuite is a suite of end-to-end tests for forum routes
type ForumRouterTestSuite struct {
	BaseE2eTestSuite
	server          *TestServer
	helpers         TestHelpers
	testCategoryIDs []string // IDs categories for tests
}

// SetupSuite is called before any Forum Router tests in the suite are run.
func (suite *ForumRouterTestSuite) SetupSuite() {
	suite.BaseE2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
	suite.helpers = TestHelpers{
		Suite:  &suite.Suite,
		Server: suite.server,
	}
}

// TearDownSuite is called after all Forum Router tests in the suite have been run.
func (suite *ForumRouterTestSuite) TearDownSuite() {
	// Shutdown the server after a set of tests
	err := suite.server.Shutdown()
	suite.NoError(err)

	suite.BaseE2eTestSuite.TearDownSuite()
}

func (suite *ForumRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)

	// Verify if the roles were created
	var count int64
	err = suite.db.Model(&models.RoleEntity{}).Count(&count).Error
	suite.NoError(err)
	suite.Greater(count, int64(0), "No roles were created")

	// Create test categories using the helper
	categoryIDs, err := suite.helpers.CreateTestCategories(suite.db)
	suite.NoError(err)
	suite.testCategoryIDs = categoryIDs
	suite.NotEmpty(suite.testCategoryIDs, "Test categories were not created")
}

func (suite *ForumRouterTestSuite) TearDownTest() {
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

// TestGetForumsByCategoryID tests retrieving forums by category ID
func (suite *ForumRouterTestSuite) TestGetForumsByCategoryID() {
	t := suite.T()

	t.Run("should successfully get forums by category ID", func(t *testing.T) {
		// Create an admin user for creating forums
		adminID, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)
		suite.NotEmpty(adminID)

		// Get a category
		var category models.Category
		result := suite.db.First(&category)
		suite.NoError(result.Error)

		// Create a forum in that category
		forumName := "Test Forum"
		forumDescription := "Test Forum Description"
		forumType := "information"
		forumIsActive := true
		forumRolesAllowed := []string{
			"anonymous", "user", "administrator", "moderator",
		}
		categoryID := category.ID.String()

		createForumPayload := request.NewForum{
			Name:         &forumName,
			Description:  &forumDescription,
			Type:         &forumType,
			IsActive:     &forumIsActive,
			CategoryID:   &categoryID,
			RolesAllowed: forumRolesAllowed,
		}

		// Options with admin tokens for authorization
		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Create the forum using the API as admin
		resp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPost,
			"/forums/protected/create-new-forum",
			createForumPayload,
			http.StatusCreated,
			options,
		)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var createResponse response.StandardResponse
		err = json.Unmarshal(body, &createResponse)
		suite.NoError(err)

		// Verify that the forum was created
		data, ok := createResponse.Data.(map[string]interface{})
		suite.True(ok)
		forumID, ok := data["id"].(string)
		suite.True(ok)
		suite.NotEmpty(forumID)

		// Now test the endpoint to get forums by category ID
		respGet, err := suite.helpers.MakeRequest(
			http.MethodGet,
			fmt.Sprintf("/forums/get-forums-by-category-id/%s", category.ID.String()),
			nil,
			http.StatusOK,
		)
		suite.NoError(err)

		bodyGet, err := io.ReadAll(respGet.Body)
		suite.NoError(err)
		defer respGet.Body.Close()

		var getResponse response.StandardResponse
		err = json.Unmarshal(bodyGet, &getResponse)
		suite.NoError(err)

		// Verify that forums are returned correctly
		forumsData, ok := getResponse.Data.([]interface{})
		suite.True(ok)
		suite.Greater(len(forumsData), 0, "No forums were returned")

		// Verify that the created forum is in the response
		foundForum := false
		for _, forumData := range forumsData {
			forum, ok := forumData.(map[string]interface{})
			suite.True(ok)

			if id, ok := forum["id"].(string); ok && id == forumID {
				suite.Equal(forumName, forum["name"], "Forum name does not match")
				suite.Equal(forumDescription, forum["description"], "Forum description does not match")
				suite.Equal(forumType, forum["type"], "Forum type does not match")
				suite.Equal(categoryID, forum["categoryId"], "Category ID does not match")
				foundForum = true
				break
			}
		}

		suite.True(foundForum, "Created forum was not found in the response")
	})

	t.Run("should return an empty list for categories without forums", func(t *testing.T) {
		// Create a new category without forums
		newCategory := models.Category{
			Name:        "Empty Category",
			Description: "Category with no forums",
		}

		result := suite.db.Create(&newCategory)
		suite.NoError(result.Error)

		// Test the endpoint for a category with no forums
		respGet, err := suite.helpers.MakeRequest(
			http.MethodGet,
			fmt.Sprintf("/forums/get-forums-by-category-id/%s", newCategory.ID.String()),
			nil,
			http.StatusOK,
		)
		suite.NoError(err)

		bodyGet, err := io.ReadAll(respGet.Body)
		suite.NoError(err)
		defer respGet.Body.Close()

		var getResponse response.StandardResponse
		err = json.Unmarshal(bodyGet, &getResponse)
		suite.NoError(err)

		// Verify that an empty list or null is returned
		if getResponse.Data != nil {
			forumsData, ok := getResponse.Data.([]interface{})
			suite.True(ok)
			suite.Equal(0, len(forumsData), "Should return an empty list for categories without forums")
		}
	})

	t.Run("should return a 400 error for invalid category ID", func(t *testing.T) {
		// Test with an invalid category ID
		respGet, err := suite.helpers.MakeRequest(
			http.MethodGet,
			"/forums/get-forums-by-category-id/invalid-uuid",
			nil,
			http.StatusBadRequest,
		)
		suite.NoError(err)

		bodyGet, err := io.ReadAll(respGet.Body)
		suite.NoError(err)
		defer respGet.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(bodyGet, &errorResponse)
		suite.NoError(err)

		suite.Equal("ID provided is not a valid UUID type", errorResponse.ErrorMessage)
	})
}

// TestUpdateForum tests the forum update functionality
func (suite *ForumRouterTestSuite) TestUpdateForum() {
	t := suite.T()

	t.Run("should update a forum successfully with admin permissions", func(t *testing.T) {
		// Create an admin user for forum operations
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Get a category for creating the initial forum
		var category models.Category
		result := suite.db.First(&category)
		suite.NoError(result.Error)

		// Create a forum to update later
		initialForumName := "Initial Forum"
		initialForumDescription := "Initial Forum Description"
		initialForumType := "discussion"
		initialForumIsActive := true
		initialRolesAllowed := []string{"administrator", "moderator"}
		categoryID := category.ID.String()

		// Create initial forum via the API
		createForumPayload := request.NewForum{
			Name:         &initialForumName,
			Description:  &initialForumDescription,
			Type:         &initialForumType,
			IsActive:     &initialForumIsActive,
			RolesAllowed: initialRolesAllowed,
			CategoryID:   &categoryID,
		}

		// Options with admin tokens for authorization
		options := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Create the forum using the API as admin
		createResp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPost,
			"/forums/protected/create-new-forum",
			createForumPayload,
			http.StatusCreated,
			options,
		)
		suite.NoError(err)

		createBody, err := io.ReadAll(createResp.Body)
		suite.NoError(err)
		defer createResp.Body.Close()

		var createResponse response.StandardResponse
		err = json.Unmarshal(createBody, &createResponse)
		suite.NoError(err)

		// Get the forum ID
		createData, ok := createResponse.Data.(map[string]interface{})
		suite.True(ok)
		forumID, ok := createData["id"].(string)
		suite.True(ok)
		suite.NotEmpty(forumID)

		// Prepare data for updating the forum
		updatedForumName := "Updated Forum"
		updatedForumDescription := "Updated Forum Description"
		updatedForumType := "announcement"
		updatedForumIsActive := false
		updatedRolesAllowed := []string{"anonymous", "user", "moderator", "administrator"}

		updateForumPayload := request.NewForum{
			Name:         &updatedForumName,
			Description:  &updatedForumDescription,
			Type:         &updatedForumType,
			IsActive:     &updatedForumIsActive,
			RolesAllowed: updatedRolesAllowed,
			CategoryID:   &categoryID, // Keep the same category
		}

		// Update the forum via the API
		updateResp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			fmt.Sprintf("/forums/protected/update-forum/%s", forumID),
			updateForumPayload,
			http.StatusOK,
			options,
		)
		suite.NoError(err)

		updateBody, err := io.ReadAll(updateResp.Body)
		suite.NoError(err)
		defer updateResp.Body.Close()

		var updateResponse response.StandardResponse
		err = json.Unmarshal(updateBody, &updateResponse)
		suite.NoError(err)
		suite.Equal("UPDATED", updateResponse.Message)

		// Verify that the forum was updated in the database
		var updatedForum models.Forum
		err = suite.db.First(&updatedForum, "id = ?", forumID).Error
		suite.NoError(err)
		suite.Equal(updatedForumName, updatedForum.Name)
		suite.Equal(updatedForumDescription, updatedForum.Description)
		suite.Equal(updatedForumType, updatedForum.Type)
		suite.Equal(updatedForumIsActive, updatedForum.IsActive)
	})

	t.Run("should fail to update a forum without admin permissions", func(t *testing.T) {
		// First create a forum as admin
		_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

		// Get a category for the forum
		var category models.Category
		result := suite.db.First(&category)
		suite.NoError(result.Error)

		// Create a forum to attempt to update later
		initialForumName := "Admin Forum"
		initialForumDescription := "Admin Forum Description"
		initialForumType := "discussion"
		initialForumIsActive := true
		categoryID := category.ID.String()

		createForumPayload := request.NewForum{
			Name:        &initialForumName,
			Description: &initialForumDescription,
			Type:        &initialForumType,
			IsActive:    &initialForumIsActive,
			CategoryID:  &categoryID,
		}

		// Options with admin tokens for authorization
		adminOptions := &RequestOptions{
			Token:        adminToken,
			RefreshToken: adminRefreshToken,
		}

		// Create the forum as admin
		createResp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPost,
			"/forums/protected/create-new-forum",
			createForumPayload,
			http.StatusCreated,
			adminOptions,
		)
		suite.NoError(err)

		createBody, err := io.ReadAll(createResp.Body)
		suite.NoError(err)
		defer createResp.Body.Close()

		var createResponse response.StandardResponse
		err = json.Unmarshal(createBody, &createResponse)
		suite.NoError(err)

		// Get the forum ID
		createData, ok := createResponse.Data.(map[string]interface{})
		suite.True(ok)
		forumID, ok := createData["id"].(string)
		suite.True(ok)
		suite.NotEmpty(forumID)

		// Now register a regular user
		_, userToken, userRefreshToken := suite.helpers.RegisterAndLoginUser("regularuser", "regular@example.com", "password123")

		// Attempt to update the forum as a regular user
		updatedForumName := "Unauthorized Update"
		updatedForumDescription := "This should not be updated"

		updateForumPayload := request.NewForum{
			Name:        &updatedForumName,
			Description: &updatedForumDescription,
		}

		userOptions := &RequestOptions{
			Token:        userToken,
			RefreshToken: userRefreshToken,
		}

		// Attempt to update as regular user - should get 403 Forbidden
		updateResp, err := suite.helpers.MakeRequestWithOptions(
			http.MethodPut,
			fmt.Sprintf("/forums/protected/update-forum/%s", forumID),
			updateForumPayload,
			http.StatusForbidden, // Expect permission denied
			userOptions,
		)
		suite.NoError(err)
		defer updateResp.Body.Close()

		// Verify that the forum was not updated in the database
		var forum models.Forum
		err = suite.db.First(&forum, "id = ?", forumID).Error
		suite.NoError(err)
		suite.Equal(initialForumName, forum.Name, "Forum name should not have changed")
		suite.Equal(initialForumDescription, forum.Description, "Forum description should not have changed")
	})

	// TODO: Add endpoint for parcial update of forum - PATCH
	// t.Run("should update some fields of a forum", func(t *testing.T) {
	// 	// Create an admin user for forum operations
	// 	_, adminToken, adminRefreshToken := suite.helpers.CreateAdminUserAndLogin(suite.db)

	// 	// Get a category for creating the initial forum
	// 	var category models.Category
	// 	result := suite.db.First(&category)
	// 	suite.NoError(result.Error)

	// 	// Create a forum to update later
	// 	initialForumName := "Partial Update Forum"
	// 	initialForumDescription := "Partial Update Forum Description"
	// 	initialForumType := "discussion"
	// 	initialForumIsActive := true
	// 	categoryID := category.ID.String()

	// 	// Create initial forum via the API
	// 	createForumPayload := request.NewForum{
	// 		Name:        &initialForumName,
	// 		Description: &initialForumDescription,
	// 		Type:        &initialForumType,
	// 		IsActive:    &initialForumIsActive,
	// 		CategoryID:  &categoryID,
	// 	}

	// 	// Options with admin tokens for authorization
	// 	options := &RequestOptions{
	// 		Token:        adminToken,
	// 		RefreshToken: adminRefreshToken,
	// 	}

	// 	// Create the forum using the API as admin
	// 	createResp, err := suite.helpers.MakeRequestWithOptions(
	// 		http.MethodPost,
	// 		"/forums/protected/create-new-forum",
	// 		createForumPayload,
	// 		http.StatusCreated,
	// 		options,
	// 	)
	// 	suite.NoError(err)

	// 	createBody, err := io.ReadAll(createResp.Body)
	// 	suite.NoError(err)
	// 	defer createResp.Body.Close()

	// 	var createResponse response.StandardResponse
	// 	err = json.Unmarshal(createBody, &createResponse)
	// 	suite.NoError(err)

	// 	// Get the forum ID
	// 	createData, ok := createResponse.Data.(map[string]interface{})
	// 	suite.True(ok)
	// 	forumID, ok := createData["id"].(string)
	// 	suite.True(ok)
	// 	suite.NotEmpty(forumID)

	// 	// Prepare partial update - only changing some fields
	// 	updatedForumName := "New Name Only"
	// 	updatedRolesAllowed := []string{"anonymous", "user", "administrator"}

	// 	// Only include fields we want to update
	// 	updateForumPayload := request.NewForum{
	// 		Name:         &updatedForumName,
	// 		RolesAllowed: updatedRolesAllowed,
	// 		// Other fields intentionally omitted to test partial updates
	// 	}

	// 	// Update the forum via the API
	// 	updateResp, err := suite.helpers.MakeRequestWithOptions(
	// 		http.MethodPut,
	// 		fmt.Sprintf("/forums/protected/update-forum/%s", forumID),
	// 		updateForumPayload,
	// 		http.StatusOK,
	// 		options,
	// 	)
	// 	suite.NoError(err)

	// 	updateBody, err := io.ReadAll(updateResp.Body)
	// 	suite.NoError(err)
	// 	defer updateResp.Body.Close()

	// 	var updateResponse response.StandardResponse
	// 	err = json.Unmarshal(updateBody, &updateResponse)
	// 	suite.NoError(err)
	// 	suite.Equal("UPDATED", updateResponse.Message)

	// 	// Verify that only specified fields were updated in the database
	// 	var updatedForum models.Forum
	// 	err = suite.db.First(&updatedForum, "id = ?", forumID).Error
	// 	suite.NoError(err)
	// 	suite.Equal(updatedForumName, updatedForum.Name, "Name should have been updated")
	// 	suite.Equal(initialForumType, updatedForum.Type, "Type should remain unchanged")
	// 	suite.Equal(initialForumIsActive, updatedForum.IsActive, "IsActive should remain unchanged")

	// 	// Check that it matches the roles allowed we set
	// 	rolesEqual := suite.Equal(len(updatedRolesAllowed), len(updatedForum.RolesAllowed), "RolesAllowed count should match")
	// 	if rolesEqual {
	// 		for i, role := range updatedRolesAllowed {
	// 			suite.Equal(role, updatedForum.RolesAllowed[i], "RolesAllowed entry should match")
	// 		}
	// 	}
	// })
}

// Entry point for the Forum Router test suite
func TestForumRouterSuite(t *testing.T) {
	suite.Run(t, new(ForumRouterTestSuite))
}
