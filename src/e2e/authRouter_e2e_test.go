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

type AuthRouterTestSuite struct {
	BaseE2eTestSuite
	server *TestServer
	helpers TestHelpers
}

// SetupSuite is called before any Router tests in the suite are run.
func (suite *AuthRouterTestSuite) SetupSuite() {
	suite.BaseE2eTestSuite.SetupSuite()
	suite.server = NewTestServer(suite.db, suite.rdClient)
	suite.helpers = TestHelpers{
		Suite:  &suite.Suite,
		Server: suite.server,
	}
}

// TearDownSuite is called after all Router tests in the suite have been run, regardless of whether they passed or failed.
func (suite *AuthRouterTestSuite) TearDownSuite() {
	// Shutdown the server after a set of tests
	err := suite.server.Shutdown()
	suite.NoError(err)

	suite.BaseE2eTestSuite.TearDownSuite()
}

func (suite *AuthRouterTestSuite) SetupTest() {
	// Create predetermined roles at the beginning of each test
	err := database.CreateDefaultRoles(suite.db)
	suite.NoError(err)

	// Verificar que se crearon los roles
	var count int64
	err = suite.db.Model(&models.RoleEntity{}).Count(&count).Error
	suite.NoError(err)
	suite.Greater(count, int64(0), "No roles were created")
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

func (suite *AuthRouterTestSuite) TestRegisterEndpoint() {
	t := suite.T()

	t.Run("should register new user successfully", func(t *testing.T) {
		// Prepare the registration payload
		registerPayload := request.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		resp, _ := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)

		// close response body after reading it to avoid resource leak
		defer resp.Body.Close()

		var registerResponse struct {
			Data response.RegisterResponse `json:"data"`
			Message string `json:"message"`
		}
		err = json.Unmarshal(body, &registerResponse)
		suite.NoError(err)

		suite.NotEmpty(registerResponse.Data.UserID)
		suite.NotEmpty(registerResponse.Data.AccessToken)
		suite.NotEmpty(registerResponse.Data.RefreshToken)

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
		_, err := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload, http.StatusOK)
		suite.NoError(err)

		// Try to register a user with the same email
		registerPayload2 := request.RegisterRequest{
			Username: "testuser2",
			Email:    "existing@example.com",
			Password: "password456",
		}
		resp, err := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload2, http.StatusConflict) // Esperamos un conflicto (409)
		suite.NoError(err)

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

func (suite *AuthRouterTestSuite) TestLoginEndpoint() {
	t := suite.T()

	t.Run("should login user successfully", func(t *testing.T) {
		// First, register a user
		registerPayload := request.RegisterRequest{
			Username: "loginuser",
			Email:    "login@example.com",
			Password: "password123",
		}
		_, err := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload, http.StatusOK)
		suite.NoError(err)

		// Now attempt to login with this user
		loginPayload := request.LoginRequest{
			Username: "loginuser",
			Password: "password123",
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/login", loginPayload, http.StatusOK)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var loginResponse struct {
			Data response.LoginResponse `json:"data"`
			Message string `json:"message"`
		}
		err = json.Unmarshal(body, &loginResponse)
		suite.NoError(err)

		// Verify the token data is returned
		suite.NotEmpty(loginResponse.Data.AccessToken)
		suite.NotEmpty(loginResponse.Data.RefreshToken)
		suite.Equal("Successfully logged in", loginResponse.Message)
	})

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		// Register a user first
		registerPayload := request.RegisterRequest{
			Username: "failuser",
			Email:    "fail@example.com",
			Password: "password123",
		}
		_, err := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload, http.StatusOK)
		suite.NoError(err)

		// Attempt login with wrong password
		loginPayload := request.LoginRequest{
			Username: "failuser",
			Password: "wrongpassword",
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/login", loginPayload, http.StatusUnauthorized)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message
		suite.Equal("UNAUTHORIZED", errorResponse.ErrorMessage)
	})

	t.Run("should fail with non-existing user", func(t *testing.T) {
		loginPayload := request.LoginRequest{
			Username: "nonexistentuser",
			Password: "password123",
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/login", loginPayload, http.StatusUnauthorized)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message
		suite.Equal("UNAUTHORIZED", errorResponse.ErrorMessage)
	})
}

func (suite *AuthRouterTestSuite) TestRefreshTokenEndpoint() {
	t := suite.T()

	t.Run("should refresh token successfully", func(t *testing.T) {
		// First, register a user to get valid tokens
		registerPayload := request.RegisterRequest{
			Username: "refreshuser",
			Email:    "refresh@example.com",
			Password: "password123",
		}
		resp, err := suite.helpers.MakeRequest("POST", "/auth/register", registerPayload, http.StatusOK)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var registerResponse struct {
			Data response.RegisterResponse `json:"data"`
			Message string `json:"message"`
		}
		err = json.Unmarshal(body, &registerResponse)
		suite.NoError(err)

		refreshToken := registerResponse.Data.RefreshToken
		suite.NotEmpty(refreshToken)

		// Now attempt to refresh the token
		refreshPayload := request.RefreshToken{
			Refresh: refreshToken,
		}

		resp, err = suite.helpers.MakeRequest("POST", "/auth/refresh-token", refreshPayload, http.StatusOK)
		suite.NoError(err)

		body, err = io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var refreshResponse struct {
			Data response.RefreshTokenResponse `json:"data"`
			Message string `json:"message"`
		}
		err = json.Unmarshal(body, &refreshResponse)
		suite.NoError(err)

		// Verify the new access token is returned
		suite.NotEmpty(refreshResponse.Data.AccessToken)
		suite.Equal("successfully refreshed", refreshResponse.Message)
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		// Attempt to refresh with an invalid token
		refreshPayload := request.RefreshToken{
			Refresh: "invalid.refresh.token",
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/refresh-token", refreshPayload, http.StatusUnauthorized)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message
		suite.Equal("UNAUTHORIZED", errorResponse.ErrorMessage)
	})

	t.Run("should fail with empty refresh token", func(t *testing.T) {
		// Attempt to refresh with an empty token
		refreshPayload := request.RefreshToken{
			Refresh: "",
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/refresh-token", refreshPayload, http.StatusUnauthorized)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message
		suite.Equal("UNAUTHORIZED", errorResponse.ErrorMessage)
	})

	t.Run("should fail with expired refresh token", func(t *testing.T) {
		// This is an example of an expired JWT token (created with "test-jwt-key" and expired in 2020)
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1NzYyMzkwMjJ9.xOsZT8J4xdQZsnPQKrz17SMmGxlygcKOCiXRgBJEt3M"
		
		refreshPayload := request.RefreshToken{
			Refresh: expiredToken,
		}

		resp, err := suite.helpers.MakeRequest("POST", "/auth/refresh-token", refreshPayload, http.StatusUnauthorized)
		suite.NoError(err)

		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)
		defer resp.Body.Close()

		var errorResponse response.StandardError
		err = json.Unmarshal(body, &errorResponse)
		suite.NoError(err)

		// Verify the error message
		suite.Equal("UNAUTHORIZED", errorResponse.ErrorMessage)
	})
}

func TestAuthRouterSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(AuthRouterTestSuite))
}
