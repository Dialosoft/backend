package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	AbstractRepository[*models.UserEntity, uuid.UUID]
	// FindByUsername retrieves a user by their username from the database, including the associated role.
	// Returns a UserEntity pointer and an error if the user is not found or the operation fails.
	FindByUsername(username string) (*models.UserEntity, error)
}

type userRepositoryImpl struct {
	*abstractRepositoryImpl[*models.UserEntity, uuid.UUID]
}

func NewUserRepository(gormDB *gorm.DB) UserRepository {
	repo := &userRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(gormDB, repo)
	return repo
}

func (repo *userRepositoryImpl) FindByUsername(username string) (*models.UserEntity, error) {
	return repo.FirstByKey("username", username)
}

func (repo *userRepositoryImpl) GetPreloads() []string {
	return []string{"Role"}
}
