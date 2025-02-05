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

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}
