package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/domain/models"
)

// RedisRepository Mock
type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// CacheServiceTestSuite It is the test suite that groups the tests related to CacheService.
type CacheServiceTestSuite struct {
	suite.Suite
	mockRepo *MockRedisRepository
	service  CacheService
}

// SetupTest is executed before each general suite test.
func (suite *CacheServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockRedisRepository)
	suite.service = NewCacheService(suite.mockRepo)
}

// TearDownTest is executed after each general suite test.
func (suite *CacheServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CacheServiceTestSuite) TestInvalidateRefreshToken() {
	tests := []struct {
		name          string
		token         string
		mockReturnErr error
		expectedErr   string
	}{
		{
			name:          "success InvalidateRefreshToken",
			token:         "testToken",
			mockReturnErr: nil,
			expectedErr:   "",
		},
		{
			name:          "error setting token in cache",
			token:         "testToken",
			mockReturnErr: errors.New("cache error"),
			expectedErr:   "cache error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "blacklist:" + tt.token

			// Reset Mock's expectations before each subtest
			suite.mockRepo.ExpectedCalls = nil

			suite.mockRepo.On("Set", context.Background(), cacheKey, "true", time.Hour*720).Return(tt.mockReturnErr)

			err := suite.service.InvalidateRefreshToken(tt.token)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr)
			}
		})
	}
}

func (suite *CacheServiceTestSuite) TestIsTokenBlacklisted() {
	tests := []struct {
		name          string
		token         string
		mockExists    bool
		mockReturnErr error
		expected      bool
	}{
		{name: "token blacklisted", token: "blacklistedToken", mockExists: true, mockReturnErr: nil, expected: true},
		{name: "token not blacklisted", token: "validToken", mockExists: false, mockReturnErr: nil, expected: false},
		{name: "error checking token", token: "errorToken", mockExists: false, mockReturnErr: errors.New("redis error"), expected: false},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "blacklist:" + tt.token

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Exists", context.Background(), cacheKey).Return(tt.mockExists, tt.mockReturnErr)

			result := suite.service.IsTokenBlacklisted(tt.token)

			assert.Equal(suite.T(), tt.expected, result)
		})
	}
}

func (suite *CacheServiceTestSuite) TestSetUserInfoByID() {
	userID := uuid.New()

	tests := []struct {
		name          string
		userEntity    *models.UserEntity
		mockReturnErr error
		expectedErr   string
	}{
		{name: "success set user info", userEntity: &models.UserEntity{Username: "testuser"}, mockReturnErr: nil, expectedErr: ""},
		{name: "error setting user info", userEntity: &models.UserEntity{Username: "testuser"}, mockReturnErr: errors.New("redis error"), expectedErr: "redis error"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "user:" + userID.String()

			// Marshal the user entity to JSON and remove the password.
			jsonData, err := json.Marshal(tt.userEntity)
			assert.NoError(suite.T(), err)

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Set", context.Background(), cacheKey, string(jsonData), time.Hour*24).Return(tt.mockReturnErr)

			err = suite.service.SetUserInfoByID(userID, tt.userEntity)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr)
			}
		})
	}
}

func (suite *CacheServiceTestSuite) TestGetUserInfoByID() {
	userID := uuid.New()

	tests := []struct {
		name          string
		mockReturn    string
		mockReturnErr error
		expected      *models.UserEntity
		expectedErr   string
	}{
		{
			name:          "success retrieving user info",
			mockReturn:    `{"id":"` + userID.String() + `","username":"testuser"}`,
			mockReturnErr: nil,
			expected:      &models.UserEntity{ID: userID, Username: "testuser"},
			expectedErr:   "",
		},
		{
			name:          "error retrieving user info from cache",
			mockReturn:    "",
			mockReturnErr: errors.New("redis error"),
			expected:      nil,
			expectedErr:   "redis error",
		},
		{
			name:          "error unmarshalling user info",
			mockReturn:    `invalid json`,
			mockReturnErr: nil,
			expected:      nil,
			expectedErr:   "invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "user:" + userID.String()

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Get", context.Background(), cacheKey).Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetUserInfoByID(userID)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expected, result)
			} else {
				assert.Nil(suite.T(), result)
				assert.EqualError(suite.T(), err, tt.expectedErr)
			}
		})
	}
}

func (suite *CacheServiceTestSuite) TestSetRefreshTokenByID() {
	userID := uuid.New()

	tests := []struct {
		name          string
		token         string
		mockReturnErr error
		expectedErr   string
	}{
		{name: "success setting refresh token", token: "testRefreshToken", mockReturnErr: nil, expectedErr: ""},
		{name: "error setting refresh token", token: "testRefreshToken", mockReturnErr: errors.New("redis error"), expectedErr: "redis error"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "refreshToken:" + userID.String()

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Set", context.Background(), cacheKey, tt.token, time.Hour*120).Return(tt.mockReturnErr)

			err := suite.service.SetRefreshTokenByID(userID, tt.token)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr)
			}
		})
	}
}

func (suite *CacheServiceTestSuite) TestGetRefreshTokenByID() {
	userID := uuid.New()

	tests := []struct {
		name          string
		mockReturn    string
		mockReturnErr error
		expected      string
		expectedErr   string
	}{
		{
			name:          "success retrieving refresh token",
			mockReturn:    "testRefreshToken",
			mockReturnErr: nil,
			expected:      "testRefreshToken",
			expectedErr:   "",
		},
		{
			name:          "error retrieving refresh token from cache",
			mockReturn:    "",
			mockReturnErr: errors.New("redis error"),
			expected:      "",
			expectedErr:   "redis error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "refreshToken:" + userID.String()

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Get", context.Background(), cacheKey).Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetRefreshTokenByID(userID)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expected, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr)
				assert.Empty(suite.T(), result)
			}
		})
	}
}

// TestDeleteRefreshTokenByID prueba el método DeleteRefreshTokenByID.
func (suite *CacheServiceTestSuite) TestDeleteRefreshTokenByID() {
	userID := uuid.New()

	tests := []struct {
		name          string
		mockReturnErr error
		expectedErr   string
	}{
		{name: "success deleting refresh token", mockReturnErr: nil, expectedErr: ""},
		{name: "error deleting refresh token", mockReturnErr: errors.New("redis error"), expectedErr: "redis error"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cacheKey := "refreshToken:" + userID.String()

			suite.mockRepo.ExpectedCalls = nil
			suite.mockRepo.On("Delete", context.Background(), cacheKey).Return(tt.mockReturnErr)

			err := suite.service.DeleteRefreshTokenByID(userID)

			if tt.expectedErr == "" {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr)
			}
		})
	}
}

// Create a new instance of the cacsevicetestsuite structure and executes the test suite
func TestCacheServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CacheServiceTestSuite))
}
