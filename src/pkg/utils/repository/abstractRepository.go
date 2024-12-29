package repository

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ID is a generic type that represents the primary key of the entity, only accepting uint or uuid.UUID.
type ID interface {
	uint | uuid.UUID
}

// Identifiable is a generic interface that represents an entity that has an ID.
type Identifiable[K ID] interface {
	GetID() K
}

// T is a generic type that represents a database entity.
// K is a generic type that represents the primary key of the entity, only accepting uint or uuid.UUID.
type AbstractRepository[T Identifiable[K], K ID] interface {

	// FindAll retrieves all entities of type T from the database.
	FindAll() ([]T, error)

	// FindByID retrieves a single entity of type T by its ID.
	FindByID(id K) (T, error)

	// FirstByKey retrieves a single entity of type T by a specific field (key),thats mean
	// only the first Match!
	// The `key` parameter specifies the field to search, and `value` is the value to match.
	//
	// if you want to find all use:
	//	 FindAllByKey(key, value)
	FirstByKey(key, value string) (T, error)

	// FindAllByKey retrieves all entities of type T by a specific field (key)
	// The `key` parameter specifies the field to search, and `value` is the value to match.
	FindAllByKey(key, value string) ([]T, error)

	// Create inserts a new entity of type T into the database and returns its ID.
	// The operation can optionally be executed within a transaction.
	Create(tx *gorm.DB, newEntity T) (T, error)

	// Update modifies an existing entity of type T identified by its ID.
	// The operation can optionally be executed within a transaction.
	Update(tx *gorm.DB, id K, newEntity T) error

	// Delete marks an entity of type T as deleted (soft delete) by its ID.
	// The operation can optionally be executed within a transaction.
	Delete(tx *gorm.DB, id K) error

	// Restore unmarks an entity of type T as deleted (restore) by its ID.
	// The operation can optionally be executed within a transaction.
	Restore(tx *gorm.DB, id K) error

	// GetPreloads returns the default preloads for the repository.
	// 	This need to be overriden by the concrete implementation!!
	// by default is nil
	GetPreloads() []string

	// GetType returns the types defined of the repository.
	GetType() string
}

type abstractRepositoryImpl[T Identifiable[K], K ID] struct {
	gorm *gorm.DB
	self AbstractRepository[T, K]
}

// FindAll implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) FindAll() ([]T, error) {
	var entities []T

	if repo.self == nil {
		return entities, errors.New("self reference is nil in repository")
	}

	preloads := repo.self.GetPreloads()
	db := applyPreloads(repo.gorm, preloads)

	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return entities, nil
}

// FindByID implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) FindByID(id K) (T, error) {
	var entity T

	if repo.self == nil {
		return entity, errors.New("self reference is nil in repository")
	}

	preloads := repo.self.GetPreloads()
	db := applyPreloads(repo.gorm, preloads)

	if err := db.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// FirstByKey implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) FirstByKey(key, value string) (T, error) {
	var entity T

	if repo.self == nil {
		return entity, errors.New("self reference is nil in repository")
	}

	preloads := repo.self.GetPreloads()
	db := applyPreloads(repo.gorm, preloads)
	query := fmt.Sprintf("%s = ?", key)

	if err := db.Where(query, value).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// FindByKey implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) FindAllByKey(key, value string) ([]T, error) {
	var entities []T

	if repo.self == nil {
		return entities, errors.New("self reference is nil in repository")
	}

	preloads := repo.self.GetPreloads()
	db := applyPreloads(repo.gorm, preloads)
	query := fmt.Sprintf("%s = ?", key)

	if err := db.Where(query, value).Find(&entities).Error; err != nil {
		return entities, err
	}
	return entities, nil
}

func (repo *abstractRepositoryImpl[T, K]) Create(tx *gorm.DB, newEntity T) (T, error) {
	if err := repo.transCheck(tx).Create(&newEntity).Error; err != nil {
		var zeroValue T
		return zeroValue, err
	}

	return newEntity, nil
}

// Update implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) Update(tx *gorm.DB, id K, newEntity T) error {
	entity := createInstance[T]()

	if err := repo.transCheck(tx).
		Model(entity).
		Where("id = ?", id).
		Updates(&newEntity).
		Error; err != nil {
		return err
	}

	return nil
}

// Delete implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) Delete(tx *gorm.DB, id K) error {
	entity := createInstance[T]()

	if err := repo.transCheck(tx).
		Where("id = ?", id).
		Delete(entity).
		Error; err != nil {
		return err
	}
	return nil
}

// Restore implements AbstractRepository.
func (repo *abstractRepositoryImpl[T, K]) Restore(tx *gorm.DB, id K) error {
	entity := createInstance[T]()

	result := repo.transCheck(tx).
		Unscoped().
		Model(entity).
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *abstractRepositoryImpl[T, K]) GetPreloads() []string {
	return nil // Default: no preloads
}

func (repo *abstractRepositoryImpl[T, K]) GetType() string {
	tType := reflect.TypeOf(new(T)).Elem().String()
	kType := reflect.TypeOf(new(K)).Elem().String()

	return fmt.Sprintf("abstractRepositoryImpl[T: %s, K: %s]", tType, kType)
}

// Helper function createInstance dynamically creates a new instance of type T.
func createInstance[T any]() *T {
	var instance T
	return &instance
}

// Helper function applyPreloads applies preloading to a GORM query if any preloads are provided.
func applyPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// Helper function to check if it is within a transactional context to use the
// transaction or use the current repository
func (repo *abstractRepositoryImpl[T, K]) transCheck(tx *gorm.DB) *gorm.DB {
	db := tx
	if db == nil {
		db = repo.gorm
		return db
	}

	return db
}

// NewAbstractRepository creates a new instance of AbstractRepositoryImpl with the provided gormDB and self.
//
// Deprecated: use CreateRepository instead
func NewAbstractRepository[T Identifiable[K], K ID](gorm *gorm.DB, self AbstractRepository[T, K]) AbstractRepository[T, K] {
	return &abstractRepositoryImpl[T, K]{gorm: gorm, self: self}
}

// CreateRepository creates a new instance of AbstractRepositoryImpl with the provided gormDB and self.
//
// example:
//
//		type AccountRepository struct {
//			stdlib.AbstractRepository[*models.Account, uint]
//		}
//
//		func NewAccountRepository(gormDB *gorm.DB) *AccountRepository {
//		return &AccountRepository{
//			AbstractRepository: stdlib.CreateRepository(gormDB, AccountRepository{}),
//		}
//	}
func CreateRepository[T Identifiable[K], K ID](gormDB *gorm.DB, self AbstractRepository[T, K]) *abstractRepositoryImpl[T, K] {
	repo := &abstractRepositoryImpl[T, K]{
		gorm: gormDB,
		self: self,
	}
	return repo
}
