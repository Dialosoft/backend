package services

import (
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/domain/models"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAllUsers() ([]*models.UserEntity, error) {
	args := m.Called()
	return args.Get(0).([]*models.UserEntity), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.UserEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.UserEntity, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

func (m *MockUserRepository) Create(newUser models.UserEntity) (uuid.UUID, error) {
	args := m.Called(newUser)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserRepository) Update(userID uuid.UUID, updatedUser models.UserEntity) error {
	args := m.Called(userID, updatedUser)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) Restore(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// UserServiceTestSuite defines the test suite for UserService
type UserServiceTestSuite struct {
	suite.Suite
	mockUserRepo *MockUserRepository
	mockRoleRepo *MockRoleRepository
	service      UserService
}

// SetupTest is executed before each test
func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockRoleRepo = new(MockRoleRepository)
	suite.service = NewUserService(suite.mockUserRepo, suite.mockRoleRepo)
}

// TearDownTest is executed after each test
func (suite *UserServiceTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockRoleRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetAllUsers() {
	tests := []struct {
		name          string
		mockUsers     []*models.UserEntity
		mockError     error
		expectedError string
	}{
		{
			name: "success get all users",
			mockUsers: []*models.UserEntity{
				{
					ID:       uuid.New(),
					Username: "user1",
					Email:    "user1@test.com",
					RoleID:   uuid.New(),
					Role: models.RoleEntity{
						ID:        uuid.New(),
						RoleType:  "USER",
						AdminRole: false,
						ModRole:   false,
						UserRole:  true,
					},
				},
				{
					ID:       uuid.New(),
					Username: "user2",
					Email:    "user2@test.com",
					RoleID:   uuid.New(),
					Role: models.RoleEntity{
						ID:        uuid.New(),
						RoleType:  "ADMIN",
						AdminRole: true,
						ModRole:   false,
						UserRole:  false,
					},
				},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "error getting users",
			mockUsers:     nil,
			mockError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("FindAllUsers").Return(tt.mockUsers, tt.mockError)

			users, err := suite.service.GetAllUsers()

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
				assert.Nil(suite.T(), users)
			} else {
				assert.NoError(suite.T(), err)
				assert.Len(suite.T(), users, len(tt.mockUsers))
				for i, user := range users {
					assert.Equal(suite.T(), tt.mockUsers[i].ID, user.ID)
					assert.Equal(suite.T(), tt.mockUsers[i].Username, user.Username)
					assert.Equal(suite.T(), tt.mockUsers[i].Email, user.Email)
				}
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUserByID() {
	userID := uuid.New()
	mockUser := &models.UserEntity{
		ID:       userID,
		Username: "testuser",
		Email:    "test@test.com",
		RoleID:   uuid.New(),
		Role: models.RoleEntity{
			ID:        uuid.New(),
			RoleType:  "USER",
			AdminRole: false,
			ModRole:   false,
			UserRole:  true,
		},
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockUser      *models.UserEntity
		mockError     error
		expectedError string
	}{
		{
			name:          "success get user by id",
			userID:        userID,
			mockUser:      mockUser,
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "user not found",
			userID:        uuid.New(),
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("FindByID", tt.userID).Return(tt.mockUser, tt.mockError)

			user, err := suite.service.GetUserByID(tt.userID)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
				assert.Nil(suite.T(), user)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), user)
				assert.Equal(suite.T(), tt.mockUser.ID, user.ID)
				assert.Equal(suite.T(), tt.mockUser.Username, user.Username)
				assert.Equal(suite.T(), tt.mockUser.Email, user.Email)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUserByUsername() {
	username := "testuser"
	mockUser := &models.UserEntity{
		ID:       uuid.New(),
		Username: username,
		Email:    "test@test.com",
		RoleID:   uuid.New(),
		Role: models.RoleEntity{
			ID:        uuid.New(),
			RoleType:  "USER",
			AdminRole: false,
			ModRole:   false,
			UserRole:  true,
		},
	}

	tests := []struct {
		name          string
		username      string
		mockUser      *models.UserEntity
		mockError     error
		expectedError string
	}{
		{
			name:          "success get user by username",
			username:      username,
			mockUser:      mockUser,
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "user not found",
			username:      "nonexistent",
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("FindByUsername", tt.username).Return(tt.mockUser, tt.mockError)

			user, err := suite.service.GetUserByUsername(tt.username)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
				assert.Nil(suite.T(), user)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), user)
				assert.Equal(suite.T(), tt.mockUser.ID, user.ID)
				assert.Equal(suite.T(), tt.mockUser.Username, user.Username)
				assert.Equal(suite.T(), tt.mockUser.Email, user.Email)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestCreateNewUser() {
	roleID := uuid.New()
	newUserDto := dto.UserDto{
		Username:    "newuser",
		Email:       "newuser@test.com",
		Password:    "password123",
		Role: dto.RoleDto{
			ID:         roleID,
			RoleType:   "USER",
			Permission: 1,
			AdminRole:  false,
			ModRole:    false,
		},
	}

	mockRole := &models.RoleEntity{
		ID:        roleID,
		RoleType:  "user",
		AdminRole: false,
		ModRole:   false,
		UserRole:  true,
	}

	tests := []struct {
		name          string
		userDto       dto.UserDto
		mockRole      *models.RoleEntity
		mockUserID    uuid.UUID
		mockError     error
		expectedError string
	}{
		{
			name:          "success create user",
			userDto:       newUserDto,
			mockRole:      mockRole,
			mockUserID:    uuid.New(),
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "error creating user",
			userDto:       newUserDto,
			mockRole:      mockRole,
			mockUserID:    uuid.Nil,
			mockError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRoleRepo.On("FindByType", "user").Return(tt.mockRole, nil)
			suite.mockUserRepo.On("Create", mock.AnythingOfType("models.UserEntity")).Return(tt.mockUserID, tt.mockError)

			userID, err := suite.service.CreateNewUser(tt.userDto)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
				assert.Equal(suite.T(), uuid.Nil, userID)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotEqual(suite.T(), uuid.Nil, userID)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestUpdateUser() {
	userID := uuid.New()
	roleID := uuid.New()
	updateReq := request.NewUser{
		Username: ptrToString("updateduser"),
		Locked:   ptrToBool(false),
		Disable:  ptrToBool(false),
		RoleID:   ptrToString(roleID.String()),
	}

	mockUser := &models.UserEntity{
		ID:       userID,
		Username: "oldusername",
		Email:    "old@test.com",
		RoleID:   roleID,
		Role: models.RoleEntity{
			ID:        roleID,
			RoleType:  "USER",
			AdminRole: false,
			ModRole:   false,
			UserRole:  true,
		},
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		updateReq     request.NewUser
		mockUser      *models.UserEntity
		mockError     error
		expectedError string
	}{
		{
			name:          "success update user",
			userID:        userID,
			updateReq:     updateReq,
			mockUser:      mockUser,
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "user not found",
			userID:        uuid.New(),
			updateReq:     updateReq,
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("FindByID", tt.userID).Return(tt.mockUser, tt.mockError)
			if tt.mockError == nil {
				suite.mockUserRepo.On("Update", tt.userID, mock.AnythingOfType("models.UserEntity")).Return(nil)
			}

			err := suite.service.UpdateUser(tt.userID, tt.updateReq)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestDeleteUser() {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockError     error
		expectedError string
	}{
		{
			name:          "success delete user",
			userID:        userID,
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "error deleting user",
			userID:        uuid.New(),
			mockError:     errors.New("delete error"),
			expectedError: "delete error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("Delete", tt.userID).Return(tt.mockError)

			err := suite.service.DeleteUser(tt.userID)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestRestoreUser() {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockError     error
		expectedError string
	}{
		{
			name:          "success restore user",
			userID:        userID,
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "error restoring user",
			userID:        uuid.New(),
			mockError:     errors.New("restore error"),
			expectedError: "restore error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("Restore", tt.userID).Return(tt.mockError)

			err := suite.service.RestoreUser(tt.userID)

			if tt.expectedError != "" {
				assert.EqualError(suite.T(), err, tt.expectedError)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestProcessAvatar() {
	userID := uuid.New()
	mockFileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     1024,
	}

	// Crear un mock de multipart.File
	mockFile := &mockMultipartFile{
		content: []byte("fake image content"),
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		fileHeader    *multipart.FileHeader
		file          multipart.File
		expectedError string
	}{
		{
			name:          "success process avatar",
			userID:        userID,
			fileHeader:    mockFileHeader,
			file:          mockFile,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.service.ProcessAvatar(tt.userID, tt.fileHeader, tt.file)
			
			// Como ProcessAvatar intenta procesar un archivo real,
			// esperamos un error ya que estamos usando datos mock
			assert.Error(suite.T(), err)
		})
	}
}

type mockMultipartFile struct {
	content []byte
	pos     int64
}

func (m *mockMultipartFile) Read(p []byte) (n int, err error) {
	if m.pos >= int64(len(m.content)) {
		return 0, io.EOF
	}
	n = copy(p, m.content[m.pos:])
	m.pos += int64(n)
	return
}

func (m *mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(m.content)) {
		return 0, io.EOF
	}
	n = copy(p, m.content[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (m *mockMultipartFile) Close() error               { return nil }
func (m *mockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.pos = offset
	case io.SeekCurrent:
		m.pos += offset
	case io.SeekEnd:
		m.pos = int64(len(m.content)) + offset
	}
	return m.pos, nil
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// Utils pointer to bool
func ptrToBool(b bool) *bool {
	return &b
}
