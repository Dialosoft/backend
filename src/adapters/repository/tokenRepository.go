package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TokenRepository defines a set of methods for managing tokens in the system.
// Each method provides operations related to the TokenEntity model.
type TokenRepository interface {
	AbstractRepository[*models.TokenEntity, uuid.UUID]

	// FindTokenByUserID retrieves a TokenEntity by the associated user's unique identifier (UUID).
	// Returns a pointer to the TokenEntity if found, or an error otherwise.
	FindTokenByUserID(userID uuid.UUID) (*models.TokenEntity, error)
}

type tokenRepositoryImpl struct {
	*abstractRepositoryImpl[*models.TokenEntity, uuid.UUID]
}

// FindTokenByUserID implements TokenRepository.
func (repo *tokenRepositoryImpl) FindTokenByUserID(userID uuid.UUID) (*models.TokenEntity, error) {
	return repo.FirstByKey("user_id", userID.String())
}

// NewTokenRepository creates a new instance of TokenRepository
func NewTokenRepository(db *gorm.DB) TokenRepository {
	repo := &tokenRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(db, repo)
	return repo
}
