package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	testUtils "github.com/Dialosoft/src/pkg/utils/test"
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
	service          ForumService
	mockForumRepo    *MockForumRepository
	mockCategoryRepo *MockCategoryRepository
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
			expectedResult: []response.ForumResponse(nil),
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
				assert.Equal(suite.T(), len(tt.expectedResult), len(result))
				assert.Equal(suite.T(), tt.expectedResult, result)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
				assert.Nil(suite.T(), result)
			}
		})
	}
}


func (suite *ForumServiceTestSuite) TestCreateForum() {
    forumID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")

    tests := []struct {
        name          string
        newRequest    request.NewForum
        mockReturn    *models.Forum
        mockReturnErr error
        expectedResult uuid.UUID
        expectedErr    error
    }{
        {
            name: "success creating forum",
            newRequest: request.NewForum{
                Name:         testUtils.PtrToString("Test Forum"),
                Description:  testUtils.PtrToString("Test Description"),
                Type:         testUtils.PtrToString("discussion"),
                IsActive:     testUtils.PtrToBool(true),
                RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
                CategoryID:   testUtils.PtrToString("8213280e-2000-403a-b375-cdcda6488450"), // Add CategoryID
            },
            mockReturn: &models.Forum{
                ID:           forumID,
                Name:         "Test Forum",
                Description:  "Test Description",
                Type:         "discussion",
                IsActive:     true,
                RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
                CategoryID:   "8213280e-2000-403a-b375-cdcda6488450",
            },
            mockReturnErr:  nil,
            expectedResult: forumID,
            expectedErr:    nil,
        },
        {
            name: "error creating forum",
            newRequest: request.NewForum{
                Name:         testUtils.PtrToString("Test Forum"),
                Description:  testUtils.PtrToString("Test Description"),
                Type:         testUtils.PtrToString("discussion"),
                IsActive:     testUtils.PtrToBool(true),
                RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
                CategoryID:   testUtils.PtrToString("8213280e-2000-403a-b375-cdcda6488450"), // Add CategoryID
            },
            mockReturn:    nil,
            mockReturnErr: errors.New("database error"),
            expectedResult: uuid.UUID{},
            expectedErr:    errors.New("database error"),
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Reset Mock's expectations before each subtest
            suite.mockForumRepo.ExpectedCalls = nil

            // Configure the mock for the creation case
            suite.mockForumRepo.On("Create", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.Forum")).
                Return(tt.mockReturn, tt.mockReturnErr)

            result, err := suite.service.CreateForum(tt.newRequest)

            if tt.expectedErr == nil {
                assert.NoError(suite.T(), err)
                assert.Equal(suite.T(), tt.expectedResult, result)
            } else {
                assert.EqualError(suite.T(), err, tt.expectedErr.Error())
                assert.Equal(suite.T(), tt.expectedResult, result)
            }
        })
    }
}


