package repository

import (
	"errors"

	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ForumRepository defines a set of methods for managing forums in the system.
// Each method provides operations related to the ForumEntity model.
type ForumRepository interface {

	// FindAll retrieves all forums from the database.
	// Returns a slice of pointers to ForumEntity and an error if something goes wrong.
	FindAll() ([]*models.Forum, error)

	// FindAllWithDeleted retrieves all forums from the database, including deleted ones.
	// Returns a slice of pointers to ForumEntity and an error if something goes wrong.
	FindAllWithDeleted() ([]*models.Forum, error)

	// FindByID retrieves a forum by its unique identifier (UUID).
	// Returns a pointer to the ForumEntity if found, or an error otherwise.
	FindByID(uuid uuid.UUID) (*models.Forum, error)

	// FindByIDWithDeleted retrieves a forum by its unique identifier (UUID), including the associated category.
	// Returns a pointer to the ForumEntity if found, or an error otherwise.
	FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error)

	// FindByName retrieves a forum by its name.
	// Returns a pointer to the ForumEntity if found, or an error otherwise.
	FindByName(name string) (*models.Forum, error)

	// FindAllByCategoryID retrieves all forums by their category ID.
	// Returns a slice of ForumEntity pointers and an error if something goes wrong.
	FindAllByCategoryID(categoryID uuid.UUID) ([]models.Forum, error)

	// Create inserts a new forum into the database.
	// Returns the UUID of the newly created forum and an error if the operation fails.
	Create(forum models.Forum) (uuid.UUID, error)

	// Update modifies an existing forum in the database identified by its ID.
	// Returns an error if the update fails or the forum is not found.
	Update(forum models.Forum) error

	// UpdateCategoryOwner updates the category owner of a forum identified by its ID.
	// Returns an error if the update fails or the forum is not found.
	UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error

	// Delete removes a forum from the database identified by its ID.
	// Returns an error if the deletion fails or the forum is not found.
	Delete(uuid uuid.UUID) error

	// Restore restores a previously deleted forum in the database identified by its ID.
	// Returns an error if the restoration fails or the forum is not found.
	Restore(uuid uuid.UUID) error
}

type forumRepositoryImpl struct {
	db *gorm.DB
}

// Create implements ForumRepository.
func (repo *forumRepositoryImpl) Create(forum models.Forum) (uuid.UUID, error) {
	var category models.Category
	result := repo.db.Find(&category, "id = ?", forum.CategoryID)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	forum.CategoryID = category.ID.String()
	forum.Category = category

	result = repo.db.Create(&forum)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return category.ID, nil
}

// Delete implements ForumRepository.
func (repo *forumRepositoryImpl) Delete(uuid uuid.UUID) error {
	result := repo.db.Delete(&models.Forum{}, uuid.String())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindAll implements ForumRepository.
func (repo *forumRepositoryImpl) FindAll() ([]*models.Forum, error) {
	var forums []*models.Forum
	result := repo.db.Preload("Category").Find(&forums)
	if result.Error != nil {
		return nil, result.Error
	}

	return forums, nil
}

// FindAllWithDeleted implements ForumRepository.
func (repo *forumRepositoryImpl) FindAllWithDeleted() ([]*models.Forum, error) {
	var forums []*models.Forum
	result := repo.db.Unscoped().Find(&forums)
	if result.Error != nil {
		return nil, result.Error
	}

	return forums, nil
}

// FindByID implements ForumRepository.
func (repo *forumRepositoryImpl) FindByID(uuid uuid.UUID) (*models.Forum, error) {
	var forum models.Forum
	result := repo.db.Preload("Category").First(&forum, "id = ?", uuid.String())
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &forum, nil
}

// FindByIDWithDeleted implements ForumRepository.
func (repo *forumRepositoryImpl) FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error) {
	var forum models.Forum
	result := repo.db.Unscoped().First(&forum, "id = ?", uuid.String())
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &forum, nil
}

// FindByName implements ForumRepository.
func (repo *forumRepositoryImpl) FindByName(name string) (*models.Forum, error) {
	var forum models.Forum
	result := repo.db.First(&forum, "name = ?", name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &forum, nil
}

func (repo *forumRepositoryImpl) FindAllByCategoryID(categoryID uuid.UUID) ([]models.Forum, error) {
	var forums []models.Forum
	if result := repo.db.Where("category_id = ?", categoryID).Find(&forums); result.Error != nil {
		return nil, result.Error
	}

	return forums, nil
}

// Restore implements ForumRepository.
func (repo *forumRepositoryImpl) Restore(uuid uuid.UUID) error {

	result := repo.db.Unscoped().Model(&models.Forum{}).Where("id = ?", uuid).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update implements ForumRepository.
func (repo *forumRepositoryImpl) Update(forum models.Forum) error {

	result := repo.db.Model(forum).Updates(forum)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateCategoryOwner implements ForumRepository.
func (repo *forumRepositoryImpl) UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error {
	var category models.Category
	var forum models.Forum

	resultCat := repo.db.Find(&category, "id = ?", categoryID)
	if resultCat.Error != nil {
		return resultCat.Error
	}

	resultFor := repo.db.Find(&forum, "id = ?", id)
	if resultFor.Error != nil {
		return resultFor.Error
	}

	forum.CategoryID = category.ID.String()
	forum.Category = category

	resultSave := repo.db.Save(&forum)
	if resultSave.Error != nil {
		return resultSave.Error
	}

	return nil
}

func NewForumRepository(db *gorm.DB) ForumRepository {
	return &forumRepositoryImpl{db: db}
}
