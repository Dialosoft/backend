package services

import (
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/google/uuid"
)

// ForumService defines the methods for managing forums in the system.
type ForumService interface {
	// GetAllForums retrieves a list of all available forums.
	// Returns a slice of ForumDto or an error if something goes wrong.
	GetAllForums() ([]response.ForumResponse, error)

	// GetForumByID retrieves a specific forum by its unique ID.
	// Returns the ForumDto or an error if the forum is not found.
	GetForumByID(id uuid.UUID) (response.ForumResponse, error)

	// GetForumByName retrieves a specific forum by its name.
	// Returns the ForumDto or an error if the forum is not found.
	GetForumByName(name string) (response.ForumResponse, error)

	// GetForumsByCategoryIDAndAllowed retrieves a list of forums by their category ID and user role.
	// Returns a slice of ForumResponse or an error if something goes wrong.
	GetForumsByCategoryIDAndAllowed(categoryID uuid.UUID, userRole string) ([]response.ForumResponse, error)

	// CreateForum adds a new forum based on the provided ForumDto.
	// Returns the UUID of the newly created forum or an error if creation fails.
	CreateForum(newRequest request.NewForum) (uuid.UUID, error)

	// UpdateForum updates an existing forum's information by its ID.
	// The updated data is provided via the NewForum request structure.
	// Returns an error if the update fails or the forum is not found.
	UpdateForum(id uuid.UUID, req request.NewForum) error

	// DeleteForum removes a forum by its ID.
	// Returns an error if the deletion fails or the forum is not found.
	DeleteForum(id uuid.UUID) error

	// RestoreForum restores a previously deleted forum by its ID.
	// Returns an error if the restoration fails or the forum is not found.
	RestoreForum(id uuid.UUID) error
}

type forumServiceImpl struct {
	forumRepository    repository.ForumRepository
	categoryRepository repository.CategoryRepository
}

func (service *forumServiceImpl) GetForumsByCategoryIDAndAllowed(categoryID uuid.UUID, userRole string) ([]response.ForumResponse, error) {
	var forumsResponse []response.ForumResponse
	forums, err := service.forumRepository.FindAllByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	for _, forum := range forums {
		if forum.RolesAllowed == nil {
			forumsResponse = append(forumsResponse, mapper.ForumEntityToForumResponse(&forum))
		}

		for _, role := range forum.RolesAllowed {
			if role == userRole {
				forumsResponse = append(forumsResponse, mapper.ForumEntityToForumResponse(&forum))
			}
		}
	}

	return forumsResponse, nil
}

// CreateForum implements ForumService.
func (service *forumServiceImpl) CreateForum(newRequest request.NewForum) (uuid.UUID, error) {
	forumEntity := mapper.ForumNewRequestToForumEntity(newRequest)
	forumEntity.IsActive = true

	forumUUID, err := service.forumRepository.Create(forumEntity)
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
func (service *forumServiceImpl) GetAllForums() ([]response.ForumResponse, error) {
	var forumsResponses []response.ForumResponse

	forums, err := service.forumRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, forum := range forums {
		forumsResponses = append(forumsResponses, mapper.ForumEntityToForumResponse(forum))
	}

	return forumsResponses, nil
}

// GetForumByID implements ForumService.
func (service *forumServiceImpl) GetForumByID(id uuid.UUID) (response.ForumResponse, error) {
	forum, err := service.forumRepository.FindByID(id)
	if err != nil {
		return response.ForumResponse{}, err
	}

	return mapper.ForumEntityToForumResponse(forum), nil
}

// GetForumByName implements ForumService.
func (service *forumServiceImpl) GetForumByName(name string) (response.ForumResponse, error) {
	forum, err := service.forumRepository.FindByName(name)
	if err != nil {
		return response.ForumResponse{}, err
	}

	return mapper.ForumEntityToForumResponse(forum), nil
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

	// {
	// 	if req.Name != nil {
	// 		forum.Name = *req.Name
	// 	}

	// 	if req.Description != nil {
	// 		forum.Description = *req.Description
	// 	}

	// 	if req.IsActive != nil {
	// 		forum.IsActive = *req.IsActive
	// 	}

	// 	if req.Type != nil {
	// 		forum.Type = *req.Type
	// 	}

	// 	if req.CategoryID != nil {
	// 		forum.CategoryID = *req.CategoryID
	// 	}
	// }

	err = service.forumRepository.Update(*forum)
	if err != nil {
		return err
	}

	return nil
}

func NewForumService(forumRepository repository.ForumRepository, categoryRepository repository.CategoryRepository) ForumService {
	return &forumServiceImpl{forumRepository: forumRepository, categoryRepository: categoryRepository}
}