func (suite *ForumServiceTestSuite) TestUpdateForum() {
	forumID := uuid.MustParse("16ceae31-e0be-41f5-b688-4740492e8acc")
	categoryID := uuid.MustParse("26ceae31-e0be-41f5-b688-4740492e8acc")

	tests := []struct {
		name           string
		forumID        uuid.UUID
		updateRequest  request.NewForum
		mockFindReturn *models.Forum
		mockFindErr    error
		mockUpdateErr  error
		expectedErr    error
	}{
		{
			name:    "success updating forum",
			forumID: forumID,
			updateRequest: request.NewForum{
				Name:         testUtils.PtrToString("Updated Forum"),
				Description:  testUtils.PtrToString("Updated Description"),
				Type:         testUtils.PtrToString("discussion"),
				IsActive:     testUtils.PtrToBool(true),
				RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
				CategoryID:   testUtils.PtrToString(categoryID.String()),
			},
			mockFindReturn: &models.Forum{
				ID:           forumID,
				Name:         "Original Forum",
				Description:  "Original Description",
				Type:         "information",
				IsActive:     true,
				RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
				CategoryID:   categoryID.String(),
			},
			mockFindErr:   nil,
			mockUpdateErr: nil,
			expectedErr:   nil,
		},
		{
			name:    "required field missing",
			forumID: forumID,
			updateRequest: request.NewForum{
				Description: testUtils.PtrToString("Updated Description"),
				// Name is missing
			},
			mockFindReturn: &models.Forum{
				ID:          forumID,
				Name:        "Original Forum",
				Description: "Original Description",
			},
			mockFindErr:   nil,
			mockUpdateErr: nil,
			expectedErr:   errorsUtils.ErrParameterCannotBeNull,
		},
		{
			name:    "forum not found",
			forumID: forumID,
			updateRequest: request.NewForum{
				Name:        testUtils.PtrToString("Updated Forum"),
				Description: testUtils.PtrToString("Updated Description"),
			},
			mockFindReturn: nil,
			mockFindErr:    errorsUtils.ErrNotFound,
			mockUpdateErr:  nil,
			expectedErr:    errorsUtils.ErrNotFound,
		},
		{
			name:    "error during update",
			forumID: forumID,
			updateRequest: request.NewForum{
				Name:         testUtils.PtrToString("Updated Forum"),
				Description:  testUtils.PtrToString("Updated Description"),
				Type:         testUtils.PtrToString("discussion"),
				IsActive:     testUtils.PtrToBool(true),
				RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
				CategoryID:   testUtils.PtrToString(categoryID.String()),
			},
			mockFindReturn: &models.Forum{
				ID:          forumID,
				Name:        "Original Forum",
				Description: "Original Description",
			},
			mockFindErr:   nil,
			mockUpdateErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
		{
			name:    "invalid category ID",
			forumID: forumID,
			updateRequest: request.NewForum{
				Name:         testUtils.PtrToString("Updated Forum"),
				Description:  testUtils.PtrToString("Updated Description"),
				Type:         testUtils.PtrToString("discussion"),
				IsActive:     testUtils.PtrToBool(true),
				RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
				CategoryID:   testUtils.PtrToString("invalid-uuid"),
			},
			mockFindReturn: &models.Forum{
				ID:           forumID,
				Name:         "Original Forum",
				Description:  "Original Description",
				Type:         "information",
				IsActive:     true,
				RolesAllowed: []string{"8213280e-2000-403a-b375-cdcda6488450"},
				CategoryID:   categoryID.String(),
			},
			mockFindErr:   nil,
			mockUpdateErr: nil,
			expectedErr:   errorsUtils.ErrInvalidUUID,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Reset Mock's expectations before each subtest
			suite.mockForumRepo.ExpectedCalls = nil

			// Mock FindByID call
			suite.mockForumRepo.On("FindByID", tt.forumID).Return(tt.mockFindReturn, tt.mockFindErr)

			// Just configure the update mock if a validation error is not expected
			if tt.expectedErr != errorsUtils.ErrInvalidUUID && tt.expectedErr != errorsUtils.ErrParameterCannotBeNull {
				suite.mockForumRepo.On("Update",
					mock.AnythingOfType("*gorm.DB"),
					tt.forumID,
					mock.MatchedBy(func(forum *models.Forum) bool {
						return forum.ID == tt.forumID &&
							forum.Name == *tt.updateRequest.Name &&
							forum.Description == *tt.updateRequest.Description
					}),
				).Return(tt.mockUpdateErr)
			}

			err := suite.service.UpdateForum(tt.forumID, tt.updateRequest)

			if tt.expectedErr == nil {
				assert.NoError(suite.T(), err)
			} else {
				assert.EqualError(suite.T(), err, tt.expectedErr.Error())
			}
		})
	}
}

// Create a new instance of the TestForumServiceTestSuite structure and executes the test suite
func TestForumServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ForumServiceTestSuite))
}
