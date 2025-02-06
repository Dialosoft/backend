package services

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/repository"
	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/errorsUtils"
	// testUtils "github.com/Dialosoft/src/pkg/utils/test"
)

type MockPostRepository struct {
	repository.MockAbstractRepository[*models.Post, uuid.UUID]
}

func (m *MockPostRepository) FindAllWithPagination(limit, offset int) ([]*models.Post, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Post), args.Error(1)
}

func (m *MockPostRepository) FindAllByForumID(forumID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	args := m.Called(forumID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Post), args.Error(1)
}

func (m *MockPostRepository) FindByUserID(userID uuid.UUID) ([]*models.Post, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Post), args.Error(1)
}

func (m *MockPostRepository) GetLikeCount(postID uuid.UUID) (int64, error) {
	args := m.Called(postID)
	return args.Get(0).(int64), args.Error(1)
}

type MockPostLikesRepository struct {
	repository.MockAbstractRepository[*models.PostLikes, uuid.UUID]
}

func (m *MockPostLikesRepository) FindAllByPostID(postID uuid.UUID) ([]*models.PostLikes, error) {
	args := m.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PostLikes), args.Error(1)
}

func (m *MockPostLikesRepository) FindAllByUserIDAndPostID(postID uuid.UUID, userID uuid.UUID) ([]*models.PostLikes, error) {
	args := m.Called(postID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PostLikes), args.Error(1)
}

func (m *MockPostLikesRepository) FindAllByUserID(userID uuid.UUID) ([]*models.PostLikes, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PostLikes), args.Error(1)
}

func (m *MockPostLikesRepository) SaveLike(postID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(postID, userID)
	return args.Error(0)
}

func (m *MockPostLikesRepository) RemoveLike(postID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(postID, userID)
	return args.Error(0)
}

type PostServiceTestSuite struct {
	suite.Suite
	service       PostService
	mockPostRepo  *MockPostRepository
	mockPostLikes *MockPostLikesRepository
	mockUserRepo  *MockUserRepository
	sampleUserID  uuid.UUID
	sampleForumID uuid.UUID
	samplePostID  uuid.UUID
	sampleUser    *models.UserEntity
	samplePost    *models.Post
}

