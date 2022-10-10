package storage

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/carfloresf/financial-chat/internal/model"
)

func TestDB_CreateUser(t *testing.T) {

	type args struct {
		username     string
		passwordHash string
	}

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		want        int
		wantErr     bool
	}{
		{"success-create-user",
			func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(insertUserQuery))).
					WithArgs("test", "hash", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args{
				username:     "test",
				passwordHash: "hash",
			},
			1,
			false,
		},
		{"fail-create-user",
			func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(insertUserQuery))).
					WithArgs("test2", "hash", sqlmock.AnyArg()).
					WillReturnError(fmt.Errorf("error"))
			},
			args{
				username:     "test2",
				passwordHash: "hash",
			},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			var storage DB
			storage.Conn = db

			tt.mockClosure(mock)

			got, err := storage.CreateUser(tt.args.username, tt.args.passwordHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDB_GetUser(t *testing.T) {
	type args struct {
		username string
	}
	now := time.Now()
	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		want        *model.User
		wantErr     bool
	}{
		{"success-get-user",
			func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(getUserQuery))).
					WithArgs("test-username").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at"}).
						AddRow(1, "test-username", "test-hash", now))
			},
			args{"test-username"},
			&model.User{ID: 1, Username: "test-username", PasswordHash: "test-hash", CreatedAt: now},
			false},
		{"fail-get-user",
			func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(getUserQuery))).
					WithArgs("fail-username").
					WillReturnError(fmt.Errorf("error"))
			},
			args{"fail-username"},
			nil,
			true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			var storage DB
			storage.Conn = db

			tt.mockClosure(mock)

			got, err := storage.GetUser(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
