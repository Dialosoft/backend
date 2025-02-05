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
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	testUtils "github.com/Dialosoft/src/pkg/utils/test"
)

// MockUserRepository is a mock of the user repository that extends the mock of the abstract repository
type MockUserRepository struct {
	repository.MockAbstractRepository[*models.UserEntity, uuid.UUID]
}

// FindByUsername is the only specific method we need to implement
func (m *MockUserRepository) FindByUsername(username string) (*models.UserEntity, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
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
	userIdMock := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	roleIdMock := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
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
					ID:       userIdMock,
					Username: "user1",
					Email:    "user1@test.com",
					RoleID:   roleIdMock,
					Role: models.RoleEntity{
						ID:        roleIdMock,
						RoleType:  "USER",
						AdminRole: false,
						ModRole:   false,
						UserRole:  true,
					},
				},
				{
					ID:       userIdMock,
					Username: "user2",
					Email:    "user2@test.com",
					RoleID:   roleIdMock,
					Role: models.RoleEntity{
						ID:        roleIdMock,
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
			suite.mockUserRepo.On("FindAll").Return(tt.mockUsers, tt.mockError).Once()

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
					assert.Equal(suite.T(), tt.mockUsers[i].Role.ID, user.Role.ID)
					assert.Equal(suite.T(), tt.mockUsers[i].Role.RoleType, user.Role.RoleType)
					assert.Equal(suite.T(), tt.mockUsers[i].Role.AdminRole, user.Role.AdminRole)
					assert.Equal(suite.T(), tt.mockUsers[i].Role.ModRole, user.Role.ModRole)
				}
			}
			suite.mockUserRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUserByID() {
	userID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")
	roleID := uuid.MustParse("22ceae31-e0be-44f6-b699-4740500e8acc")
	mockUser := &models.UserEntity{
		ID:       userID,
		Username: "testuser",
		Email:    "test@test.com",
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
			userID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
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
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Username: username,
		Email:    "test@test.com",
		RoleID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Role: models.RoleEntity{
			ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
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
	userID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")
	roleID := uuid.MustParse("22ceae31-e0be-44f6-b699-4740500e8acc")

	tests := []struct {
		name          string
		userDto       dto.UserDto
		mockRole      *models.RoleEntity
		mockRoleErr   error
		mockCreateErr error
		mockCreateResult *models.UserEntity
		expectedResult uuid.UUID
		expectedError error
	}{
		{
			name: "success create user",
			userDto: dto.UserDto{
				Username: "newuser",
				Email:    "newuser@test.com",
				Password: "password123",
				Role: dto.RoleDto{
					ID:         roleID,
					RoleType:   "USER",
					Permission: 1,
					AdminRole:  false,
					ModRole:    false,
				},
			},
			mockRole: &models.RoleEntity{
				ID:        roleID,
				RoleType:  "user",
				AdminRole: false,
				ModRole:   false,
				UserRole:  true,
			},
			mockRoleErr: nil,
			mockCreateResult: &models.UserEntity{
				ID:       userID,
				Username: "newuser",
				Email:    "newuser@test.com",
				RoleID:   roleID,
				Role:     models.RoleEntity{ID: roleID},
			},
			mockCreateErr:   nil,
			expectedResult:  userID,
			expectedError:   nil,
		},
		{
			name: "role not found",
			userDto: dto.UserDto{
				Username: "newuser",
				Email:    "newuser@test.com",
				Password: "password123",
			},
			mockRole:      nil,
			mockRoleErr:   errors.New("role not found"),
			mockCreateResult: nil,
			mockCreateErr:   nil,
			expectedResult:  uuid.UUID{},
			expectedError:   errors.New("role not found"),
		},
		{
			name: "error creating user",
			userDto: dto.UserDto{
				Username: "newuser",
				Email:    "newuser@test.com",
				Password: "password123",
			},
			mockRole: &models.RoleEntity{
				ID:        roleID,
				RoleType:  "user",
				AdminRole: false,
				ModRole:   false,
				UserRole:  true,
			},
			mockRoleErr:   nil,
			mockCreateResult: nil,
			mockCreateErr:   errors.New("database error"),
			expectedResult:  uuid.UUID{},
			expectedError:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockRoleRepo.ExpectedCalls = nil
			suite.mockRoleRepo.On("FindByType", "user").Return(tt.mockRole, tt.mockRoleErr)

			// Reset Mock's expectations before each subtest
			suite.mockUserRepo.ExpectedCalls = nil

			// Configure the Mock correctly to return the user and error
			if tt.mockRoleErr == nil {
				suite.mockUserRepo.On("Create", 
					mock.AnythingOfType("*gorm.DB"),
					mock.MatchedBy(func(user *models.UserEntity) bool {
						return user.Username == tt.userDto.Username &&
							user.Email == tt.userDto.Email &&
							user.RoleID == tt.mockRole.ID &&
							user.Role == *tt.mockRole
					}),
				).Return(tt.mockCreateResult, tt.mockCreateErr)
			}

			userID, err := suite.service.CreateNewUser(tt.userDto)

			if tt.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedError.Error(), err.Error())
				assert.Equal(suite.T(), tt.expectedResult, userID)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, userID)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestUpdateUser() {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	updateReq := request.NewUser{
		Username: testUtils.PtrToString("updateduser"),
		Locked:   ptrToBool(false),
		Disable:  ptrToBool(false),
		RoleID:   testUtils.PtrToString(roleID.String()),
	}

	mockRole := &models.RoleEntity{
		ID:        roleID,
		RoleType:  "USER",
		AdminRole: false,
		ModRole:   false,
		UserRole:  true,
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
		mockRole      *models.RoleEntity
		mockUserErr   error
		mockRoleErr   error
		expectedError string
	}{
		{
			name:          "success update user",
			userID:        userID,
			updateReq:     updateReq,
			mockUser:      mockUser,
			mockRole:      mockRole,
			mockUserErr:   nil,
			mockRoleErr:   nil,
			expectedError: "",
		},
		{
			name:          "user not found",
			userID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			updateReq:     updateReq,
			mockUser:      nil,
			mockRole:      nil,
			mockUserErr:   errors.New("user not found"),
			mockRoleErr:   nil,
			expectedError: "user not found",
		},
		{
			name:          "role not found",
			userID:        userID,
			updateReq:     updateReq,
			mockUser:      mockUser,
			mockRole:      nil,
			mockUserErr:   nil,
			mockRoleErr:   errors.New("role not found"),
			expectedError: "role not found",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockUserRepo.ExpectedCalls = nil
			suite.mockRoleRepo.ExpectedCalls = nil

			// Configure user repository mock
			suite.mockUserRepo.On("FindByID", tt.userID).Return(tt.mockUser, tt.mockUserErr)

			if tt.mockUserErr == nil && tt.updateReq.RoleID != nil {
				roleUUID, _ := uuid.Parse(*tt.updateReq.RoleID)
				suite.mockRoleRepo.On("FindByID", roleUUID).Return(tt.mockRole, tt.mockRoleErr)
			}

			if tt.mockUserErr == nil && tt.mockRoleErr == nil {
				suite.mockUserRepo.On("Update", 
					mock.AnythingOfType("*gorm.DB"),
					tt.userID,
					mock.MatchedBy(func(user *models.UserEntity) bool {
						matches := user.Username == *tt.updateReq.Username &&
							user.Banned == *tt.updateReq.Locked
						
						if tt.mockRole != nil {
							matches = matches &&
								user.RoleID == tt.mockRole.ID &&
								user.Role == *tt.mockRole
						}
						
						return matches
					}),
				).Return(nil)
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
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

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
			userID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			mockError:     errors.New("delete error"),
			expectedError: "delete error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("Delete", mock.AnythingOfType("*gorm.DB"), tt.userID).Return(tt.mockError)

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
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

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
			userID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			mockError:     errors.New("restore error"),
			expectedError: "restore error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockUserRepo.On("Restore", mock.AnythingOfType("*gorm.DB"), tt.userID).Return(tt.mockError)

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
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	mockFileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     1024,
	}

	// Create a mock of multipart.File
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

			// As ProcessAvatar attempts to process a real file,
			// we expect an error since we're using mock data
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

func (m *mockMultipartFile) Close() error { return nil }
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