func (suite *PostServiceTestSuite) SetupTest() {
	suite.mockPostRepo = new(MockPostRepository)
	suite.mockPostLikes = new(MockPostLikesRepository)
	suite.mockUserRepo = new(MockUserRepository)
	suite.service = NewPostService(suite.mockPostRepo, suite.mockPostLikes, suite.mockUserRepo)

	// Initialize sample data
	suite.sampleUserID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	suite.sampleForumID = uuid.MustParse("660e8400-e29b-41d4-a716-446655440000")
	suite.samplePostID = uuid.MustParse("770e8400-e29b-41d4-a716-446655440000")

	suite.sampleUser = &models.UserEntity{
		ID:       suite.sampleUserID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.samplePost = &models.Post{
		ID:        suite.samplePostID,
		UserID:    suite.sampleUserID,
		ForumID:   suite.sampleForumID,
		Title:     "Test Post",
		Content:   "Test Content",
		User:      *suite.sampleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (suite *PostServiceTestSuite) TearDownTest() {
	suite.mockPostRepo.AssertExpectations(suite.T())
	suite.mockPostLikes.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreateNewPost() {
	tests := []struct {
		name          string
		input         request.NewPost
		userID        uuid.UUID
		mockSetup     func()
		expectedError error
	}{
		{
			name: "success - create post with valid data",
			input: request.NewPost{
				Title:   "Test Post",
				Content: "Test Content",
				ForumID: suite.sampleForumID.String(),
			},
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", suite.sampleUserID).
					Return(suite.sampleUser, nil)

				suite.mockPostRepo.On("Create", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.Post")).
					Return(suite.samplePost, nil)
			},
			expectedError: nil,
		},
		{
			name: "error - user not found",
			input: request.NewPost{
				Title:   "Test Post",
				Content: "Test Content",
				ForumID: suite.sampleForumID.String(),
			},
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", suite.sampleUserID).
					Return((*models.UserEntity)(nil), errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
		},
		{
			name: "error - invalid forum ID",
			input: request.NewPost{
				Title:   "Test Post",
				Content: "Test Content",
				ForumID: "invalid-uuid",
			},
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockUserRepo.On("FindByID", suite.sampleUserID).
					Return(suite.sampleUser, nil)
			},
			expectedError: errors.New("invalid UUID length: 12"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			result, err := suite.service.CreateNewPost(tt.userID, tt.input)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), suite.samplePost.ID, result.ID)
				assert.Equal(suite.T(), tt.input.Title, result.Title)
				assert.Equal(suite.T(), tt.input.Content, result.Content)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestLikePost() {
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success - like post",
			postID: suite.samplePostID,
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("FindAllByUserIDAndPostID", suite.samplePostID, suite.sampleUserID).
					Return([]*models.PostLikes{}, nil)
				suite.mockPostLikes.On("SaveLike", suite.samplePostID, suite.sampleUserID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error - already liked",
			postID: suite.samplePostID,
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("FindAllByUserIDAndPostID", suite.samplePostID, suite.sampleUserID).
					Return([]*models.PostLikes{{PostID: suite.samplePostID, UserID: suite.sampleUserID}}, nil)
			},
			expectedError: errorsUtils.ErrPostAlreadyLiked,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			err := suite.service.LikePost(tt.postID, tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestUnlikePost() {
	tests := []struct {
		name          string
		postID        uuid.UUID
		userID        uuid.UUID
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success - unlike post",
			postID: suite.samplePostID,
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("RemoveLike", suite.samplePostID, suite.sampleUserID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error - failed to unlike post",
			postID: suite.samplePostID,
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("RemoveLike", suite.samplePostID, suite.sampleUserID).
					Return(errors.New("failed to remove like"))
			},
			expectedError: errors.New("failed to remove like"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			err := suite.service.UnlikePost(tt.postID, tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestGetAllPostsByForum() {
	samplePosts := []*models.Post{suite.samplePost}

	tests := []struct {
		name          string
		forumID       uuid.UUID
		limit         int
		offset        int
		mockSetup     func()
		expectedPosts []response.PostResponse
		expectedError error
	}{
		{
			name:    "success - get posts with valid forum ID",
			forumID: suite.sampleForumID,
			limit:   10,
			offset:  0,
			mockSetup: func() {
				suite.mockPostRepo.On("FindAllByForumID", suite.sampleForumID, 10, 0).
					Return(samplePosts, nil)
			},
			expectedPosts: []response.PostResponse{
				{
					ID:      suite.samplePostID,
					Title:   suite.samplePost.Title,
					Content: suite.samplePost.Content,
					User: response.UserResponse{
						ID:       suite.sampleUserID,
						Username: suite.sampleUser.Username,
						Email:    suite.sampleUser.Email,
					},
					CreatedAt: suite.samplePost.CreatedAt,
					UpdatedAt: suite.samplePost.UpdatedAt,
				},
			},
			expectedError: nil,
		},
		{
			name:    "error - database error",
			forumID: suite.sampleForumID,
			limit:   10,
			offset:  0,
			mockSetup: func() {
				suite.mockPostRepo.On("FindAllByForumID", suite.sampleForumID, 10, 0).
					Return(nil, errors.New("database error"))
			},
			expectedPosts: nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			posts, err := suite.service.GetAllPostsByForum(tt.forumID, tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedPosts, posts)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestGetAllPostsAndReturnSimpleResponse() {
	samplePosts := []*models.Post{suite.samplePost}

	tests := []struct {
		name           string
		limit          int
		offset         int
		mockSetup      func()
		expectedPosts  []response.SimplePostResponse
		expectedError  error
	}{
		{
			name:   "success - get posts with pagination",
			limit:  10,
			offset: 0,
			mockSetup: func() {
				suite.mockPostRepo.On("FindAllWithPagination", 10, 0).
					Return(samplePosts, nil)
			},
			expectedPosts: []response.SimplePostResponse{
				{
					ID:        suite.samplePostID.String(),
					UserID:    suite.sampleUserID.String(),
					Title:     suite.samplePost.Title,
					CreatedAt: suite.samplePost.CreatedAt.String(),
					UpdatedAt: suite.samplePost.UpdatedAt.String(),
					DeletedAt: suite.samplePost.DeletedAt,
				},
			},
			expectedError: nil,
		},
		{
			name:   "error - database error",
			limit:  10,
			offset: 0,
			mockSetup: func() {
				suite.mockPostRepo.On("FindAllWithPagination", 10, 0).
					Return(nil, errors.New("database error"))
			},
			expectedPosts: nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			posts, err := suite.service.GetAllPostsAndReturnSimpleResponse(tt.limit, tt.offset)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
				assert.Nil(suite.T(), posts)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedPosts, posts)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestGetLikeCount() {
	tests := []struct {
		name           string
		postID        uuid.UUID
		mockSetup      func()
		expectedCount  int64
		expectedError  error
	}{
		{
			name:   "success - get like count",
			postID: suite.samplePostID,
			mockSetup: func() {
				suite.mockPostRepo.On("GetLikeCount", suite.samplePostID).
					Return(int64(5), nil)
			},
			expectedCount: 5,
			expectedError: nil,
		},
		{
			name:   "error - database error",
			postID: suite.samplePostID,
			mockSetup: func() {
				suite.mockPostRepo.On("GetLikeCount", suite.samplePostID).
					Return(int64(0), errors.New("database error"))
			},
			expectedCount: 0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			count, err := suite.service.GetLikeCount(tt.postID)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
				assert.Equal(suite.T(), int64(0), count)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedCount, count)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestUpdatePostTitle() {
	tests := []struct {
		name          string
		postID        uuid.UUID
		title         string
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "success - update post title",
			postID: suite.samplePostID,
			title:  "New Title",
			mockSetup: func() {
				suite.mockPostRepo.On("FindByID", suite.samplePostID).
					Return(suite.samplePost, nil)
				
				updatedPost := *suite.samplePost
				updatedPost.Title = "New Title"
				
				suite.mockPostRepo.On("Update", mock.Anything, suite.samplePostID, &updatedPost).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error - post not found",
			postID: suite.samplePostID,
			title:  "New Title",
			mockSetup: func() {
				var nilPost *models.Post
				suite.mockPostRepo.On("FindByID", suite.samplePostID).
					Return(nilPost, gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "error - update failed",
			postID: suite.samplePostID,
			title:  "New Title",
			mockSetup: func() {
				suite.mockPostRepo.On("FindByID", suite.samplePostID).
					Return(suite.samplePost, nil)
				
				updatedPost := *suite.samplePost
				updatedPost.Title = "New Title"
				
				suite.mockPostRepo.On("Update", mock.Anything, suite.samplePostID, &updatedPost).
					Return(errors.New("update failed"))
			},
			expectedError: errors.New("update failed"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			err := suite.service.UpdatePostTitle(tt.postID, tt.title)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func (suite *PostServiceTestSuite) TestGetPostLikesByUserID() {
	samplePostLikes := &models.PostLikes{
		PostID: suite.samplePostID,
		UserID: suite.sampleUserID,
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func()
		expectedIDs   []uuid.UUID
		expectedError error
	}{
		{
			name:   "success - get liked posts",
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("FindAllByUserID", suite.sampleUserID).
					Return([]*models.PostLikes{samplePostLikes}, nil)
			},
			expectedIDs:   []uuid.UUID{suite.samplePostID},
			expectedError: nil,
		},
		{
			name:   "error - database error",
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("FindAllByUserID", suite.sampleUserID).
					Return(nil, errors.New("database error"))
			},
			expectedIDs:   nil,
			expectedError: errors.New("database error"),
		},
		{
			name:   "error - no likes found",
			userID: suite.sampleUserID,
			mockSetup: func() {
				suite.mockPostLikes.On("FindAllByUserID", suite.sampleUserID).
					Return([]*models.PostLikes{}, nil)
			},
			expectedIDs:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			tt.mockSetup()

			ids, err := suite.service.GetPostLikesByUserID(tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(suite.T(), err, tt.expectedError.Error())
				assert.Nil(suite.T(), ids)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedIDs, ids)
			}
		})
	}
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}
