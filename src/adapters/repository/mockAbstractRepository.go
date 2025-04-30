package repository

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAbstractRepository is a mock implementation of the AbstractRepository interface
type MockAbstractRepository[T Identifiable[K], K ID] struct {
	mock.Mock
}

func (m *MockAbstractRepository[T, K]) FindAll() ([]T, error) {
	args := m.Called()
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockAbstractRepository[T, K]) FindByID(id K) (T, error) {
	args := m.Called(id)
	return args.Get(0).(T), args.Error(1)
}

func (m *MockAbstractRepository[T, K]) FirstByKey(key, value string) (T, error) {
	args := m.Called(key, value)
	return args.Get(0).(T), args.Error(1)
}

func (m *MockAbstractRepository[T, K]) FindAllByKey(key, value string) ([]T, error) {
	args := m.Called(key, value)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockAbstractRepository[T, K]) Create(tx *gorm.DB, newEntity T) (T, error) {
	args := m.Called(tx, newEntity)
	return args.Get(0).(T), args.Error(1)
}

func (m *MockAbstractRepository[T, K]) Update(tx *gorm.DB, id K, newEntity T) error {
	args := m.Called(tx, id, newEntity)
	return args.Error(0)
}

func (m *MockAbstractRepository[T, K]) Delete(tx *gorm.DB, id K) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockAbstractRepository[T, K]) Restore(tx *gorm.DB, id K) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockAbstractRepository[T, K]) GetPreloads() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockAbstractRepository[T, K]) GetType() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAbstractRepository[T, K]) GetKeyIdName() string {
    args := m.Called()
    return args.String(0)
}
