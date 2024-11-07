package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/Dialosoft/src/services"

)

// Mock del CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindAll() ([]*models.Category, error) {
	args := m.Called()
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByName(name string) (*models.Category, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Create(category models.Category) (uuid.UUID, error) {
	args := m.Called(category)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockCategoryRepository) Update(category models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) Restore(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// RoleRepository Mock
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindAllRoles() ([]*models.RoleEntity, error) {
	args := m.Called()
	return args.Get(0).([]*models.RoleEntity), args.Error(1)
}

// Test GetAllCategories
func TestCategoryService_GetAllCategories(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	t.Run("success", func(t *testing.T) {
		mockCategories := []*models.Category{
			{Name: "Test Category 1"},
			{Name: "Test Category 2"},
		}
		mockRepo.On("FindAll").Return(mockCategories, nil)

		result, err := service.GetAllCategories()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("FindAll").Return(nil, errors.New("database error"))

		result, err := service.GetAllCategories()

		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

// Test GetCategoryByID
func TestCategoryService_GetCategoryByID(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	t.Run("success", func(t *testing.T) {
		categoryID := uuid.New()
		mockCategory := &models.Category{ID: categoryID, Name: "Test Category"}
		
		mockRepo.On("FindByID", categoryID).Return(mockCategory, nil)

		result, err := service.GetCategoryByID(categoryID)

		assert.NoError(t, err)
		assert.Equal(t, "Test Category", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		categoryID := uuid.New()
		
		mockRepo.On("FindByID", categoryID).Return(nil, errorsUtils.ErrNotFound)

		result, err := service.GetCategoryByID(categoryID)

		assert.Nil(t, result)
		assert.EqualError(t, err, errorsUtils.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

// Test CreateCategory
func TestCategoryService_CreateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := services.NewCategoryService(mockRepo, mockRoleRepo)

	t.Run("success", func(t *testing.T) {
		newCatReq := request.NewCategory{
			Name:           strPtr("New Category"),
			Description:    strPtr("Description"),
			RolesAllowedID: []string{},
		}
		
        roleEntities := []*models.RoleEntity{
			{ID: uuid.New(), RoleType: "administrator"},
			{ID: uuid.New(), RoleType: "user"},
        }

        mockRoleRepo.On("FindAllRoles").Return(roleEntities, nil)

        newCatEntity := models.Category{
            Name:         "New Category",
            Description:  "Description",
            RolesAllowed: []string{},
        }

        mockRepo.On("Create", newCatEntity).Return(uuid.New(), nil)

        id, err := service.CreateCategory(newCatReq)

        assert.NoError(t, err)
        assert.NotEqual(t, uuid.Nil, id)

        mockRoleRepo.AssertExpectations(t)
        mockRepo.AssertExpectations(t)
    })

    t.Run("invalid role UUID", func(t *testing.T) {
        newCatReq := request.NewCategory{
            Name:           strPtr("New Category"),
            Description:    strPtr("Description"),
            RolesAllowedID: []string{"invalid-uuid"},
        }
        
        roleEntities := []*models.RoleEntity{
			{ID: uuid.New(), RoleType: "administrator"},
			{ID: uuid.New(), RoleType: "user"},
        }

        mockRoleRepo.On("FindAllRoles").Return(roleEntities, nil)

        _, err := service.CreateCategory(newCatReq)

        assert.EqualError(t, err, errorsUtils.ErrInvalidUUID.Error())
    })
}


// Test GetCategoryByName
func TestCategoryService_GetCategoryByName(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	tests := []struct {
		name        string
		categoryName string
		mockReturn  *models.Category
		mockError   error
		expectedErr error
	}{
		{
			name:        "success",
			categoryName: "Test Category",
			mockReturn:  &models.Category{Name: "Test Category"},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "category not found",
			categoryName: "Nonexistent Category",
			mockReturn:  nil,
			mockError:   errorsUtils.ErrNotFound,
			expectedErr: errorsUtils.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("FindByName", tt.categoryName).Return(tt.mockReturn, tt.mockError)

			result, err := service.GetCategoryByName(tt.categoryName)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn.Name, result.Name)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Test GetAllCategoriesAllowedByRole
func TestCategoryService_GetAllCategoriesAllowedByRole(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	tests := []struct {
		name         string
		roleID       string
		mockCategories []*models.Category
		expectedCount int
		expectedErr  error
	}{
		{
			name:         "success with matching roles",
			roleID:       "administrator",
			mockCategories: []*models.Category{
				{Name: "Public Category", RolesAllowed: []string{}},
				{Name: "Admin Category", RolesAllowed: []string{"administrator"}},
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
        {
            name:         "no matching roles",
            roleID:       "user",
            mockCategories: []*models.Category{
                {Name: "Admin Category", RolesAllowed: []string{"administrator"}},
            },
            expectedCount: 0,
            expectedErr:   nil,
        },
        {
            name:         "error retrieving categories",
            roleID:       "administrator",
            mockCategories: nil,
            expectedCount: 0,
            expectedErr:   errors.New("database error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo.On("FindAll").Return(tt.mockCategories, tt.expectedErr)

            result, err := service.GetAllCategoriesAllowedByRole(tt.roleID)

            if tt.expectedErr != nil {
                assert.EqualError(t, err, tt.expectedErr.Error())
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.Len(t, result, tt.expectedCount)
            }
            mockRepo.AssertExpectations(t)
        })
    }
}

// Test UpdateCategory
func TestCategoryService_UpdateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	t.Run("success", func(t *testing.T) {
        categoryID := uuid.New()
        existingCategory := &models.Category{ID: categoryID, Name: "Old Name", Description: "Old Description"}
        updatedRequest := request.NewCategory{Name: strPtr("New Name"), Description: strPtr("New Description")}

        mockRepo.On("FindByID", categoryID).Return(existingCategory, nil)
        updatedEntity := models.Category{ID: categoryID, Name: "New Name", Description: "New Description"}
        mockRepo.On("Update", updatedEntity).Return(nil)

        err := service.UpdateCategory(categoryID, updatedRequest)

        assert.NoError(t, err)
        mockRepo.AssertExpectations(t)
    })

    t.Run("category not found", func(t *testing.T) {
        categoryID := uuid.New()
        updatedRequest := request.NewCategory{Name: strPtr("New Name"), Description: strPtr("New Description")}

        mockRepo.On("FindByID", categoryID).Return(nil, errorsUtils.ErrNotFound)

        err := service.UpdateCategory(categoryID, updatedRequest)

        assert.EqualError(t, err, errorsUtils.ErrNotFound.Error())
    })

    t.Run("error updating category", func(t *testing.T) {
        categoryID := uuid.New()
        existingCategory := &models.Category{ID: categoryID}
        updatedRequest := request.NewCategory{Name: strPtr("New Name")}

        mockRepo.On("FindByID", categoryID).Return(existingCategory, nil)
        mockRepo.On("Update", mock.Anything).Return(errors.New("update error"))

        err := service.UpdateCategory(categoryID, updatedRequest)

        assert.EqualError(t, err, "update error")
    })
}


// Test DeleteCategory
func TestCategoryService_DeleteCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	tests := []struct {
		name        string
		categoryID  uuid.UUID
		mockError   error
		expectedErr error
	}{
		{
			name:        "success",
			categoryID:  uuid.New(),
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "category not found",
			categoryID:  uuid.New(),
			mockError:   errorsUtils.ErrNotFound,
			expectedErr: errorsUtils.ErrNotFound,
		},
		{
			name:        "internal server error",
			categoryID:  uuid.New(),
			mockError:   errors.New("internal server error"),
			expectedErr: errors.New("internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("Delete", tt.categoryID).Return(tt.mockError)

			err := service.DeleteCategory(tt.categoryID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Test RestoreCategory
func TestCategoryService_RestoreCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := services.NewCategoryService(mockRepo, nil)

	tests := []struct {
		name        string
		categoryID  uuid.UUID
		mockError   error
		expectedErr error
	}{
		{
			name:        "success",
			categoryID:  uuid.New(),
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "category not found",
			categoryID:  uuid.New(),
			mockError:   errorsUtils.ErrNotFound,
			expectedErr: errorsUtils.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("Restore", tt.categoryID).Return(tt.mockError)

			err := service.RestoreCategory(tt.categoryID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create a string pointer
func strPtr(s string) *string {
	return &s
}