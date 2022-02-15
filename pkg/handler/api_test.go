package handler

import (
	"bytes"
	mock_pkg "for_avito_tech_with_gin/pkg/mocks"
	"for_avito_tech_with_gin/pkg/service"
	mock_service "for_avito_tech_with_gin/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserBehavior func(s *mock_service.MockUser)
type mockCalculatorBehavior func(s *mock_pkg.MockCurrencyCalculator)

type testSkillet struct {
	name                   string
	inputBody              string
	inputQueryParams       string
	mockUserBehavior       mockUserBehavior
	mockCalculatorBehavior mockCalculatorBehavior
	expectedStatusCode     int
	expectedRequestBody    string
}

func TestHandler_addFundsHandler(t *testing.T) {
	testData := []testSkillet{
		{
			name:      "OK",
			inputBody: `{"id":348, "sum": 2700}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().AddFunds(348, float32(2700)).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "",
		},
		{
			name:                "Invalid Body",
			inputBody:           `{"id":348}`,
			mockUserBehavior:    func(s *mock_service.MockUser) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid body."}`,
		},
		{
			name:      "Negative Sum",
			inputBody: `{"id":34, "sum": -10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().AddFunds(34, float32(-10)).Return(&service.NegativeSum{})
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"sum can't be negative or 0."}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"id":14589, "sum": 10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().AddFunds(14589, float32(10)).Return(&service.InternalServerError{})
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"internal server error."}`,
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servi := mock_service.NewMockUser(c)
			testCase.mockUserBehavior(servi)

			services := &service.Service{User: servi}
			handler := NewHandler(services)

			// test server
			r := gin.New()
			r.POST("/api/v1/add_funds", handler.addFundsHandler)

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/add_funds", bytes.NewBufferString(testCase.inputBody))

			// perform request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_writeOffFundsHandler(t *testing.T) {
	testData := []testSkillet{
		{
			name:      "OK",
			inputBody: `{"id":348, "sum": 2700}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().WriteOffFunds(348, float32(2700)).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "",
		},
		{
			name:                "Invalid Body",
			inputBody:           `{"id":348}`,
			mockUserBehavior:    func(s *mock_service.MockUser) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid body."}`,
		},
		{
			name:      "Negative Sum",
			inputBody: `{"id":34, "sum": -10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().WriteOffFunds(34, float32(-10)).Return(&service.NegativeSum{})
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"sum can't be negative or 0."}`,
		},
		{
			name:      "User Not Found",
			inputBody: `{"id":91, "sum": 10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().WriteOffFunds(91, float32(10)).Return(&service.UserNotFound{Id: 91})
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"user 91 does not exist."}`,
		},
		{
			name:      "Insufficient Funds",
			inputBody: `{"id":23, "sum": 10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().WriteOffFunds(23, float32(10)).Return(&service.InsufficientFunds{Id: 23})
			},
			expectedStatusCode:  http.StatusPreconditionFailed,
			expectedRequestBody: `{"message":"user 23 has insufficient funds."}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"id":14589, "sum": 10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().WriteOffFunds(14589, float32(10)).Return(&service.InternalServerError{})
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"internal server error."}`,
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servi := mock_service.NewMockUser(c)
			testCase.mockUserBehavior(servi)

			services := &service.Service{User: servi}
			handler := NewHandler(services)

			// test server
			r := gin.New()
			r.POST("/api/v1/write_off_funds", handler.writeOffFundsHandler)

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/write_off_funds", bytes.NewBufferString(testCase.inputBody))

			// perform request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_fundsTransferHandler(t *testing.T) {
	testData := []testSkillet{
		{
			name:      "OK",
			inputBody: `{"sender_id":348, "receiver_id": 4389, "sum": 2700}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(348, 4389, float32(2700)).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "",
		},
		{
			name:                "Invalid Body",
			inputBody:           `{"sender_id":348, "receiver_id": 4389}`,
			mockUserBehavior:    func(s *mock_service.MockUser) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid body."}`,
		},
		{
			name:      "Negative Sum",
			inputBody: `{"sender_id":34, "receiver_id": 89, "sum": -10}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(34, 89, float32(-10)).Return(&service.NegativeSum{})
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"sum can't be negative or 0."}`,
		},
		{
			name:      "Equal Sender And Receiver",
			inputBody: `{"sender_id":34, "receiver_id": 34, "sum": 1000}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(34, 34, float32(1000)).Return(&service.SameId{})
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"user cannot send money to himself."}`,
		},
		{
			name:      "User Not Found",
			inputBody: `{"sender_id":91, "receiver_id": 12, "sum": 599}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(91, 12, float32(599)).Return(&service.UserNotFound{Id: 91})
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"user 91 does not exist."}`,
		},
		{
			name:      "Insufficient Funds",
			inputBody: `{"sender_id":23, "receiver_id": 24, "sum": 1000}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(23, 24, float32(1000)).Return(&service.InsufficientFunds{Id: 23})
			},
			expectedStatusCode:  http.StatusPreconditionFailed,
			expectedRequestBody: `{"message":"user 23 has insufficient funds."}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"sender_id":14589, "receiver_id": 4389, "sum": 3500}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().FundsTransfer(14589, 4389, float32(3500)).Return(&service.InternalServerError{})
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"internal server error."}`,
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servi := mock_service.NewMockUser(c)
			testCase.mockUserBehavior(servi)

			services := &service.Service{User: servi}
			handler := NewHandler(services)

			// test server
			r := gin.New()
			r.POST("/api/v1/funds_transfer", handler.fundsTransferHandler)

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/funds_transfer", bytes.NewBufferString(testCase.inputBody))

			// perform request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getBalance(t *testing.T) {
	testData := []testSkillet{
		{
			name:      "OK",
			inputBody: `{"id":348}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetBalance(348).Return(float32(100), nil)
			},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {},
			expectedStatusCode:     http.StatusOK,
			expectedRequestBody:    `{"balance":100}`,
		},
		{
			name:                   "Invalid Body",
			inputBody:              `{}`,
			mockUserBehavior:       func(s *mock_service.MockUser) {},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {},
			expectedStatusCode:     http.StatusBadRequest,
			expectedRequestBody:    `{"message":"invalid body."}`,
		},
		{
			name:             "OK With Query Param",
			inputBody:        `{"id":34}`,
			inputQueryParams: "?currency=USD",
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetBalance(34).Return(float32(100), nil)
			},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {
				s.EXPECT().ConvertRubTo("USD", float32(100)).Return(1.3, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"balance":1.3}`,
		},
		{
			name:             "Invalid Query Param",
			inputBody:        `{"id":34}`,
			inputQueryParams: "?currency=XRP",
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetBalance(34).Return(float32(100), nil)
			},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {
				s.EXPECT().ConvertRubTo("XRP", float32(100)).Return(0.0, &service.WrongParam{Param: "currency"})
			},
			expectedStatusCode:  http.StatusPreconditionFailed,
			expectedRequestBody: `{"message":"wrong currency param."}`,
		},
		{
			name:      "Internal Server Error",
			inputBody: `{"id":14589}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetBalance(14589).Return(float32(0), &service.InternalServerError{})
			},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {},
			expectedStatusCode:     http.StatusInternalServerError,
			expectedRequestBody:    `{"message":"internal server error."}`,
		},
		{
			name:      "User Not Found",
			inputBody: `{"id":91}`,
			mockUserBehavior: func(s *mock_service.MockUser) {
				s.EXPECT().GetBalance(91).Return(float32(0), &service.UserNotFound{Id: 91})
			},
			mockCalculatorBehavior: func(s *mock_pkg.MockCurrencyCalculator) {},
			expectedStatusCode:     http.StatusNotFound,
			expectedRequestBody:    `{"message":"user 91 does not exist."}`,
		},
	}

	t.Parallel()
	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// init deps
			c := gomock.NewController(t)
			defer c.Finish()

			userService := mock_service.NewMockUser(c)
			testCase.mockUserBehavior(userService)

			services := &service.Service{User: userService}
			handler := NewHandler(services)

			calculator := mock_pkg.NewMockCurrencyCalculator(c)
			testCase.mockCalculatorBehavior(calculator)

			// test server
			r := gin.New()
			r.GET("/api/v1/get_balance", handler.getBalanceHandler(calculator))

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/get_balance"+testCase.inputQueryParams,
				bytes.NewBufferString(testCase.inputBody))

			// perform request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
