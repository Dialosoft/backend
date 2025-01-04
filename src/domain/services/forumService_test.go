package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
)

// MockForumRepository mocks the ForumRepository interface
type MockForumRepository struct {
	repository.MockAbstractRepository[*models.Forum, uuid.UUID]
}

func (m *MockForumRepository) FindByName(name string) (*models.Forum, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Forum), args.Error(1)
}

func (m *MockForumRepository) FindAllWithDeleted() ([]*models.Forum, error) {
	args := m.Called()
	return args.Get(0).([]*models.Forum), args.Error(1)
}

func (m *MockForumRepository) FindByIDWithDeleted(uuid uuid.UUID) (*models.Forum, error) {
	args := m.Called(uuid)
	return args.Get(0).(*models.Forum), args.Error(1)
}

func (m *MockForumRepository) FindAllByCategoryID(categoryID uuid.UUID) ([]*models.Forum, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]*models.Forum), args.Error(1)
}

func (m *MockForumRepository) UpdateCategoryOwner(id uuid.UUID, categoryID uuid.UUID) error {
	args := m.Called(id, categoryID)
	return args.Error(0)
}

// ForumServiceTestSuite contains the test suite for ForumService
type ForumServiceTestSuite struct {
	suite.Suite
	service            ForumService
	mockForumRepo      *MockForumRepository
	mockCategoryRepo   *MockCategoryRepository
}

func (suite *ForumServiceTestSuite) SetupTest() {
	suite.mockForumRepo = new(MockForumRepository)
	suite.mockCategoryRepo = new(MockCategoryRepository)
	suite.service = NewForumService(suite.mockForumRepo, suite.mockCategoryRepo)
}

func (suite *ForumServiceTestSuite) TearDownTest() {
	suite.mockForumRepo.AssertExpectations(suite.T())
	suite.mockCategoryRepo.AssertExpectations(suite.T())
}


func (suite *ForumServiceTestSuite) TestGetAllForums() {
	forumID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		mockReturn     []*models.Forum
		mockReturnErr  error
		expectedResult []response.ForumResponse
		expectedErr    error
	}{
		{
			name: "success retrieving all forums",
			mockReturn: []*models.Forum{
				{
					ID:          forumID,
					Name:        "General Discussion",
					Description: "Discuss anything here",
				},
			},
			mockReturnErr: nil,
			expectedResult: []response.ForumResponse{
				mapper.ForumEntityToForumResponse(&models.Forum{
					ID:          forumID,
					Name:        "General Discussion",
					Description: "Discuss anything here",
				}),
			},
			expectedErr: nil,
		},
		{
			name:           "error retrieving forums",
			mockReturn:     nil,
			mockReturnErr:  errors.New("database error"),
			expectedResult: nil,
			expectedErr:    errors.New("database error"),
		},
		{
			name:           "empty list of forums",
			mockReturn:     []*models.Forum{},
			mockReturnErr:  nil,
			expectedResult: []response.ForumResponse{},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockForumRepo.ExpectedCalls = nil

			suite.mockForumRepo.On("FindAll").Return(tt.mockReturn, tt.mockReturnErr)

			result, err := suite.service.GetAllForums()

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), len(tt.expectedResult), len(result)) // Verificar cantidad
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}