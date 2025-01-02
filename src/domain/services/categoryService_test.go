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
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
)

type MockCategoryRepository struct {
	repository.MockAbstractRepository[*models.Category, uuid.UUID]
}

func (m *MockCategoryRepository) FindByName(name string) (*models.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAllIncludingDeleted() ([]*models.Category, error) {
	args := m.Called()
	return args.Get(0).([]*models.Category), args.Error(1)
}


// ****************************************
// ****************************************

type MockRoleRepository struct {
	repository.MockAbstractRepository[*models.RoleEntity, uuid.UUID]
}

func (m *MockRoleRepository) FindByType(roleType string) (*models.RoleEntity, error) {
	args := m.Called(roleType)
	return args.Get(0).(*models.RoleEntity), args.Error(1)
}


type CategoryServiceTestSuite struct {
	suite.Suite
	mockCategoryRepo *MockCategoryRepository
	mockRoleRepo     *MockRoleRepository
	service          CategoryService
}

func (suite *CategoryServiceTestSuite) SetupTest() {
	suite.mockCategoryRepo = new(MockCategoryRepository)
	suite.mockRoleRepo = new(MockRoleRepository)
	suite.service = NewCategoryService(suite.mockCategoryRepo, suite.mockRoleRepo)
}

func (suite *CategoryServiceTestSuite) TearDownTest() {
	suite.mockCategoryRepo.AssertExpectations(suite.T())
	suite.mockRoleRepo.AssertExpectations(suite.T())
}

func (suite *CategoryServiceTestSuite) TestGetAllCategories() {
	// We define a fixed UUID to use in all tests
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		mockReturn     []*models.Category
		mockReturnErr  error
		expectedResult []response.CategoryResponse
		expectedErr    error
	}{
		{
			name: "success retrieving all categories",
			mockReturn: []*models.Category{
				{
					ID:          categoryID,
					Name:        "Electronics",
					Description: "Devices and gadgets",
				},
			},
			mockReturnErr: nil,
			expectedResult: []response.CategoryResponse{
				mapper.CategoryEntityToCategoryResponse(&models.Category{
					ID:          categoryID,
					Name:        "Electronics",
					Description: "Devices and gadgets",
				}),
			},
			expectedErr: nil,
		},
		{
			name:           "error retrieving categories",
			mockReturn:     nil,
			mockReturnErr:  errors.New("database error"),
			expectedResult: nil,
			expectedErr:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("FindAll").Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetAllCategories()

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestGetCategoryByID() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		id             uuid.UUID
		mockReturn     *models.Category
		mockReturnErr  error
		expectedResult *dto.CategoryDto
		expectedErr    error
	}{
		{
			name: "success retrieving category by ID",
			id:   categoryID,
			mockReturn: &models.Category{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			},
			mockReturnErr: nil,
			expectedResult: mapper.CategoryEntityToCategoryDto(&models.Category{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			}),
			expectedErr: nil,
		},
		{
			name:           "category not found",
			id:             categoryID,
			mockReturn:     nil,
			mockReturnErr:  errorsUtils.ErrNotFound,
			expectedResult: nil,
			expectedErr:    errorsUtils.ErrNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("FindByID", tt.id).Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetCategoryByID(tt.id)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestGetCategoryByName() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		categoryName   string
		mockReturn     *models.Category
		mockReturnErr  error
		expectedResult *dto.CategoryDto
		expectedErr    error
	}{
		{
			name:         "success retrieving category by name",
			categoryName: "Electronics",
			mockReturn: &models.Category{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			},
			mockReturnErr: nil,
			expectedResult: &dto.CategoryDto{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			},
			expectedErr: nil,
		},
		{
			name:           "category not found",
			categoryName:   "NonExistentCategory",
			mockReturn:     nil,
			mockReturnErr:  errorsUtils.ErrNotFound,
			expectedResult: nil,
			expectedErr:    errorsUtils.ErrNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("FindByName", tt.categoryName).Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetCategoryByName(tt.categoryName)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestGetAllCategoriesAllowedByRole() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		roleID         string
		mockReturn     []*models.Category
		mockReturnErr  error
		expectedResult []response.CategoryResponse
		expectedErr    error
	}{
		{
			name:   "success retrieving categories allowed by role",
			roleID: "administrator",
			mockReturn: []*models.Category{
				{
					ID:           categoryID,
					Name:         "Electronics",
					Description:  "Devices and gadgets",
					RolesAllowed: []string{"administrator"},
				},
				{
					ID:           categoryID,
					Name:         "Books",
					Description:  "Books and literature",
					RolesAllowed: []string{"administrator", "user"},
				},
			},
			mockReturnErr: nil,
			expectedResult: []response.CategoryResponse{
				mapper.CategoryEntityToCategoryResponse(&models.Category{
					ID:           categoryID,
					Name:         "Electronics",
					Description:  "Devices and gadgets",
					RolesAllowed: []string{"administrator"},
				}),
				mapper.CategoryEntityToCategoryResponse(&models.Category{
					ID:           categoryID,
					Name:         "Books",
					Description:  "Books and literature",
					RolesAllowed: []string{"administrator", "user"},
				}),
			},
			expectedErr: nil,
		},
		{
			name:           "no categories found for role",
			roleID:         "guest",
			mockReturn:     []*models.Category{},
			mockReturnErr:  nil,
			expectedResult: []response.CategoryResponse{},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("FindAll").Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetAllCategoriesAllowedByRole(tt.roleID)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestCreateCategory() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")
	roleID := uuid.MustParse("22ceae31-e0be-44f6-b699-4740500e8acc")

	tests := []struct {
		name             string
		newCategory      request.NewCategory
		mockRoles        []*models.RoleEntity
		mockRolesErr     error
		mockCreateErr    error
		mockCreateResult *models.Category
		expectedResult   uuid.UUID
		expectedErr      error
	}{
		{
			name: "success creating category",
			newCategory: request.NewCategory{
				Name:           ptrToString("Electronics"),
				Description:    ptrToString("Devices and gadgets"),
				RolesAllowedID: []string{roleID.String()},
			},
			mockRoles: []*models.RoleEntity{
				{ID: roleID, RoleType: "administrator"},
			},
			mockRolesErr:  nil,
			mockCreateErr: nil,
			mockCreateResult: &models.Category{
				ID:           categoryID,
				Name:         "Electronics",
				Description:  "Devices and gadgets",
				RolesAllowed: []string{roleID.String()},
			},
			expectedResult: categoryID,
			expectedErr:    nil,
		},
		{
			name: "invalid role UUID provided",
			newCategory: request.NewCategory{
				Name:           ptrToString("Electronics"),
				Description:    ptrToString("Devices and gadgets"),
				RolesAllowedID: []string{"invalid-uuid"},
			},
			mockRolesErr:     nil,
			mockCreateErr:    nil,
			mockCreateResult: nil,
			expectedResult:   uuid.UUID{},
			expectedErr:      errorsUtils.ErrInvalidUUID,
		},
		{
			name: "role not found in system",
			newCategory: request.NewCategory{
				Name:           ptrToString("Electronics"),
				Description:    ptrToString("Devices and gadgets"),
				RolesAllowedID: []string{roleID.String()},
			},
			mockRoles:        []*models.RoleEntity{}, // No roles found
			mockRolesErr:     nil,
			mockCreateErr:    nil,
			mockCreateResult: nil,
			expectedResult:   uuid.UUID{},
			expectedErr:      errorsUtils.ErrNotFound,
		},
		{
			name: "error during category creation in repository",
			newCategory: request.NewCategory{
				Name:           ptrToString("Electronics 2"),
				Description:    ptrToString("Devices and gadgets"),
				RolesAllowedID: []string{roleID.String()},
			},
			mockRoles:        []*models.RoleEntity{{ID: roleID, RoleType: "administrator"}},
			mockRolesErr:     nil,
			mockCreateErr:    errors.New("database error"),
			mockCreateResult: nil,
			expectedResult:   uuid.UUID{},
			expectedErr:      errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockRoleRepo.ExpectedCalls = nil
			suite.mockRoleRepo.On("FindAll").Return(tt.mockRoles, tt.mockRolesErr)

			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			// Configure the Mock correctly to return the category and error
			if tt.expectedErr == nil {
				suite.mockCategoryRepo.On("Create",
					mock.AnythingOfType("*gorm.DB"),
					mock.MatchedBy(func(category *models.Category) bool {
						return category.Name == *tt.newCategory.Name &&
							category.Description == *tt.newCategory.Description &&
							len(category.RolesAllowed) == len(tt.newCategory.RolesAllowedID) &&
							category.RolesAllowed[0] == tt.newCategory.RolesAllowedID[0]
					}),
				).Return(tt.mockCreateResult, tt.mockCreateErr)
			} else {
				suite.mockCategoryRepo.On("Create",
					mock.AnythingOfType("*gorm.DB"),
					mock.AnythingOfType("*models.Category"),
				).Return(tt.mockCreateResult, tt.mockCreateErr)
			}

			result, err := suite.service.CreateCategory(tt.newCategory)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Equal(suite.T(), tt.expectedResult, result)
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestUpdateCategory() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		categoryID     uuid.UUID
		updateReq      request.NewCategory
		mockFindReturn *models.Category
		mockFindErr    error
		mockUpdateErr  error
		expectedErr    error
	}{
		{
			name:       "success updating category",
			categoryID: categoryID,
			updateReq: request.NewCategory{
				Name:        ptrToString("Updated Electronics"),
				Description: ptrToString("Updated description"),
			},
			mockFindReturn: &models.Category{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			},
			mockFindErr:   nil,
			mockUpdateErr: nil,
			expectedErr:   nil,
		},
		{
			name:       "category not found",
			categoryID: categoryID,
			updateReq: request.NewCategory{
				Name:        ptrToString("Updated Electronics"),
				Description: ptrToString("Updated description"),
			},
			mockFindReturn: nil,
			mockFindErr:    errorsUtils.ErrNotFound,
			mockUpdateErr:  nil,
			expectedErr:    errorsUtils.ErrNotFound,
		},
		{
			name:       "error during update",
			categoryID: categoryID,
			updateReq: request.NewCategory{
				Name:        ptrToString("Updated Electronics"),
				Description: ptrToString("Updated description"),
			},
			mockFindReturn: &models.Category{
				ID:          categoryID,
				Name:        "Electronics",
				Description: "Devices and gadgets",
			},
			mockFindErr:   nil,
			mockUpdateErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("FindByID", tt.categoryID).Return(tt.mockFindReturn, tt.mockFindErr)

			if tt.mockFindReturn != nil {
				suite.mockCategoryRepo.On("Update",
					mock.AnythingOfType("*gorm.DB"),
					tt.categoryID,
					mock.MatchedBy(func(category *models.Category) bool {
						return category.Name == *tt.updateReq.Name &&
							category.Description == *tt.updateReq.Description &&
							category.ID == tt.categoryID
					}),
				).Return(tt.mockUpdateErr)
			}

			err := suite.service.UpdateCategory(tt.categoryID, tt.updateReq)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestDeleteCategory() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name          string
		categoryID    uuid.UUID
		mockDeleteErr error
		expectedErr   error
	}{
		{
			name:          "success deleting category",
			categoryID:    categoryID,
			mockDeleteErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "error during deletion",
			categoryID:    categoryID,
			mockDeleteErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("Delete",
				mock.AnythingOfType("*gorm.DB"),
				tt.categoryID,
			).Return(tt.mockDeleteErr)

			err := suite.service.DeleteCategory(tt.categoryID)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
			}
		})
	}
}

func (suite *CategoryServiceTestSuite) TestRestoreCategory() {
	categoryID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name          string
		categoryID    uuid.UUID
		mockReturnErr error
		expectedErr   error
	}{
		{
			name:          "success restoring category",
			categoryID:    categoryID,
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "error restoring category",
			categoryID:    categoryID,
			mockReturnErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockCategoryRepo.ExpectedCalls = nil

			suite.mockCategoryRepo.On("Restore",
				mock.AnythingOfType("*gorm.DB"),
				tt.categoryID,
			).Return(tt.mockReturnErr)

			err := suite.service.RestoreCategory(tt.categoryID)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
			}
		})
	}
}

// Create a new instance of the CategoryServiceTestSuite structure and executes the test suite
func TestCategoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}

// Utils pointer to string
func ptrToString(s string) *string {
	return &s
}
