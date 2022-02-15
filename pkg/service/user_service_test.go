package service

import (
	"for_avito_tech_with_gin/pkg/model"
	"for_avito_tech_with_gin/pkg/repository"
	mock_repository "for_avito_tech_with_gin/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepositoryBehavior func(s *mock_repository.MockUser)

func TestUserService_AddFunds(t *testing.T) {
	testData := []struct {
		name                   string
		userId                 int
		sum                    float32
		mockRepositoryBehavior mockRepositoryBehavior
		expectedError          error
	}{
		{
			name:   "OK When Exist",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().UpdateBalance(17, float32(5000)).Return(&model.User{}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "OK When Not Exist",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, nil)
				s.EXPECT().CreateUser(17, float32(5000)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:                   "Incorrect Sum",
			userId:                 17,
			sum:                    0,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {},
			expectedError:          &NegativeSum{},
		},
		{
			name:   "Error in IsUserExist",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:   "Error in CreateUser",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, nil)
				s.EXPECT().CreateUser(17, float32(5000)).Return(errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:   "Error in UpdateBalance",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().UpdateBalance(17, float32(5000)).Return(nil, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockUser(c)
			testCase.mockRepositoryBehavior(repo)

			services := NewUserService(&repository.Repository{User: repo})

			// test
			err := services.AddFunds(testCase.userId, testCase.sum)

			// assert
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUserService_WriteOffFunds(t *testing.T) {
	testData := []struct {
		name                   string
		userId                 int
		sum                    float32
		mockRepositoryBehavior mockRepositoryBehavior
		expectedError          error
	}{
		{
			name:   "OK",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 20000}, nil)
				s.EXPECT().UpdateBalance(17, float32(-5000)).Return(&model.User{}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "User Not Exist",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, nil)
			},
			expectedError: &UserNotFound{Id: 17},
		},
		{
			name:                   "Incorrect Sum",
			userId:                 17,
			sum:                    -900,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {},
			expectedError:          &NegativeSum{},
		},
		{
			name:   "Error in IsUserExist",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:   "Error in GetUser",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(nil, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:   "Error in UpdateBalance",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 20000}, nil)
				s.EXPECT().UpdateBalance(17, float32(-5000)).Return(nil, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:   "Insufficient sum",
			userId: 17,
			sum:    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 300}, nil)
			},
			expectedError: &InsufficientFunds{Id: 17},
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockUser(c)
			testCase.mockRepositoryBehavior(repo)

			services := NewUserService(&repository.Repository{User: repo})

			// test
			err := services.WriteOffFunds(testCase.userId, testCase.sum)

			// assert
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUserService_FundsTransfer(t *testing.T) {
	testData := []struct {
		name                   string
		senderId               int
		receiverId             int
		sum                    float32
		mockRepositoryBehavior mockRepositoryBehavior
		expectedError          error
	}{
		{
			name:       "OK 1",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(true, nil)
				s.EXPECT().CreateFundsTransaction(17, 18, float32(5000)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "OK 2",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(false, nil)
				s.EXPECT().CreateUser(18, float32(0)).Return(nil)
				s.EXPECT().CreateFundsTransaction(17, 18, float32(5000)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:                   "Same User",
			senderId:               17,
			receiverId:             17,
			sum:                    5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {},
			expectedError:          &SameId{},
		},
		{
			name:                   "Incorrect Sum",
			senderId:               17,
			receiverId:             18,
			sum:                    -900,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {},
			expectedError:          &NegativeSum{},
		},
		{
			name:       "Error in IsUserExist 1",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:       "Sender User Not Exist",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, nil)
			},
			expectedError: &UserNotFound{Id: 17},
		},
		{
			name:       "Error in GetUser 1",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(nil, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:       "Sender Doesnt Have Enough Money",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 300}, nil)
			},
			expectedError: &InsufficientFunds{Id: 17},
		},
		{
			name:       "Error in IsUserExist 2",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(false, errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:       "Error in CreateUser",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(false, nil)
				s.EXPECT().CreateUser(18, float32(0)).Return(errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:       "Error in CreateFundsTransaction 1",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(false, nil)
				s.EXPECT().CreateUser(18, float32(0)).Return(nil)
				s.EXPECT().CreateFundsTransaction(17, 18, float32(5000)).Return(errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
		{
			name:       "Error in CreateFundsTransaction 2",
			senderId:   17,
			receiverId: 18,
			sum:        5000,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 30000}, nil)
				s.EXPECT().IsUserExist(18).Return(true, nil)
				s.EXPECT().CreateFundsTransaction(17, 18, float32(5000)).Return(errors.Errorf("lol kek cheburek."))
			},
			expectedError: &InternalServerError{},
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockUser(c)
			testCase.mockRepositoryBehavior(repo)

			services := NewUserService(&repository.Repository{User: repo})

			// test
			err := services.FundsTransfer(testCase.senderId, testCase.receiverId, testCase.sum)

			// assert
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUserService_GetBalance(t *testing.T) {
	testData := []struct {
		name                   string
		userId                 int
		mockRepositoryBehavior mockRepositoryBehavior
		expectedBalance        float32
		expectedError          error
	}{
		{
			name:   "OK",
			userId: 17,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(&model.User{Id: 17, UserId: 17, Balance: 10000}, nil)
			},
			expectedBalance: 10000,
			expectedError:   nil,
		},
		{
			name:   "Error in IsUserExist",
			userId: 17,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, errors.Errorf("lol kek cheburek."))
			},
			expectedBalance: 0,
			expectedError:   &InternalServerError{},
		},
		{
			name:   "User Not Found",
			userId: 17,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(false, nil)
			},
			expectedBalance: 0,
			expectedError:   &UserNotFound{Id: 17},
		},
		{
			name:   "Error in GetUser",
			userId: 17,
			mockRepositoryBehavior: func(s *mock_repository.MockUser) {
				s.EXPECT().IsUserExist(17).Return(true, nil)
				s.EXPECT().GetUser(17).Return(nil, errors.Errorf("lol kek cheburek."))
			},
			expectedBalance: 0,
			expectedError:   &InternalServerError{},
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockUser(c)
			testCase.mockRepositoryBehavior(repo)

			services := NewUserService(&repository.Repository{User: repo})

			// test
			balance, err := services.GetBalance(testCase.userId)

			// assert
			assert.Equal(t, testCase.expectedBalance, balance)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
