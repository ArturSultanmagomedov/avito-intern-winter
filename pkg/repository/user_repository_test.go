package repository

import (
	"fmt"
	"for_avito_tech_with_gin/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		userId  int
		balance float32
	}

	testData := []struct {
		name             string
		args             args
		mockSqlxBehavior func(args args)
		wantError        bool
	}{
		{
			name: "OK",
			args: args{
				userId:  71,
				balance: 1000,
			},
			mockSqlxBehavior: func(args args) {
				mock.ExpectQuery(`insert into users \(user_id, balance\) values \(\$1, \$2\);`).WithArgs(args.userId, args.balance).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantError: false,
		},
		{
			name: "ERR",
			args: args{
				userId:  71,
				balance: 1000,
			},
			mockSqlxBehavior: func(args args) {
				mock.ExpectQuery(`insert into users \(user_id, balance\) values \(\$1, \$2\);`).WithArgs(args.userId, args.balance).
					WillReturnError(fmt.Errorf("error"))
			},
			wantError: true,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockSqlxBehavior(testCase.args)

			err := repo.CreateUser(testCase.args.userId, testCase.args.balance)

			// assert
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		userId int
	}

	testData := []struct {
		name             string
		args             args
		mockSqlxBehavior func(args args, user model.User)
		expectedUser     model.User
		wantError        bool
	}{
		{
			name: "OK",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args, user model.User) {
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1`).WithArgs(args.userId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
						AddRow(user.Id, user.UserId, user.Balance))
			},
			expectedUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 2000,
			},
			wantError: false,
		},
		{
			name: "ERR",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args, user model.User) {
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1`).WithArgs(args.userId).
					WillReturnError(fmt.Errorf("error"))
			},
			expectedUser: model.User{},
			wantError:    true,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockSqlxBehavior(testCase.args, testCase.expectedUser)

			user, err := repo.GetUser(testCase.args.userId)

			// assert
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedUser, *user)
			}
		})
	}
}

func TestUserRepository_IsUserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		userId int
	}

	testData := []struct {
		name             string
		args             args
		mockSqlxBehavior func(args args)
		expectedResult   bool
		wantError        bool
	}{
		{
			name: "OK True",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args) {
				mock.ExpectQuery(`select count\(1\) from users where user_id = \$1`).WithArgs(args.userId).
					WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(1))
			},
			expectedResult: true,
			wantError:      false,
		},
		{
			name: "OK False",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args) {
				mock.ExpectQuery(`select count\(1\) from users where user_id = \$1`).WithArgs(args.userId).
					WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))
			},
			expectedResult: false,
			wantError:      false,
		},
		{
			name: "ERR",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args) {
				mock.ExpectQuery(`select count\(1\) from users where user_id = \$1`).WithArgs(args.userId).
					WillReturnError(fmt.Errorf("error"))
			},
			wantError: true,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockSqlxBehavior(testCase.args)

			ex, err := repo.IsUserExist(testCase.args.userId)

			// assert
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, ex)
			}
		})
	}
}

func TestUserRepository_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		userId int
		sum    float32
	}

	testData := []struct {
		name             string
		args             args
		originalUser     model.User
		mockSqlxBehavior func(args args, origUser model.User, exUser model.User)
		expectedUser     model.User
		wantError        bool
	}{
		{
			name: "OK +",
			args: args{
				userId: 71,
				sum:    20,
			},
			mockSqlxBehavior: func(args args, origUser model.User, exUser model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.userId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(origUser.Id, origUser.UserId, origUser.Balance))
				mock.ExpectQuery(`update users set balance = \$1 where user_id = \$2;`).
					WithArgs(origUser.Balance+args.sum, args.userId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(exUser.Id, exUser.UserId, exUser.Balance))
				mock.ExpectCommit()
			},
			originalUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 80,
			},
			expectedUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 100,
			},
			wantError: false,
		},
		{
			name: "OK -",
			args: args{
				userId: 71,
				sum:    -20,
			},
			mockSqlxBehavior: func(args args, origUser model.User, exUser model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.userId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(origUser.Id, origUser.UserId, origUser.Balance))
				mock.ExpectQuery(`update users set balance = \$1 where user_id = \$2;`).
					WithArgs(origUser.Balance+args.sum, args.userId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(exUser.Id, exUser.UserId, exUser.Balance))
				mock.ExpectCommit()
			},
			originalUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 80,
			},
			expectedUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 60,
			},
			wantError: false,
		},
		{
			name: "Error in Begin",
			mockSqlxBehavior: func(args args, user model.User, exUser model.User) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			wantError: true,
		},
		{
			name: "ERR bad args",
			mockSqlxBehavior: func(args args, origUser model.User, exUser model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(nil).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(origUser.Id, origUser.UserId, origUser.Balance).RowError(1, fmt.Errorf("some error")))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Get 1",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args, origUser model.User, exUser model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.userId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Get 2",
			args: args{
				userId: 71,
			},
			mockSqlxBehavior: func(args args, origUser model.User, exUser model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.userId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(origUser.Id, origUser.UserId, origUser.Balance))
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.userId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockSqlxBehavior(testCase.args, testCase.originalUser, testCase.expectedUser)

			user, err := repo.UpdateBalance(testCase.args.userId, testCase.args.sum)

			// assert
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedUser, *user)
			}
		})
	}
}

