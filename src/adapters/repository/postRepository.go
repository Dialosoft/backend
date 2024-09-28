package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	FindAll() ([]*models.Post, error)
	FindByID(ID uuid.UUID) (*models.Post, error)
	FindByUserID(userID uuid.UUID) ([]*models.Post, error)
	GetLikeCount(postID uuid.UUID) (int64, error)
	Create(post models.Post) (*models.Post, error)
	Update(postID uuid.UUID, updatedPost models.Post) error
	Delete(postID uuid.UUID) error
	Restore(postID uuid.UUID) error
}

type postRepositoryImpl struct {
	db *gorm.DB
}

// FindAll implements PostRepository.
func (repo *postRepositoryImpl) FindAll() ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.db.Preload("users").Find(&posts).Error; err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return posts, nil
}

// FindByID implements PostRepository.
func (repo *postRepositoryImpl) FindByID(ID uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := repo.db.Preload("users").Where("id = ?", ID.String()).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// FindByUserID implements PostRepository.
func (repo *postRepositoryImpl) FindByUserID(userID uuid.UUID) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.db.Preload("users").Where("user_id = ?", userID.String()).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// GetLikeCount implements PostRepository.
func (repo *postRepositoryImpl) GetLikeCount(postID uuid.UUID) (int64, error) {
	var likesCount int64
	if err := repo.db.Model(models.PostLikes{}).Where("post_id = ?", postID).Count(&likesCount).Error; err != nil {
		return 0, err
	}

	return likesCount, nil
}

// Create implements PostRepository.
func (repo *postRepositoryImpl) Create(post models.Post) (*models.Post, error) {
	result := repo.db.Save(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}

// Update implements PostRepository.
func (repo *postRepositoryImpl) Update(postID uuid.UUID, updatedPost models.Post) error {
	result := repo.db.Model(&models.Post{}).
		Where("id = ?", postID).
		Updates(updatedPost)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete implements PostRepository.
func (repo *postRepositoryImpl) Delete(postID uuid.UUID) error {
	return repo.db.Delete(&models.Post{}, postID.String()).Error
}

// Restore implements PostRepository.
func (repo *postRepositoryImpl) Restore(postID uuid.UUID) error {
	result := repo.db.Unscoped().
		Model(&models.Post{}).
		Where("id = ?", postID).
		Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepositoryImpl{db: db}
}
