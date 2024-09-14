package repository

import (
	"errors"

	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUsers() ([]models.UserEntity, error)
	FindByID(id uuid.UUID) (*models.UserEntity, error)
	FindByUsername(username string) (*models.UserEntity, error)
	Create(newUser models.UserEntity) (uuid.UUID, error)
	Update(userId uuid.UUID, updatedUser models.UserEntity) error
	Delete(userId uuid.UUID) error
	Restore(userId uuid.UUID) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// GetAllUsers implements UserRepository.
func (repo *userRepositoryImpl) FindAllUsers() ([]models.UserEntity, error) {
	var users []models.UserEntity
	if err := repo.db.Preload("Role").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID implements UserRepository.
func (repo *userRepositoryImpl) FindByID(id uuid.UUID) (*models.UserEntity, error) {
	var user models.UserEntity
	if err := repo.db.Preload("Role").
		Where("id = ?", id.String()).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername implements UserRepository.
func (repo *userRepositoryImpl) FindByUsername(username string) (*models.UserEntity, error) {
	var user models.UserEntity
	if err := repo.db.Preload("Role").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create implements UserRepository.
func (repo *userRepositoryImpl) Create(newUser models.UserEntity) (uuid.UUID, error) {
	result := repo.db.Create(&newUser)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}
	return newUser.ID, nil
}

// Update implements UserRepository.
func (repo *userRepositoryImpl) Update(userID uuid.UUID, updatedUser models.UserEntity) error {
	result := repo.db.Model(&models.UserEntity{}).
		Where("id = ?", userID).
		Updates(updatedUser)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no user found with the given id")
	}

	return nil
}

// Delete implements UserRepository.
func (repo *userRepositoryImpl) Delete(userId uuid.UUID) error {
	return repo.db.Delete(&models.UserEntity{}, userId).Error
}

// Restore implements UserRepository.
func (repo *userRepositoryImpl) Restore(userId uuid.UUID) error {
	result := repo.db.Unscoped().Model(&models.UserEntity{}).Where("id = ?", userId).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
