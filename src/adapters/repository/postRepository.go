package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {

	// FindAll retrieves all posts from the database.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindAll(limit, offset int) ([]*models.Post, error)

	// FindByID retrieves a post by its unique identifier (UUID).
	// Returns a pointer to the PostEntity if found, or an error otherwise.
	FindByID(ID uuid.UUID) (*models.Post, error)

	// FindByUserID retrieves all posts created by a specific user.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindByUserID(userID uuid.UUID) ([]*models.Post, error)

	// FindAllByForumID retrieves all posts from a specific forum.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindAllByForumID(forumID uuid.UUID, limit, offset int) ([]*models.Post, error)

	// GetLikeCount returns the number of likes for a specific post.
	// Returns the number of likes and an error if the operation fails.
	GetLikeCount(postID uuid.UUID) (int64, error)

	// Create inserts a new post into the database.
	// Returns the pointer to the newly created post and an error if the operation fails.
	Create(post models.Post) (*models.Post, error)

	// Update modifies an existing post in the database identified by its ID.
	// Returns an error if the update fails or the post is not found.
	Update(postID uuid.UUID, updatedPost models.Post) error

	// Delete removes a post from the database identified by its ID.
	// Returns an error if the deletion fails or the post is not found.
	Delete(postID uuid.UUID) error

	// Restore restores a previously deleted post in the database identified by its ID.
	// Returns an error if the restoration fails or the post is not found.
	Restore(postID uuid.UUID) error
}

type postRepositoryImpl struct {
	db *gorm.DB
}

func (repo *postRepositoryImpl) FindAllByForumID(forumID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.db.Preload("User").
		Preload("User.Role").
		Where("forum_id = ?", forumID.String()).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// FindAll implements PostRepository.
func (repo *postRepositoryImpl) FindAll(limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.db.Preload("User").Preload("User.Role").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errorsUtils.ErrNoPostsObtained
	}

	return posts, nil
}

// FindByID implements PostRepository.
func (repo *postRepositoryImpl) FindByID(ID uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := repo.db.Preload("User").Preload("User.Role").Where("id = ?", ID.String()).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// FindByUserID implements PostRepository.
func (repo *postRepositoryImpl) FindByUserID(userID uuid.UUID) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.db.Preload("User").Preload("User.Role").Where("user_id = ?", userID.String()).Find(&posts).Error; err != nil {
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
