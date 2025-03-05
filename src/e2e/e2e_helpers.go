//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/utils/jsonWebToken"
)

// TestHelpers provides common test helper methods for e2e tests
type TestHelpers struct {
	*suite.Suite // Changed to pointer to avoid copying the mutex
	Server       *TestServer
}

// RequestOptions contains options for making HTTP requests
type RequestOptions struct {
	Headers      map[string]string
	Token        string
	RefreshToken string
}

// MakeRequest creates and sends an HTTP request to the test server
// It handles JSON marshaling, request creation, and validates the expected status code
func (h *TestHelpers) MakeRequest(method, path string, payload interface{}, expectedStatus int) (*http.Response, error) {
	return h.MakeRequestWithOptions(method, path, payload, expectedStatus, nil)
}

// MakeRequestWithOptions creates and sends an HTTP request to the test server with custom options
// It handles JSON marshaling, request creation with custom headers, and validates the expected status code
func (h *TestHelpers) MakeRequestWithOptions(method, path string, payload interface{}, expectedStatus int, options *RequestOptions) (*http.Response, error) {
	var req *http.Request
	var err error

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		h.NoError(err)
		req, err = http.NewRequest(method, GetAPIBasePath(path), bytes.NewBuffer(jsonPayload))
		h.NoError(err)
		req.Header.Set("Content-Type", "application/json")
	} else {
		// For requests without payload (GET, DELETE, etc.)
		req, err = http.NewRequest(method, GetAPIBasePath(path), nil)
		h.NoError(err)
	}

	// Apply custom headers if provided
	if options != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}

		// Add Authorization header if token is provided
		if options.Token != "" {
			req.Header.Set("Authorization", "Bearer "+options.Token)
		}

		// Add X-Refresh-Token header if refresh token is provided
		if options.RefreshToken != "" {
			req.Header.Set("X-Refresh-Token", options.RefreshToken)
		}
	}

	resp, err := h.Server.App.Test(req, -1) // No timeout
	h.NoError(err)

	h.Equal(expectedStatus, resp.StatusCode)
	return resp, nil
}

// LoginUser performs a login with the given credentials and returns the access token
func (h *TestHelpers) LoginUser(username, password string) string {
	// Login with the user
	loginPayload := request.LoginRequest{
		Username: username,
		Password: password,
	}

	// Make login request
	resp, err := h.MakeRequest(http.MethodPost, "/auth/login", loginPayload, http.StatusOK)
	h.NoError(err)
	defer resp.Body.Close()

	// Parse login response
	body, err := io.ReadAll(resp.Body)
	h.NoError(err)

	var loginResponse response.StandardResponse
	err = json.Unmarshal(body, &loginResponse)
	h.NoError(err)

	// Extract access token
	data := loginResponse.Data.(map[string]interface{})
	accessToken := data["accessToken"].(string)

	return accessToken
}

// LoginUserFull performs a login and returns user ID, access token, and refresh token
func (h *TestHelpers) LoginUserFull(username, password string) (string, string, string) {
	// Login with the user
	loginPayload := request.LoginRequest{
		Username: username,
		Password: password,
	}

	// Make login request
	resp, err := h.MakeRequest(http.MethodPost, "/auth/login", loginPayload, http.StatusOK)
	h.NoError(err)
	defer resp.Body.Close()

	// Parse login response
	body, err := io.ReadAll(resp.Body)
	h.NoError(err)

	var loginResponse response.StandardResponse
	err = json.Unmarshal(body, &loginResponse)
	h.NoError(err)

	// Check if data exists and isn't nil
	h.NotNil(loginResponse.Data, "Login response data should not be nil")

	// Extract data into LoginResponse type
	loginResponseData, ok := loginResponse.Data.(map[string]interface{})
	h.True(ok, "Failed to convert login response data to map[string]interface{}")

	// Extract access and refresh tokens - these are directly in the response
	accessToken, ok := loginResponseData["accessToken"].(string)
	h.True(ok, "Failed to extract accessToken from login response")

	refreshToken, ok := loginResponseData["refreshToken"].(string)
	h.True(ok, "Failed to extract refreshToken from login response")

	// Extract user ID from the JWT token using the utility function
	secretKey := h.Server.Config.JWTKey
	userID, err := jsonWebToken.ExtractUserIDFromToken(accessToken, secretKey)
	h.NoError(err, "Failed to extract user ID from JWT token")

	return userID, accessToken, refreshToken
}

// RegisterAndLoginUser registers a new user, logs them in, and returns ID and tokens
func (h *TestHelpers) RegisterAndLoginUser(username, email, password string) (string, string, string) {
	// Register a new user
	registerPayload := request.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Make registration request
	resp, err := h.MakeRequest(http.MethodPost, "/auth/register", registerPayload, http.StatusCreated)
	h.NoError(err)
	defer resp.Body.Close()

	// Login with the new user and get full information
	return h.LoginUserFull(username, password)
}

// GetRoleIDByType gets the role ID for a specific role type from the database
func (h *TestHelpers) GetRoleIDByType(db *gorm.DB, roleType string) string {
	var role models.RoleEntity
	result := db.Where("role_type = ?", roleType).First(&role)
	h.NoError(result.Error)
	return role.ID.String()
}

// Create admin user and login with it to get admin token
func (h *TestHelpers) CreateAdminUserAndLogin(db *gorm.DB) (string, string, string) {
	// Register an admin user first
	adminUserName := "admin_test_user"
	adminEmail := "admin@example.com"
	adminPassword := "adminPassword123"
	adminID, _, _ := h.RegisterAndLoginUser(adminUserName, adminEmail, adminPassword)

	// Get admin role ID
	adminRoleID := h.GetRoleIDByType(db, "administrator")

	// Get the admin user without admin role yet
	var adminUser models.UserEntity
	result := db.First(&adminUser, "id = ?", adminID)
	h.NoError(result.Error)

	// Convert admin role ID to UUID
	uuidAdminRoleID, err := uuid.Parse(adminRoleID)
	h.NoError(err)

	// Update the admin's role to be an actual administrator
	adminUser.RoleID = uuidAdminRoleID

	result = db.Save(&adminUser)
	h.NoError(result.Error)

	// Get a fresh admin token with admin role
	adminID, adminToken, adminRefreshToken := h.LoginUserFull(adminUserName, adminPassword)

	return adminID, adminToken, adminRefreshToken
}
