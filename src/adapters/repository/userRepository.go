package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUsers() ([]*models.UserEntity, error)
	FindByID(id uuid.UUID) (*models.UserEntity, error)
	FindByUsername(username string) (*models.UserEntity, error)
	Create(newUser models.UserEntity) (uuid.UUID, error)
	Update(userID uuid.UUID, updatedUser models.UserEntity) error
	Delete(userID uuid.UUID) error
	Restore(userID uuid.UUID) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// FindAllUsers retrieves all users from the database, including their associated roles.
// Returns a slice of UserEntity pointers and an error if the operation fails.
// If no users are found, returns gorm.ErrRecordNotFound.
//
//	gorm.ErrRecordNotFound = "ErrRecordNotFound record not found error"
func (repo *userRepositoryImpl) FindAllUsers() ([]*models.UserEntity, error) {
	var users []*models.UserEntity
	if err := repo.db.Preload("Role").
		Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return users, nil
}

// FindByID retrieves a user by their UUID from the database, including the associated role.
// Returns a UserEntity pointer and an error if the user is not found or the operation fails.
func (repo *userRepositoryImpl) FindByID(id uuid.UUID) (*models.UserEntity, error) {
	var user models.UserEntity
	if err := repo.db.Preload("Role").
		Where("id = ?", id.String()).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by their username from the database, including the associated role.
// Returns a UserEntity pointer and an error if the user is not found or the operation fails.
func (repo *userRepositoryImpl) FindByUsername(username string) (*models.UserEntity, error) {
	var user models.UserEntity
	if err := repo.db.Preload("Role").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user into the database.
// Returns the UUID of the newly created user and an error if the operation fails.
func (repo *userRepositoryImpl) Create(newUser models.UserEntity) (uuid.UUID, error) {
	result := repo.db.Create(&newUser)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}
	return newUser.ID, nil
}

// Update modifies an existing user in the database identified by userID.
// Returns an error if the update fails or if the user does not exist.
//
//	gorm.ErrRecordNotFound = "ErrRecordNotFound record not found error"
func (repo *userRepositoryImpl) Update(userID uuid.UUID, updatedUser models.UserEntity) error {
	result := repo.db.Model(&models.UserEntity{}).
		Where("id = ?", userID).
		Updates(updatedUser)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes a user from the database identified by userID.
// Returns an error if the deletion fails.
func (repo *userRepositoryImpl) Delete(userID uuid.UUID) error {
	return repo.db.Delete(&models.UserEntity{}, userID).Error
}

// Restore restores a soft-deleted user in the database identified by userID.
// Returns an error if the restore operation fails.
func (repo *userRepositoryImpl) Restore(userID uuid.UUID) error {
	result := repo.db.Unscoped().
		Model(&models.UserEntity{}).
		Where("id = ?", userID).
		Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
