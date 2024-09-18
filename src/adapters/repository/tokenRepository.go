package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository interface {
	FindAllTokens() ([]*models.TokenEntity, error)
	FindTokenByID(tokenID uuid.UUID) (*models.TokenEntity, error)
	FindTokenByUserID(userID uuid.UUID) (*models.TokenEntity, error)
	Save(tokenEntity models.TokenEntity) error
	Update(tokenID uuid.UUID, tokenEntity models.TokenEntity) error
	Delete(tokenID uuid.UUID) error
}

type tokenRepositoryImpl struct {
	db *gorm.DB
}

// FindAllTokens implements TokenRepository.
func (repo *tokenRepositoryImpl) FindAllTokens() ([]*models.TokenEntity, error) {
	var tokens []*models.TokenEntity
	if err := repo.db.Find(&tokens).Error; err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return tokens, nil
}

// FindTokenByID implements TokenRepository.
func (repo *tokenRepositoryImpl) FindTokenByID(tokenID uuid.UUID) (*models.TokenEntity, error) {
	var token models.TokenEntity
	if err := repo.db.First(&token, "id = ?", tokenID.String()).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

// FindTokenByUserID implements TokenRepository.
func (repo *tokenRepositoryImpl) FindTokenByUserID(userID uuid.UUID) (*models.TokenEntity, error) {
	var token models.TokenEntity
	if err := repo.db.First(&token, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

// Save implements TokenRepository.
func (repo *tokenRepositoryImpl) Save(tokenEntity models.TokenEntity) error {
	result := repo.db.Create(&tokenEntity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update implements TokenRepository.
func (repo *tokenRepositoryImpl) Update(tokenID uuid.UUID, updatedToken models.TokenEntity) error {
	result := repo.db.Model(&models.RoleEntity{}).
		Where("id = ?", tokenID).
		Updates(updatedToken)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete implements TokenRepository.
func (repo *tokenRepositoryImpl) Delete(tokenID uuid.UUID) error {
	return repo.db.Delete(&models.RoleEntity{}, tokenID).Error
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryImpl{db: db}
}
