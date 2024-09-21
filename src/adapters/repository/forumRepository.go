package repository

import (
	"errors"

	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumRepository interface {
	FindAll() ([]*models.Forum, error)
	FindAllWithDeleted() ([]*models.Forum, error)
	FindByID(uuid uuid.UUID) (*models.Forum, error)
	FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error)
	FindByName(name string) (*models.Forum, error)
	Create(forum models.Forum) error
	Update(forum models.Forum) error
	UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error
	Delete(uuid uuid.UUID) error
	Restore(uuid uuid.UUID) error
}

type forumRepositoryImpl struct {
	db *gorm.DB
}

// Create implements ForumRepository.
func (repo *forumRepositoryImpl) Create(forum models.Forum) error {
	var category models.Category
	result := repo.db.Find(&category, "id = ?", forum.CategoryID)
	if result.Error != nil {
		return result.Error
	}

	forum.CategoryID = category.ID.String()
	forum.Category = category

	result = repo.db.Create(forum)
	if result.Error != nil {
		return result.Error
	}

	return nil
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
