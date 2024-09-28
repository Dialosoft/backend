package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/domain/models"
)

func ForumDtoToForumEntity(forumDto *dto.ForumDto) *models.Forum {
	ForumEntity := models.Forum{
		ID:          forumDto.ID,
		Name:        forumDto.Name,
		Description: forumDto.Description,
		IsActive:    forumDto.IsActive,
		Type:        forumDto.Type,
		CategoryID:  forumDto.CategoryID,
		CreatedAt:   forumDto.CreatedAt,
		UpdatedAt:   forumDto.UpdatedAt,
	}

	return &ForumEntity
}

func ForumEntityToForumDto(forumEntity *models.Forum) *dto.ForumDto {
	ForumDto := dto.ForumDto{
		ID:          forumEntity.ID,
		Name:        forumEntity.Name,
		Description: forumEntity.Description,
		IsActive:    forumEntity.IsActive,
		Type:        forumEntity.Type,
		CategoryID:  forumEntity.CategoryID,
		CreatedAt:   forumEntity.CreatedAt,
		UpdatedAt:   forumEntity.UpdatedAt,
	}

	return &ForumDto
}
