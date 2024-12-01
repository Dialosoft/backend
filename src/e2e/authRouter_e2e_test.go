//go:build e2e

package e2e

import (
	"bytes"
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

type AuthRouterTestSuite struct {
	E2eTestSuite
	server *TestServer
}

// SetupSuite is called before any Router tests in the suite are run.
func (suite *AuthRouterTestSuite) SetupSuite() {
	suite.E2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
}

// TearDownSuite is called after all Router tests in the suite have been run, regardless of whether they passed or failed.
func (suite *AuthRouterTestSuite) TearDownSuite() {
	// Shutdown the server after a set of tests
    err := suite.server.Shutdown()
    suite.NoError(err)

    suite.E2eTestSuite.TearDownSuite()
}

func (suite *AuthRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)
}

func (suite *AuthRouterTestSuite) TearDownTest() {

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

func (suite *AuthRouterTestSuite) makeRequest(method, path string, payload interface{}, expectedStatus int) (*http.Response, error) {
	jsonPayload, err := json.Marshal(payload)
	suite.NoError(err)

	req, err := http.NewRequest(method, GetAPIBasePath(path), bytes.NewBuffer(jsonPayload))
	suite.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.server.App.Test(req, -1) // Remove timeout
	suite.NoError(err)

	suite.Equal(expectedStatus, resp.StatusCode)
	return resp, nil
}

func (suite *AuthRouterTestSuite) TestRegisterEndpoint() {
	t := suite.T()

	t.Run("should register new user successfully", func(t *testing.T) {
		// Prepare the registration payload
		registerPayload := request.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		resp, _ := suite.makeRequest("POST", "/auth/register", registerPayload, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)

		// close response body after reading it to avoid resource leak
		defer resp.Body.Close()

		var standardResponse struct {
			Message string                    `json:"message"`
			Data    response.RegisterResponse `json:"data"`
		}

		err = json.Unmarshal(body, &standardResponse)
		suite.NoError(err)

		// Now you can access the data through standardResponse.Data
		suite.NotEmpty(standardResponse.Data.UserID)
		suite.NotEmpty(standardResponse.Data.AccessToken)
		suite.NotEmpty(standardResponse.Data.RefreshToken)

		// Verify user in the database
		var user models.UserEntity
		err = suite.db.Where("email = ?", registerPayload.Email).First(&user).Error
		suite.NoError(err)
		suite.Equal(registerPayload.Username, user.Username)
		suite.Equal(registerPayload.Email, user.Email)

		// Verify that there are no related cache tickets
		cacheKey := "user:" + user.ID.String()
		exists, err := suite.rdClient.Exists(suite.testCtx, cacheKey).Result()
		suite.NoError(err)
		suite.Equal(int64(0), exists)
	})

	t.Run("should return error when registering with existing email", func(t *testing.T) {
		// Register a user first to have an existing email
		registerPayload := request.RegisterRequest{
			Username: "existinguser",
			Email:    "existing@example.com",
			Password: "password123",
		}
		_, err := suite.makeRequest("POST", "/auth/register", registerPayload, http.StatusOK)
		suite.NoError(err)

		// Try to register a user with the same email
		registerPayload2 := request.RegisterRequest{
			Username: "testuser2",
			Email:    "existing@example.com",
			Password: "password456",
		}
		resp, err := suite.makeRequest("POST", "/auth/register", registerPayload2, http.StatusConflict) // Esperamos un conflicto (409)
		suite.NoError(err)                                                                              // El error debería ser nulo porque esperamos un código de estado 409

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message.Adjust "Email Already Exists" according to your real answer.
		suite.Equal("CONFLICT", errorResponse.ErrorMessage)
	})
}

func TestAuthRouterSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(AuthRouterTestSuite))
}
