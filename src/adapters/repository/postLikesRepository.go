package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostLikesRepository interface {
	FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error)
	FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error)
	FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error)
	Save(postID uuid.UUID, userID uuid.UUID) error
	Remove(postID uuid.UUID, userID uuid.UUID) error
}

type postLikesRepositoryImpl struct {
	db *gorm.DB
}

// FindAll implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error) {
	var postLikes []*models.PostLikes
	result := repo.db.Find(&postLikes, "post_id = ?", postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return postLikes, nil
}

// FindAllByUserID implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error) {
	var postLikes []*models.PostLikes
	if err := repo.db.Find(&postLikes, "post_id = ? AND user_id = ?", postID, userID).Error; err != nil {
		return nil, err
	}
	return postLikes, nil
}

func (repo *postLikesRepositoryImpl) FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error) {
	var postLikes []*models.PostLikes
	if err := repo.db.Find(&postLikes, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return postLikes, nil
}

// Save implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) Save(postID uuid.UUID, userID uuid.UUID) error {
	if err := repo.db.Create(&models.PostLikes{
		PostID: postID,
		UserID: userID,
	}).Error; err != nil {
		return err
	}
	return nil
}

// Remove implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) Remove(postID uuid.UUID, userID uuid.UUID) error {
	return repo.db.Delete(models.PostLikes{}, "post_id = ? AND user_id = ?", postID, userID).Error
}

func NewPostLikesRepository(db *gorm.DB) PostLikesRepository {
	return &postLikesRepositoryImpl{db: db}
}
