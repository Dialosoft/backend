package services

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostService provides an interface for managing posts in the system.
type PostService interface {
	// GetAllPosts retrieves a list of posts with pagination options.
	// limit specifies the maximum number of posts, and offset skips the given number of posts.
	GetAllPosts(limit, offset int) ([]response.PostResponse, error)

	// GetPostByID fetches a post based on its unique postID.
	GetPostByID(postID uuid.UUID) (*response.PostResponse, error)

	// GetPostsByUserID fetches all posts created by a specific user.
	GetPostsByUserID(userID uuid.UUID) ([]response.PostResponse, error)

	// GetAllPostsByForum retrieves posts from a specific forum with pagination options.
	GetAllPostsByForum(forumID uuid.UUID, limit, offset int) ([]response.PostResponse, error)

	// GetAllPostsAndReturnSimpleResponse retrieves posts with simplified response data and pagination options.
	GetAllPostsAndReturnSimpleResponse(limit, offset int) ([]response.SimplePostResponse, error)

	// GetLikeCount returns the number of likes for a specific post.
	GetLikeCount(postID uuid.UUID) (int64, error)

	// CreateNewPost creates a new post by a user.
	CreateNewPost(UserID uuid.UUID, post request.NewPost) (response.PostResponse, error)

	// UpdatePostTitle updates the title of a post identified by its postID.
	UpdatePostTitle(postID uuid.UUID, title string) error

	// UpdatePostContent updates the content of a post identified by its postID.
	UpdatePostContent(postID uuid.UUID, content string) error

	// DeletePost deletes a post identified by its postID.
	DeletePost(postID uuid.UUID) error

	// RestorePost restores a previously deleted post identified by its postID.
	RestorePost(postID uuid.UUID) error

	// LikePost allows a user to like a post identified by postID.
	LikePost(postID uuid.UUID, userID uuid.UUID) error

	// UnlikePost allows a user to remove their like from a post identified by postID.
	UnlikePost(postID uuid.UUID, userID uuid.UUID) error

	// GetPostLikesByUserID retrieves a list of post IDs that a user has liked.
	GetPostLikesByUserID(userID uuid.UUID) ([]uuid.UUID, error)
}

type postServiceImpl struct {
	postRepository repository.PostRepository
	postLikesRepo  repository.PostLikesRepository
	userRepository repository.UserRepository
}

// CreateNewPost implements PostService.
func (service *postServiceImpl) CreateNewPost(UserID uuid.UUID, post request.NewPost) (response.PostResponse, error) {
	userEntity, err := service.userRepository.FindByID(UserID)
	if err != nil {
		return response.PostResponse{}, err
	}

	forumUUID, err := uuid.Parse(post.ForumID)
	if err != nil {
		return response.PostResponse{}, err
	}

	postEntity := models.Post{
		UserID:  userEntity.ID,
		ForumID: forumUUID,
		User:    *userEntity,
		Title:   post.Title,
		Content: post.Content,
	}

	newPostEntity, err := service.postRepository.Create(postEntity)
	if err != nil {
		return response.PostResponse{}, err
	}
	return mapper.PostEntityToPostResponse(newPostEntity), nil
}

func (service *postServiceImpl) GetAllPostsByForum(forumID uuid.UUID, limit, offset int) ([]response.PostResponse, error) {
	var postResponses []response.PostResponse
	postsModels, err := service.postRepository.FindAllByForumID(forumID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, postModel := range postsModels {
		postResponses = append(postResponses, mapper.PostEntityToPostResponse(postModel))
	}

	return postResponses, nil
}

// GetAllPosts implements PostService.
func (service *postServiceImpl) GetAllPosts(limit, offset int) ([]response.PostResponse, error) {
	var postResponses []response.PostResponse
	postsModels, err := service.postRepository.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, postModel := range postsModels {
		postResponses = append(postResponses, mapper.PostEntityToPostResponse(postModel))
	}

	return postResponses, nil
}

// GetPostByID implements PostService.
func (service *postServiceImpl) GetPostByID(postID uuid.UUID) (*response.PostResponse, error) {
	postModel, err := service.postRepository.FindByID(postID)
	if err != nil {
		return nil, err
	}

	postResponse := mapper.PostEntityToPostResponse(postModel)

	return &postResponse, nil
}

// GetPostsByUserID implements PostService.
func (service *postServiceImpl) GetPostsByUserID(userID uuid.UUID) ([]response.PostResponse, error) {
	var postResponses []response.PostResponse
	postsModels, err := service.postRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, postModel := range postsModels {
		postResponses = append(postResponses, mapper.PostEntityToPostResponse(postModel))
	}

	return postResponses, nil
}

func (service *postServiceImpl) GetAllPostsAndReturnSimpleResponse(limit, offset int) ([]response.SimplePostResponse, error) {
	var postResponses []response.SimplePostResponse
	postsModels, err := service.postRepository.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, postModel := range postsModels {
		postResponses = append(postResponses, response.SimplePostResponse{
			ID:        postModel.ID.String(),
			UserID:    postModel.UserID.String(),
			Title:     postModel.Title,
			CreatedAt: postModel.CreatedAt.String(),
			UpdatedAt: postModel.UpdatedAt.String(),
			DeletedAt: postModel.DeletedAt,
		})
	}

	return postResponses, nil
}

func (service *postServiceImpl) GetLikeCount(postID uuid.UUID) (int64, error) {
	return service.postRepository.GetLikeCount(postID)
}

// UpdatePost implements PostService.
func (service *postServiceImpl) UpdatePostTitle(postID uuid.UUID, title string) error {

	modelPost, err := service.postRepository.FindByID(postID)
	if err != nil {
		return err
	}

	if modelPost.Title != title {
		modelPost.Title = title
	}

	return service.postRepository.Update(postID, *modelPost)
}

// UpdatePost implements PostService.
func (service *postServiceImpl) UpdatePostContent(postID uuid.UUID, content string) error {
	modelPost, err := service.postRepository.FindByID(postID)
	if err != nil {
		return err
	}

	if modelPost.Content != content {
		modelPost.Content = content
	}

	return service.postRepository.Update(postID, *modelPost)
}

// LikePost implements PostService.
func (service *postServiceImpl) LikePost(postID uuid.UUID, userID uuid.UUID) error {
	return service.postLikesRepo.Save(postID, userID)
}

// UnlikePost implements PostService.
func (service *postServiceImpl) UnlikePost(postID uuid.UUID, userID uuid.UUID) error {
	return service.postLikesRepo.Remove(postID, userID)
}

// DeletePost implements PostService.
func (service *postServiceImpl) DeletePost(postID uuid.UUID) error {
	return service.postRepository.Delete(postID)
}

// RestorePost implements PostService.
func (service *postServiceImpl) RestorePost(postID uuid.UUID) error {
	return service.postRepository.Restore(postID)
}

func (service *postServiceImpl) GetPostLikesByUserID(userID uuid.UUID) ([]uuid.UUID, error) {
	var postsIDs []uuid.UUID

	posts, err := service.postLikesRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		postsIDs = append(postsIDs, post.PostID)
	}

	if postsIDs == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return postsIDs, nil
}

func NewPostService(postRepository repository.PostRepository, postLikesRepo repository.PostLikesRepository, userRepository repository.UserRepository) PostService {
	return &postServiceImpl{postRepository: postRepository, postLikesRepo: postLikesRepo, userRepository: userRepository}
}
