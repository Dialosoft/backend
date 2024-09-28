package mapper

import (
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
)

func PostEntityToPostResponse(postEntity *models.Post) response.PostResponse {
	return response.PostResponse{
		ID:        postEntity.ID,
		User:      UserEntityToUserResponse(&postEntity.User),
		Title:     postEntity.Title,
		Content:   postEntity.Content,
		Views:     postEntity.Views,
		Comments:  postEntity.Comments,
		CreatedAt: postEntity.CreatedAt,
		UpdatedAt: postEntity.UpdatedAt,
		DeletedAt: postEntity.DeletedAt,
	}
}

func PostResponseToPostEntity(postResponse *response.PostResponse) *models.Post {
	return &models.Post{
		ID:        postResponse.ID,
		UserID:    postResponse.User.ID,
		User:      *UserResponseToUserEntity(&postResponse.User),
		Title:     postResponse.Title,
		Content:   postResponse.Content,
		Views:     postResponse.Views,
		Comments:  postResponse.Comments,
		CreatedAt: postResponse.CreatedAt,
		UpdatedAt: postResponse.UpdatedAt,
		DeletedAt: postResponse.DeletedAt,
	}
}
