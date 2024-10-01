package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/google/uuid"
)

// CategoryService defines the methods for managing categories in the system.
type CategoryService interface {
	// GetAllCategories retrieves a list of all available categories.
	// Returns a slice of CategoryDto or an error if something goes wrong.
	GetAllCategories() ([]response.CategoryResponse, error)

	// GetCategoryByID retrieves a specific category by its unique ID.
	// Returns the CategoryDto or an error if the category is not found.
	GetCategoryByID(id uuid.UUID) (*dto.CategoryDto, error)

	// GetCategoryByName retrieves a specific category by its name.
	// Returns the CategoryDto or an error if the category is not found.
	GetCategoryByName(name string) (*dto.CategoryDto, error)

	GetAllCategoriesAllowedByRole(roleID string) ([]response.CategoryResponse, error)

	// CreateCategory adds a new category based on the provided newCategory request.
	// Returns the UUID of the newly created category or an error if creation fails.
	CreateCategory(newCategory request.NewCategory) (uuid.UUID, error)

	// UpdateCategory updates an existing category's information by its ID.
	// The updated data is provided via the NewCategory request structure.
	// Returns an error if the update fails or the category is not found.
	UpdateCategory(id uuid.UUID, req request.NewCategory) error

	// DeleteCategory removes a category by its ID.
	// Returns an error if the deletion fails or the category is not found.
	DeleteCategory(id uuid.UUID) error

	// RestoreCategory restores a previously deleted category by its ID.
	// Returns an error if the restoration fails or the category is not found.
	RestoreCategory(id uuid.UUID) error
}

type categoryServiceImpl struct {
	categoryRepository repository.CategoryRepository
	roleRepository     repository.RoleRepository
}

// CreateCategory implements CategoryService.
func (service *categoryServiceImpl) CreateCategory(newCategory request.NewCategory) (uuid.UUID, error) {
	empty := ""
	rolesMap := make(map[uuid.UUID]string)

	if newCategory.Description == nil {
		newCategory.Description = &empty
	}

	roleEntities, err := service.roleRepository.FindAllRoles()
	if err != nil {
		return uuid.UUID{}, err
	}

	for _, roleEntity := range roleEntities {
		rolesMap[roleEntity.ID] = roleEntity.RoleType
	}

	for _, roleID := range newCategory.RolesAllowedID {
		roleUUID, err := uuid.Parse(roleID)
		if err != nil {
			return uuid.UUID{}, errorsUtils.ErrInvalidUUID
		}

		if _, ok := rolesMap[roleUUID]; !ok {
			return uuid.UUID{}, errorsUtils.ErrNotFound
		}
	}

	newCategoryEntity := models.Category{
		Name:         *newCategory.Name,
		Description:  *newCategory.Description,
		RolesAllowed: newCategory.RolesAllowedID,
	}

	id, err := service.categoryRepository.Create(newCategoryEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

// DeleteCategory implements CategoryService.
func (service *categoryServiceImpl) DeleteCategory(id uuid.UUID) error {
	err := service.categoryRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCategories implements CategoryService.
func (service *categoryServiceImpl) GetAllCategories() ([]response.CategoryResponse, error) {
	var categoriesResponses []response.CategoryResponse

	categoriesEntities, err := service.categoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, category := range categoriesEntities {
		categoriesResponses = append(categoriesResponses, mapper.CategoryEntityToCategoryResponse(category))
	}

	return categoriesResponses, err
}

func (service *categoryServiceImpl) GetAllCategoriesAllowedByRole(roleID string) ([]response.CategoryResponse, error) {
	var categoriesIDs []response.CategoryResponse

	categories, err := service.categoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if len(category.RolesAllowed) == 0 {
			categoriesIDs = append(categoriesIDs, mapper.CategoryEntityToCategoryResponse(category))
		}

		for _, categoryRole := range category.RolesAllowed {
			if categoryRole == roleID {
				categoriesIDs = append(categoriesIDs, mapper.CategoryEntityToCategoryResponse(category))
			}
		}
	}

	return categoriesIDs, nil
}

// GetCategoryByID implements CategoryService.
func (service *categoryServiceImpl) GetCategoryByID(id uuid.UUID) (*dto.CategoryDto, error) {
	categoryEntity, err := service.categoryRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return mapper.CategoryEntityToCategoryDto(categoryEntity), nil
}

// GetCategoryByName implements CategoryService.
func (service *categoryServiceImpl) GetCategoryByName(name string) (*dto.CategoryDto, error) {
	categoryEntity, err := service.categoryRepository.FindByName(name)
	if err != nil {
		return nil, err
	}

	return mapper.CategoryEntityToCategoryDto(categoryEntity), nil
}

// RestoreCategory implements CategoryService.
func (service *categoryServiceImpl) RestoreCategory(id uuid.UUID) error {
	err := service.categoryRepository.Restore(id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategoryDescription implements CategoryService.
func (service *categoryServiceImpl) UpdateCategory(id uuid.UUID, req request.NewCategory) error {
	existingCategory, err := service.categoryRepository.FindByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		existingCategory.Name = *req.Name
	}

	if req.Description != nil {
		existingCategory.Description = *req.Description
	}

	err = service.categoryRepository.Update(*existingCategory)
	if err != nil {
		return err
	}

	return nil
}

func NewCategoryService(categoryRepository repository.CategoryRepository, roleRepository repository.RoleRepository) CategoryService {
	return &categoryServiceImpl{categoryRepository: categoryRepository, roleRepository: roleRepository}
}
