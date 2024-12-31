package repository

import (
	"github.com/Dialosoft/src/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	AbstractRepository[*models.Category, uuid.UUID]
	// FindByName retrieves a category by its name.
	// Returns a pointer to the CategoryEntity if found, or an error otherwise.
	FindByName(name string) (*models.Category, error)
	// FindAllIncludingDeleted retrieves all categories from the database, including deleted ones.
	// Returns a slice of pointers to CategoryEntity and an error if something goes wrong.
	FindAllIncludingDeleted() ([]*models.Category, error)
}

type categoryRepositoryImpl struct {
	*abstractRepositoryImpl[*models.Category, uuid.UUID]
}

func NewCategoryRepository(gormDB *gorm.DB) CategoryRepository {
	repo := &categoryRepositoryImpl{}
	repo.abstractRepositoryImpl = CreateRepository(gormDB, repo)
	return repo
}

// Implementación de los métodos específicos
func (repo *categoryRepositoryImpl) FindAllIncludingDeleted() ([]*models.Category, error) {
	var categories []*models.Category

	result := repo.gorm.Unscoped().Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (repo *categoryRepositoryImpl) FindByName(name string) (*models.Category, error) {
	return repo.FirstByKey("name", name)
}
