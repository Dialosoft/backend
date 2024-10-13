package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostLikesRepository defines a set of methods for managing post likes in the system.
// Each method provides operations related to the PostLikesEntity model.
type PostLikesRepository interface {

	// FindAllByPostID retrieves all likes for a specific post.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error)

	// FindAllByUserIDAndPostID retrieves all likes for a specific post by a user.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error)

	// FindAllByUserID retrieves all likes for a specific user.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error)

	// Save inserts a new like for a specific post by a user.
	// Returns an error if the operation fails.
	Save(postID uuid.UUID, userID uuid.UUID) error

	// Remove removes a like for a specific post by a user.
	// Returns an error if the operation fails.
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
