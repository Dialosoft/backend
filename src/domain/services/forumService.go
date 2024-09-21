package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
)

type ForumService interface {
	GetAllForums() ([]*models.Forum, error)
	GetForumByID(id uuid.UUID) (*models.Forum, error)
	GetForumByName(name string) (*models.Forum, error)
	CreateForum(forumDto dto.ForumDto) (uuid.UUID, error)
	UpdateForum(id uuid.UUID, forumDto dto.ForumDto) error
	DeleteForum(id uuid.UUID) error
	RestoreForum(id uuid.UUID) error
}

type forumServiceImpl struct {
	forumRepository repository.ForumRepository
}

// CreateForum implements ForumService.
func (service *forumServiceImpl) CreateForum(forumDto dto.ForumDto) (uuid.UUID, error) {
	forumEntity := mapper.ForumDtoToForumEntity(&forumDto)

	forumUUID, err := service.forumRepository.Create(*forumEntity)
	if err != nil {
		return uuid.UUID{}, err
	}

	return forumUUID, nil
}

// DeleteForum implements ForumService.
func (service *forumServiceImpl) DeleteForum(id uuid.UUID) error {
	err := service.forumRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllForums implements ForumService.
func (service *forumServiceImpl) GetAllForums() ([]*models.Forum, error) {
	forums, err := service.forumRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return forums, nil
}

// GetForumByID implements ForumService.
func (service *forumServiceImpl) GetForumByID(id uuid.UUID) (*models.Forum, error) {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

// GetForumByName implements ForumService.
func (service *forumServiceImpl) GetForumByName(name string) (*models.Forum, error) {
	forum, err := service.forumRepository.FindByName(name)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

// RestoreForum implements ForumService.
func (service *forumServiceImpl) RestoreForum(id uuid.UUID) error {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err = service.forumRepository.Restore(forum.ID); err != nil {
		return err
	}

	return nil
}

// UpdateForum implements ForumService.
func (service *forumServiceImpl) UpdateForum(id uuid.UUID, forumDto dto.ForumDto) error {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return err
	}

	{
		if forumDto.Name != "" {
			forum.Name = forumDto.Name
		}

		if forumDto.Description != "" {
			forum.Description = forumDto.Description
		}

		forum.IsActive = forumDto.IsActive

		if forumDto.Type != "" {
			forum.Type = forumDto.Type
		}

		if forumDto.PostCount != 0 {
			forum.PostCount = forumDto.PostCount
		}

		if forumDto.CategoryID != "" {
			forum.CategoryID = forumDto.CategoryID
		}
	}

	err = service.forumRepository.Update(*forum)
	if err != nil {
		return err
	}

	return nil
}

func NewForumService(forumRepository repository.ForumRepository) ForumService {
	return &forumServiceImpl{forumRepository: forumRepository}
}
