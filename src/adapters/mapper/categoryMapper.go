package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
)

func CategoryDtoToCategoryEntity(categoryDto *dto.CategoryDto) *models.Category {
	categoryEntity := models.Category{
		ID:          categoryDto.ID,
		Name:        categoryDto.Name,
		Description: categoryDto.Description,
		CreatedAt:   categoryDto.CreatedAt,
		UpdatedAt:   categoryDto.UpdatedAt,
	}

	return &categoryEntity
}

func CategoryEntityToCategoryDto(categoryEntity *models.Category) *dto.CategoryDto {
	categoryDto := dto.CategoryDto{
		ID:          categoryEntity.ID,
		Name:        categoryEntity.Name,
		Description: categoryEntity.Description,
		CreatedAt:   categoryEntity.CreatedAt,
		UpdatedAt:   categoryEntity.UpdatedAt,
	}

	return &categoryDto
}

func CategoryCreateRequestToCategoryDto(categoryRequest *request.NewCategory) *dto.CategoryDto {
	categoryDto := dto.CategoryDto{
		Name:        *categoryRequest.Name,
		Description: *categoryRequest.Description,
	}

	return &categoryDto
}

func CategoryEntityToCategoryResponse(categoryEntity *models.Category) response.CategoryResponse {
	return response.CategoryResponse{
		ID:           categoryEntity.ID,
		Name:         categoryEntity.Name,
		Description:  categoryEntity.Description,
		RolesAllowed: categoryEntity.RolesAllowed,
		CreatedAt:    categoryEntity.CreatedAt,
		UpdatedAt:    categoryEntity.UpdatedAt,
	}
}
