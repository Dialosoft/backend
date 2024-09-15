package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAllRoles() ([]models.RoleEntity, error)
	FindByID(roleID uuid.UUID) (*models.RoleEntity, error)
	FindByType(roleType string) (*models.RoleEntity, error)
	Create(newRole models.RoleEntity) (uuid.UUID, error)
	Update(roleID uuid.UUID, updatedRole models.RoleEntity) error
	Delete(roleID uuid.UUID) error
	Restore(roleID uuid.UUID) error
}

type roleRepositoryImpl struct {
	db *gorm.DB
}

// FindAllRoles implements RoleRepository.
func (repo *roleRepositoryImpl) FindAllRoles() ([]models.RoleEntity, error) {
	var roles []models.RoleEntity
	if err := repo.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return roles, nil
}

// FindByID implements RoleRepository.
func (repo *roleRepositoryImpl) FindByID(roleID uuid.UUID) (*models.RoleEntity, error) {
	var role models.RoleEntity
	if err := repo.db.First(&role, "id = ?", roleID.String()).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByType implements RoleRepository.
func (repo *roleRepositoryImpl) FindByType(roleType string) (*models.RoleEntity, error) {
	var role models.RoleEntity
	if err := repo.db.First(&role, "role_type = ?", roleType).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create implements RoleRepository.
func (repo *roleRepositoryImpl) Create(newRole models.RoleEntity) (uuid.UUID, error) {
	result := repo.db.Create(&newRole)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}
	return newRole.ID, nil
}

// Update implements RoleRepository.
func (repo *roleRepositoryImpl) Update(roleID uuid.UUID, updatedRole models.RoleEntity) error {
	result := repo.db.Model(&models.RoleEntity{}).
		Where("id = ?", roleID).
		Updates(updatedRole)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete implements RoleRepository.
func (repo *roleRepositoryImpl) Delete(roleID uuid.UUID) error {
	return repo.db.Delete(&models.RoleEntity{}, roleID).Error
}

// Restore implements RoleRepository.
func (repo *roleRepositoryImpl) Restore(roleID uuid.UUID) error {
	result := repo.db.Unscoped().
		Model(&models.RoleEntity{}).
		Where("id = ?", roleID).
		Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}
