package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {

	// FindAll retrieves all categories from the database.
	// Returns a slice of pointers to CategoryEntity and an error if something goes wrong.
	FindAll() ([]*models.Category, error)

	// FindByID retrieves a category by its unique identifier (UUID).
	// Returns a pointer to the CategoryEntity if found, or an error otherwise.
	FindByID(uuid uuid.UUID) (*models.Category, error)

	// FindByName retrieves a category by its name.
	// Returns a pointer to the CategoryEntity if found, or an error otherwise.
	FindByName(name string) (*models.Category, error)

	// FindAllIncludingDeleted retrieves all categories from the database, including deleted ones.
	// Returns a slice of pointers to CategoryEntity and an error if something goes wrong.
	FindAllIncludingDeleted() ([]*models.Category, error)

	// Create inserts a new category into the database.
	// Returns the UUID of the newly created category and an error if the operation fails.
	Create(category models.Category) (uuid.UUID, error)

	// Update modifies an existing category in the database identified by its ID.
	// Returns an error if the update fails or the category is not found.
	Update(category models.Category) error

	// Delete removes a category from the database identified by its ID.
	// Returns an error if the deletion fails or the category is not found.
	Delete(uuid uuid.UUID) error

	// Restore restores a previously deleted category in the database identified by its ID.
	// Returns an error if the restoration fails or the category is not found.
	Restore(uuid uuid.UUID) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

// Create implements CategoryRepository.
func (repo *categoryRepositoryImpl) Create(category models.Category) (uuid.UUID, error) {
	if err := repo.db.Create(&category).Error; err != nil {
		return uuid.UUID{}, err
	}

	return category.ID, nil
}

// Delete implements CategoryRepository.
func (repo *categoryRepositoryImpl) Delete(uuid uuid.UUID) error {
	panic("unimplemented")
}

// FindAll implements CategoryRepository.
func (repo *categoryRepositoryImpl) FindAll() ([]*models.Category, error) {
	var Categories []*models.Category
	result := repo.db.Find(&Categories)

	if result.Error != nil {
		return nil, result.Error
	}
	return Categories, nil
}

// FindAllIncludingDeleted implements CategoryRepository.
func (repo *categoryRepositoryImpl) FindAllIncludingDeleted() ([]*models.Category, error) {
	var categories []*models.Category

	result := repo.db.Unscoped().Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

// FindByID implements CategoryRepository.
func (repo *categoryRepositoryImpl) FindByID(uuid uuid.UUID) (*models.Category, error) {
	var category models.Category
	result := repo.db.First(&category, "id = ?", uuid.String())
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

// FindByName implements CategoryRepository.
func (repo *categoryRepositoryImpl) FindByName(name string) (*models.Category, error) {
	var category models.Category
	result := repo.db.First(&category, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

// Restore implements CategoryRepository.
func (repo *categoryRepositoryImpl) Restore(uuid uuid.UUID) error {
	result := repo.db.Unscoped().Model(&models.Category{}).Where("id = ?", uuid.String()).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update implements CategoryRepository.
func (repo *categoryRepositoryImpl) Update(category models.Category) error {
	result := repo.db.Model(&category).Where("id = ?", category.ID).Updates(category)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}
