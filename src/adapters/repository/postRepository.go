package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	AbstractRepository[*models.Post, uuid.UUID]

	// FindAll retrieves all posts from the database.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindAllWithPagination(limit, offset int) ([]*models.Post, error)

	// FindByUserID retrieves all posts created by a specific user.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindByUserID(userID uuid.UUID) ([]*models.Post, error)

	// FindAllByForumID retrieves all posts from a specific forum.
	// Returns a slice of pointers to PostEntity and an error if something goes wrong.
	FindAllByForumID(forumID uuid.UUID, limit, offset int) ([]*models.Post, error)

	// GetLikeCount returns the number of likes for a specific post.
	// Returns the number of likes and an error if the operation fails.
	GetLikeCount(postID uuid.UUID) (int64, error)
}

type postRepositoryImpl struct {
	*abstractRepositoryImpl[*models.Post, uuid.UUID]
}

func NewPostRepository(db *gorm.DB) PostRepository {
	repo := &postRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(db, repo)
	return repo
}

func (repo *postRepositoryImpl) FindAllByForumID(forumID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.TransCheck(nil).Preload("User").
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
func (repo *postRepositoryImpl) FindAllWithPagination(limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.TransCheck(nil).Preload("User").Preload("User.Role").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errorsUtils.ErrNoPostsObtained
	}

	return posts, nil
}

func (repo *postRepositoryImpl) FindByUserID(userID uuid.UUID) ([]*models.Post, error) {
	var posts []*models.Post
	if err := repo.TransCheck(nil).Preload("User").
		Preload("User.Role").
		Where("user_id = ?", userID.String()).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// GetLikeCount implements PostRepository.
func (repo *postRepositoryImpl) GetLikeCount(postID uuid.UUID) (int64, error) {
	var count int64
	if err := repo.TransCheck(nil).
		Model(&models.Post{}).
		Where("id = ?", postID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPreloads implements PostRepository.
func (repo *postRepositoryImpl) GetPreloads() []string {
	return []string{"User", "User.Role", "Forum"}
}
