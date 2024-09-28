package services

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

type PostService interface {
	GetAllPosts() ([]response.PostResponse, error)
	GetPostByID(postID uuid.UUID) (*response.PostResponse, error)
	GetPostsByUserID(userID uuid.UUID) ([]response.PostResponse, error)
	CreateNewPost(UserID uuid.UUID, post request.NewPost) (response.PostResponse, error)
	UpdatePostTitle(postID uuid.UUID, userID uuid.UUID, title string) error
	UpdatePostContent(postID uuid.UUID, userID uuid.UUID, content string) error
	DeletePost(postID uuid.UUID) error
	RestorePost(postID uuid.UUID) error
	LikePost(postID uuid.UUID, userID uuid.UUID) error
	UnlikePost(postID uuid.UUID, userID uuid.UUID) error
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

	postEntity := models.Post{
		UserID:   userEntity.ID,
		User:     *userEntity,
		Title:    post.Title,
		Content:  post.Content,
		Views:    post.Views,
		Comments: post.Comments,
	}

	newPostEntity, err := service.postRepository.Create(postEntity)
	if err != nil {
		return response.PostResponse{}, err
	}
	return mapper.PostEntityToPostResponse(newPostEntity), nil
}

// GetAllPosts implements PostService.
func (service *postServiceImpl) GetAllPosts() ([]response.PostResponse, error) {
	var postResponses []response.PostResponse
	postsModels, err := service.postRepository.FindAll()
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

// UpdatePost implements PostService.
func (service *postServiceImpl) UpdatePostTitle(postID uuid.UUID, userID uuid.UUID, title string) error {

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
func (service *postServiceImpl) UpdatePostContent(postID uuid.UUID, userID uuid.UUID, content string) error {
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

func NewPostService(postRepository repository.PostRepository, postLikesRepo repository.PostLikesRepository, userRepository repository.UserRepository) PostService {
	return &postServiceImpl{postRepository: postRepository, postLikesRepo: postLikesRepo, userRepository: userRepository}
}
