package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

// CategoryService defines the methods for managing categories in the system.
type CategoryService interface {
	// GetAllCategories retrieves a list of all available categories.
	// Returns a slice of CategoryDto or an error if something goes wrong.
	GetAllCategories() ([]*dto.CategoryDto, error)

	// GetCategoryByID retrieves a specific category by its unique ID.
	// Returns the CategoryDto or an error if the category is not found.
	GetCategoryByID(id uuid.UUID) (*dto.CategoryDto, error)

	// GetCategoryByName retrieves a specific category by its name.
	// Returns the CategoryDto or an error if the category is not found.
	GetCategoryByName(name string) (*dto.CategoryDto, error)

	// CreateCategory adds a new category based on the provided CategoryDto.
	// Returns the UUID of the newly created category or an error if creation fails.
	CreateCategory(categoryDto dto.CategoryDto) (uuid.UUID, error)

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
}

// CreateCategory implements CategoryService.
func (service *categoryServiceImpl) CreateCategory(categoryDto dto.CategoryDto) (uuid.UUID, error) {
	newCategory := models.Category{
		Name:        categoryDto.Name,
		Description: categoryDto.Description,
	}

	id, err := service.categoryRepository.Create(newCategory)
	if err != nil {
		return uuid.UUID{}, nil
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
func (service *categoryServiceImpl) GetAllCategories() ([]*dto.CategoryDto, error) {
	var categoriesDtos []*dto.CategoryDto

	categoriesEntities, err := service.categoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, v := range categoriesEntities {
		categoryDto := mapper.CategoryEntityToCategoryDto(v)
		categoriesDtos = append(categoriesDtos, categoryDto)
	}

	return categoriesDtos, err
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

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{categoryRepository: categoryRepository}
}
