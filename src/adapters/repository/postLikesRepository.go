package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/domain/models"
)

// PostLikesRepository defines a set of methods for managing post likes in the system.
// Each method provides operations related to the PostLikesEntity model.
type PostLikesRepository interface {
	AbstractRepository[*models.PostLikes, uuid.UUID]

	// FindAllByPostID retrieves all likes for a specific post.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error)

	// FindAllByUserIDAndPostID retrieves all likes for a specific post by a user.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error)

	// FindAllByUserID retrieves all likes for a specific user.
	// Returns a slice of pointers to PostLikesEntity and an error if something goes wrong.
	FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error)

	// SaveLike inserts a new like for a specific post by a user.
	// Returns an error if the operation fails.
	SaveLike(postID uuid.UUID, userID uuid.UUID) error

	// RemoveLike removes a like for a specific post by a user.
	// Returns an error if the operation fails.
	RemoveLike(postID uuid.UUID, userID uuid.UUID) error
}

type postLikesRepositoryImpl struct {
	*abstractRepositoryImpl[*models.PostLikes, uuid.UUID]
}

// NewPostLikesRepository creates a new instance of PostLikesRepository
func NewPostLikesRepository(db *gorm.DB) PostLikesRepository {
	repo := &postLikesRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(db, repo)
	return repo
}

// FindAllByPostID implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error) {
	return repo.FindAllByKey("post_id", postID.String())
}

// FindAllByUserID implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error) {
	return repo.FindAllByKey("user_id", userID.String())
}

// FindAllByUserIDAndPostID implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error) {
	var postLikes []*models.PostLikes
	if err := repo.gorm.Find(&postLikes, "post_id = ? AND user_id = ?", postID, userID).Error; err != nil {
		return nil, err
	}
	return postLikes, nil
}

// SaveLike implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) SaveLike(postID uuid.UUID, userID uuid.UUID) error {
	postLike := &models.PostLikes{
		PostID: postID,
		UserID: userID,
	}
	_, err := repo.Create(nil, postLike)
	return err
}

// RemoveLike implements PostLikesRepository.
func (repo *postLikesRepositoryImpl) RemoveLike(postID uuid.UUID, userID uuid.UUID) error {
	return repo.gorm.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.PostLikes{}).Error
}