func TestUserRepository_CreateFundsTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		senderId   int
		receiverId int
		sum        float32
	}

	testData := []struct {
		name             string
		args             args
		senderUser       model.User
		receiverUser     model.User
		mockSqlxBehavior func(args args, sender, receiver model.User)
		wantError        bool
	}{
		{
			name: "OK +",
			args: args{
				senderId:   71,
				receiverId: 56,
			},
			senderUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 2000,
			},
			receiverUser: model.User{
				Id:      56,
				UserId:  56,
				Balance: 100,
			},
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.senderId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(sender.Id, sender.UserId, sender.Balance))
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.receiverId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(receiver.Id, receiver.UserId, receiver.Balance))
				mock.ExpectExec(`update users set balance = $1 where user_id = $2;`).
					WithArgs(sender.Balance-args.sum, args.senderId).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectExec(`update users set balance = $1 where user_id = $2;`).
					WithArgs(receiver.Balance+args.sum, args.receiverId).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			wantError: true,
		},
		{
			name: "Error in Begin",
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			wantError: true,
		},
		{
			name: "ERR bad args",
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(nil).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(sender.Id, sender.UserId, sender.Balance).RowError(1, fmt.Errorf("some error")))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Get 1",
			args: args{
				senderId: 71,
			},
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.senderId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Get 2",
			args: args{
				senderId:   71,
				receiverId: 56,
			},
			senderUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 2000,
			},
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.senderId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(sender.Id, sender.UserId, sender.Balance))
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.receiverId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Exec 1",
			args: args{
				senderId:   71,
				receiverId: 56,
			},
			senderUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 2000,
			},
			receiverUser: model.User{
				Id:      56,
				UserId:  56,
				Balance: 100,
			},
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.senderId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(sender.Id, sender.UserId, sender.Balance))
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.receiverId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(receiver.Id, receiver.UserId, receiver.Balance))
				mock.ExpectExec(`update users set balance = $1 where user_id = $2;`).
					WithArgs(sender.Balance-args.sum, args.senderId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "Error in Exec 2",
			args: args{
				senderId:   71,
				receiverId: 56,
			},
			senderUser: model.User{
				Id:      71,
				UserId:  71,
				Balance: 2000,
			},
			receiverUser: model.User{
				Id:      56,
				UserId:  56,
				Balance: 100,
			},
			mockSqlxBehavior: func(args args, sender, receiver model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.senderId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(sender.Id, sender.UserId, sender.Balance))
				mock.ExpectQuery(`select id, user_id, balance from users where user_id = \$1;`).
					WithArgs(args.receiverId).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).
					AddRow(receiver.Id, receiver.UserId, receiver.Balance))
				mock.ExpectExec(`update users set balance = $1 where user_id = $2;`).
					WithArgs(sender.Balance-args.sum, args.senderId).WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectExec(`update users set balance = $1 where user_id = $2;`).
					WithArgs(receiver.Balance+args.sum, args.receiverId).WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockSqlxBehavior(testCase.args, testCase.senderUser, testCase.receiverUser)

			err := repo.CreateFundsTransaction(testCase.args.senderId, testCase.args.receiverId, testCase.args.sum)

			// assert
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
