package services

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/google/uuid"
)

type ForumService interface {
	GetAllForums() ([]*dto.ForumDto, error)
	GetForumByID(id uuid.UUID) (*dto.ForumDto, error)
	GetForumByName(name string) (*dto.ForumDto, error)
	CreateForum(forumDto dto.ForumDto) (uuid.UUID, error)
	UpdateForum(id uuid.UUID, req request.NewForum) error
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
func (service *forumServiceImpl) GetAllForums() ([]*dto.ForumDto, error) {
	var forumsDtos []*dto.ForumDto

	forums, err := service.forumRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, v := range forums {
		forumDto := mapper.ForumEntityToForumDto(v)
		forumsDtos = append(forumsDtos, forumDto)
	}

	return forumsDtos, nil
}

// GetForumByID implements ForumService.
func (service *forumServiceImpl) GetForumByID(id uuid.UUID) (*dto.ForumDto, error) {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return mapper.ForumEntityToForumDto(forum), nil
}

// GetForumByName implements ForumService.
func (service *forumServiceImpl) GetForumByName(name string) (*dto.ForumDto, error) {
	forum, err := service.forumRepository.FindByName(name)
	if err != nil {
		return nil, err
	}

	return mapper.ForumEntityToForumDto(forum), nil
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
func (service *forumServiceImpl) UpdateForum(id uuid.UUID, req request.NewForum) error {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return err
	}

	{
		if req.Name != nil {
			forum.Name = *req.Name
		}

		if req.Description != nil {
			forum.Description = *req.Description
		}

		if req.IsActive != nil {
			forum.IsActive = *req.IsActive
		}

		if req.Type != nil {
			forum.Type = *req.Type
		}

		if req.PostCount != nil {
			forum.PostCount = *req.PostCount
		}

		if req.CategoryID != nil {
			forum.CategoryID = *req.CategoryID
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
