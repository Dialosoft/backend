package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ForumRepository defines a set of methods for managing forums in the system.
// Each method provides operations related to the ForumEntity model.
type ForumRepository interface {
	AbstractRepository[*models.Forum, uuid.UUID]
	// FindByName retrieves a forum by its name.
	// Returns a pointer to the ForumEntity if found, or an error otherwise.
	FindByName(name string) (*models.Forum, error)

	// FindAllWithDeleted retrieves all forums from the database, including deleted ones.
	// Returns a slice of pointers to ForumEntity and an error if something goes wrong.
	FindAllWithDeleted() ([]*models.Forum, error)

	// FindByIDWithDeleted retrieves a forum by its unique identifier (UUID), including the associated category.
	// Returns a pointer to the ForumEntity if found, or an error otherwise.
	FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error)

	// FindAllByCategoryID retrieves all forums by their category ID.
	// Returns a slice of ForumEntity pointers and an error if something goes wrong.
	FindAllByCategoryID(categoryID uuid.UUID) ([]*models.Forum, error)

	// UpdateCategoryOwner updates the category owner of a forum identified by its ID.
	// Returns an error if the update fails or the forum is not found.
	UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error
}

type forumRepositoryImpl struct {
	*abstractRepositoryImpl[*models.Forum, uuid.UUID]
}

func NewForumRepository(gormDB *gorm.DB) ForumRepository {
	repo := &forumRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(gormDB, repo)
	return repo
}

func (repo *forumRepositoryImpl) FindByName(name string) (*models.Forum, error) {
	return repo.FirstByKey("name", name)
}

func (repo *forumRepositoryImpl) FindAllWithDeleted() ([]*models.Forum, error) {
	var forums []*models.Forum
	result := repo.gorm.Unscoped().Find(&forums)
	if result.Error != nil {
		return nil, result.Error
	}
	return forums, nil
}

func (repo *forumRepositoryImpl) FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error) {
	var forum *models.Forum
	result := repo.gorm.Unscoped().Where("id = ?", uuid).First(&forum)
	if result.Error != nil {
		return nil, result.Error
	}
	return forum, nil
}

func (repo *forumRepositoryImpl) GetPreloads() []string {
	return []string{"Category"}
}

func (repo *forumRepositoryImpl) FindAllByCategoryID(categoryID uuid.UUID) ([]*models.Forum, error) {
	return repo.FindAllByKey("category_id", categoryID.String())
}

func (repo *forumRepositoryImpl) UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error {
	result := repo.gorm.Model(&models.Forum{}).Where("id = ?", id).Update("category_id", categoryID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
