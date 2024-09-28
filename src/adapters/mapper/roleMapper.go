package mapper

import (
	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/models"
)

// RoleDtoToRoleEntity converts a RoleDto to a RoleEntity.
// It returns an error if any of the required fields (RoleType or Permission) are empty or invalid.
func RoleDtoToRoleEntity(roleDto *dto.RoleDto) *models.RoleEntity {
	roleEntity := models.RoleEntity{
		ID:         roleDto.ID,
		RoleType:   roleDto.RoleType,
		Permission: roleDto.Permission,
		AdminRole:  roleDto.AdminRole,
		ModRole:    roleDto.ModRole,
		CreatedAt:  roleDto.CreatedAt,
		UpdatedAt:  roleDto.UpdatedAt,
		DeletedAt:  roleDto.DeletedAt,
	}

	return &roleEntity
}

// RoleEntityToRoleDto converts a RoleEntity to a RoleDto.
// Returns an error if the RoleEntity has missing required fields (RoleType or Permission).
func RoleEntityToRoleDto(roleEntity *models.RoleEntity) *dto.RoleDto {
	roleDto := dto.RoleDto{
		ID:         roleEntity.ID,
		RoleType:   roleEntity.RoleType,
		Permission: roleEntity.Permission,
		AdminRole:  roleEntity.AdminRole,
		ModRole:    roleEntity.ModRole,
		CreatedAt:  roleEntity.CreatedAt,
		UpdatedAt:  roleEntity.UpdatedAt,
		DeletedAt:  roleEntity.DeletedAt,
	}

	return &roleDto
}

func RoleEntityToRoleResponse(roleEntity *models.RoleEntity) response.RoleResponse {
	return response.RoleResponse{
		ID:         roleEntity.ID,
		RoleType:   roleEntity.RoleType,
		Permission: roleEntity.Permission,
		AdminRole:  roleEntity.AdminRole,
		ModRole:    roleEntity.ModRole,
		CreatedAt:  roleEntity.CreatedAt,
		UpdatedAt:  roleEntity.UpdatedAt,
		DeletedAt:  roleEntity.DeletedAt,
	}
}

func RoleResponseToRoleEntity(roleResponse *response.RoleResponse) *models.RoleEntity {
	return &models.RoleEntity{
		ID:         roleResponse.ID,
		RoleType:   roleResponse.RoleType,
		Permission: roleResponse.Permission,
		AdminRole:  roleResponse.AdminRole,
		ModRole:    roleResponse.ModRole,
		CreatedAt:  roleResponse.CreatedAt,
		UpdatedAt:  roleResponse.UpdatedAt,
		DeletedAt:  roleResponse.DeletedAt,
	}
}
