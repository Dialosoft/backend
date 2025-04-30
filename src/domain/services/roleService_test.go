package services

import (
	"errors"
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

type MockRoleRepository struct {
	repository.MockAbstractRepository[*models.RoleEntity, uuid.UUID]
}

func (m *MockRoleRepository) FindByType(roleType string) (*models.RoleEntity, error) {
	args := m.Called(roleType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoleEntity), args.Error(1)
}

type MockRolePermissionsRepository struct {
	repository.MockAbstractRepository[*models.RolePermissions, uuid.UUID]
}

func (m *MockRolePermissionsRepository) FindByRoleID(roleID uuid.UUID) (*models.RolePermissions, error) {
	args := m.Called(roleID)
	return args.Get(0).(*models.RolePermissions), args.Error(1)
}

type RoleServiceTestSuite struct {
	suite.Suite
	service      RoleService
	mockRoleRepo *MockRoleRepository
	mockPermRepo *MockRolePermissionsRepository
}

func (suite *RoleServiceTestSuite) SetupTest() {
	suite.mockRoleRepo = new(MockRoleRepository)
	suite.mockPermRepo = new(MockRolePermissionsRepository)
	suite.service = NewRoleService(suite.mockRoleRepo, suite.mockPermRepo)
}

func (suite *RoleServiceTestSuite) TearDownTest() {
	suite.mockRoleRepo.AssertExpectations(suite.T())
	suite.mockPermRepo.AssertExpectations(suite.T())
}

func (suite *RoleServiceTestSuite) TestCreateNewRole() {
	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name          string
		input         dto.RoleDto
		mockSetup     func()
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name: "success creating role with permissions",
			input: dto.RoleDto{
				RoleType:   "ADMIN",
				Permission: 1,
				AdminRole:  true,
				ModRole:    false,
			},
			mockSetup: func() {
	
				suite.mockRoleRepo.On("FindByType", "ADMIN").
					Return(nil, errors.New("role not found"))
				
				suite.mockRoleRepo.On("Create", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.RoleEntity")).
					Return(&models.RoleEntity{
						ID:         roleID,
						RoleType:   "ADMIN",
						Permission: 1,
						AdminRole:  true,
						ModRole:    false,
						UserRole:   false,
					}, nil)

				suite.mockPermRepo.On("Create", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.RolePermissions")).
					Return(&models.RolePermissions{
						RoleID:              roleID,
						CanManageCategories: true,
						CanManageForums:     true,
						CanManageRoles:      true,
						CanManageUsers:      true,
					}, nil)
			},
			expectedID:    roleID,
			expectedError: nil,
		},
		{
			name: "error when role type already exists",
			input: dto.RoleDto{
				RoleType:   "ADMIN",
				Permission: 1,
			},
			mockSetup: func() {
				suite.mockRoleRepo.On("FindByType", "ADMIN").
					Return(&models.RoleEntity{
						ID:       uuid.New(),
						RoleType: "ADMIN",
					}, nil)
			},
			expectedID:    uuid.Nil,
			expectedError: errors.New("role type already exists"),
		},
		{
			name: "error when creating role fails",
			input: dto.RoleDto{
				RoleType:   "ADMIN",
				Permission: 1,
			},
			mockSetup: func() {

				suite.mockRoleRepo.On("FindByType", "ADMIN").
					Return(nil, errors.New("role not found"))

				suite.mockRoleRepo.On("Create", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.RoleEntity")).
					Return((*models.RoleEntity)(nil), errors.New("database error"))
			},
			expectedID:    uuid.Nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset mock expectations
			suite.mockRoleRepo.ExpectedCalls = nil
			suite.mockPermRepo.ExpectedCalls = nil

			// Setup mocks for this test case
			tt.mockSetup()

			// Execute
			createdID, err := suite.service.CreateNewRole(tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorContains(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedID, createdID)
			}
		})
	}
}

func (suite *RoleServiceTestSuite) TestSetRolePermissionsByRoleID() {
	// We define a fixed UUID to use in all tests
	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name          string
		roleID        uuid.UUID
		input         request.NewRolePermissions
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success updating role permissions",
			roleID: roleID,
			input: request.NewRolePermissions{
				CanManageCategories: testUtils.PtrToBool(true),
				CanManageForums:     testUtils.PtrToBool(true),
				CanManageRoles:      testUtils.PtrToBool(false),
				CanManageUsers:      testUtils.PtrToBool(true),
			},
			mockSetup: func() {
				// Mock FindByRoleID to return existing permissions
				suite.mockPermRepo.On("FindByRoleID", roleID).
					Return(&models.RolePermissions{
						RoleID:              roleID,
						CanManageCategories: false,
						CanManageForums:     false,
						CanManageRoles:      true,
						CanManageUsers:      false,
					}, nil)

				// Mock Update to update permissions
				suite.mockPermRepo.On("Update", mock.AnythingOfType("*gorm.DB"), roleID, mock.AnythingOfType("*models.RolePermissions")).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error when role not found",
			roleID: roleID,
			input: request.NewRolePermissions{
				CanManageCategories: testUtils.PtrToBool(false),
				CanManageForums:     testUtils.PtrToBool(false),
				CanManageRoles:      testUtils.PtrToBool(false),
				CanManageUsers:      testUtils.PtrToBool(false),
			},
			mockSetup: func() {
				// Mock FindByRoleID to return error
				suite.mockPermRepo.On("FindByRoleID", roleID).
					Return((*models.RolePermissions)(nil), errors.New("role not found"))
			},
			expectedError: errors.New("role not found"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset mock expectations
			suite.mockPermRepo.ExpectedCalls = nil

			// Setup mocks for this test case
			tt.mockSetup()

			// Execute
			err := suite.service.SetRolePermissionsByRoleID(tt.roleID, tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorContains(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *RoleServiceTestSuite) TestGetRolePermissionsByRoleID() {
	// We define a fixed UUID to use in all tests
	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name           string
		roleID         uuid.UUID
		mockSetup      func()
		expectedResult *models.RolePermissions
		expectedError  error
	}{
		{
			name:   "success getting role permissions",
			roleID: roleID,
			mockSetup: func() {
				// Mock FindByRoleID to return existing permissions
				suite.mockPermRepo.On("FindByRoleID", roleID).
					Return(&models.RolePermissions{
						RoleID:              roleID,
						CanManageCategories: true,
						CanManageForums:     true,
						CanManageRoles:      false,
						CanManageUsers:      true,
					}, nil)
			},
			expectedResult: &models.RolePermissions{
				RoleID:              roleID,
				CanManageCategories: true,
				CanManageForums:     true,
				CanManageRoles:      false,
				CanManageUsers:      true,
			},
			expectedError: nil,
		},
		{
			name:   "error when role not found",
			roleID: roleID,
			mockSetup: func() {
				// Mock FindByRoleID to return error
				suite.mockPermRepo.On("FindByRoleID", roleID).
					Return((*models.RolePermissions)(nil), errors.New("role not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("role not found"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset mock expectations
			suite.mockPermRepo.ExpectedCalls = nil

			// Setup mocks for this test case
			tt.mockSetup()

			// Execute
			result, err := suite.service.GetRolePermissionsByRoleID(tt.roleID)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorContains(suite.T(), err, tt.expectedError.Error())
				assert.Nil(suite.T(), result)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			}
		})
	}
}

func (suite *RoleServiceTestSuite) TestRestoreRole() {
	// We define a fixed UUID to use in all tests
	roleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name          string
		roleID        uuid.UUID
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success restoring role",
			roleID: roleID,
			mockSetup: func() {
				// Mock Restore to return success
				suite.mockRoleRepo.On("Restore", mock.AnythingOfType("*gorm.DB"), roleID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error when role not found or already restored",
			roleID: roleID,
			mockSetup: func() {
				// Mock Restore to return error
				suite.mockRoleRepo.On("Restore", mock.AnythingOfType("*gorm.DB"), roleID).
					Return(errors.New("role not found or already restored"))
			},
			expectedError: errors.New("role not found or already restored"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset mock expectations
			suite.mockRoleRepo.ExpectedCalls = nil

			// Setup mocks for this test case
			tt.mockSetup()

			// Execute
			err := suite.service.RestoreRole(tt.roleID)

			// Assert
			if tt.expectedError != nil {
				assert.ErrorContains(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *RoleServiceTestSuite) TestGetDefaultRoles() {
	// We define fixed UUIDs to use in all tests
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	moderatorID := uuid.MustParse("650e8400-e29b-41d4-a716-446655440000")
	adminID := uuid.MustParse("750e8400-e29b-41d4-a716-446655440000")
	anonymousID := uuid.MustParse("850e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name           string
		mockSetup      func()
		expectedResult map[string]uuid.UUID
		expectedError  error
	}{
		{
			name: "success getting all default roles",
			mockSetup: func() {
				// Mock FindByType for each role type
				suite.mockRoleRepo.On("FindByType", "user").
					Return(&models.RoleEntity{
						ID:       userID,
						RoleType: "user",
					}, nil)
				suite.mockRoleRepo.On("FindByType", "moderator").
					Return(&models.RoleEntity{
						ID:       moderatorID,
						RoleType: "moderator",
					}, nil)
				suite.mockRoleRepo.On("FindByType", "administrator").
					Return(&models.RoleEntity{
						ID:       adminID,
						RoleType: "administrator",
					}, nil)
				suite.mockRoleRepo.On("FindByType", "anonymous").
					Return(&models.RoleEntity{
						ID:       anonymousID,
						RoleType: "anonymous",
					}, nil)
			},
			expectedResult: map[string]uuid.UUID{
				"user":          userID,
				"moderator":     moderatorID,
				"administrator": adminID,
				"anonymous":     anonymousID,
			},
			expectedError: nil,
		},
		{
			name: "error when moderator role not found",
			mockSetup: func() {
				// Mock FindByType for user role
				suite.mockRoleRepo.On("FindByType", "user").
					Return(&models.RoleEntity{
						ID:       userID,
						RoleType: "user",
					}, nil)
				// Mock FindByType for anonymous role
				suite.mockRoleRepo.On("FindByType", "anonymous").
					Return(&models.RoleEntity{
						ID:       anonymousID,
						RoleType: "anonymous",
					}, nil)
				// Mock FindByType for moderator role to return error
				suite.mockRoleRepo.On("FindByType", "moderator").
					Return((*models.RoleEntity)(nil), errors.New("role not found"))
				// Mock for admin role should not be called
			},
			expectedResult: nil,
			expectedError:  errors.New("failed to get default role moderator: role not found"),
		},
		{
			name: "error when getting administrator role fails",
			mockSetup: func() {
				// Mock FindByType for user role
				suite.mockRoleRepo.On("FindByType", "user").
					Return(&models.RoleEntity{
						ID:       userID,
						RoleType: "user",
					}, nil)
				// Mock FindByType for anonymous role
				suite.mockRoleRepo.On("FindByType", "anonymous").
					Return(&models.RoleEntity{
						ID:       anonymousID,
						RoleType: "anonymous",
					}, nil)
				// Mock FindByType for moderator role
				suite.mockRoleRepo.On("FindByType", "moderator").
					Return(&models.RoleEntity{
						ID:       moderatorID,
						RoleType: "moderator",
					}, nil)
				// Mock FindByType for admin role to return error
				suite.mockRoleRepo.On("FindByType", "administrator").
					Return((*models.RoleEntity)(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("failed to get default role administrator: database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset mock expectations
			suite.mockRoleRepo.ExpectedCalls = nil

			// Setup mocks for this test case
			tt.mockSetup()

			// Execute
			result, err := suite.service.GetDefaultRoles()

			// Assert
			if tt.expectedError != nil {
				assert.ErrorContains(suite.T(), err, tt.expectedError.Error())
				assert.Nil(suite.T(), result)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			}
		})
	}
}

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}
