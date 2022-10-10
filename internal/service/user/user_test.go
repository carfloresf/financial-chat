package user

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/carfloresf/financial-chat/internal/model"
	mockuser "github.com/carfloresf/financial-chat/internal/service/user/mock"
	"github.com/carfloresf/financial-chat/internal/storage"
	mockstorage "github.com/carfloresf/financial-chat/internal/storage/mock"
)

func TestUser_Create(t *testing.T) {
	type fields struct {
		db storage.Storage
		up PasswordHasher
	}

	type args struct {
		user model.User
	}

	ctrl := gomock.NewController(t)
	mockDB := mockstorage.NewMockstorage(ctrl)
	mockDB.EXPECT().CreateUser("test", "hash").Return(1, nil)
	mockDB.EXPECT().CreateUser("carlos", "hash").Return(0, fmt.Errorf("error"))

	passwordMock := mockuser.NewMockPasswordHasher(ctrl)
	passwordMock.EXPECT().GenerateHash("test").Return("hash", nil).AnyTimes().AnyTimes()
	passwordMock.EXPECT().GenerateHash("carlos").Return("", fmt.Errorf("error")).AnyTimes()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{"success-create-user",
			fields{
				mockDB, passwordMock,
			},
			args{
				model.User{Username: "test", Password: "test"},
			},
			false,
			nil,
		},
		{"fail-create-user",
			fields{
				mockDB, passwordMock,
			},
			args{
				model.User{Username: "carlos", Password: "test"},
			},
			true,
			fmt.Errorf("error creating user: error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				db: tt.fields.db,
				ph: tt.fields.up,
			}

			err := u.Register(tt.args.user.Username, tt.args.user.Password)
			if err != nil && !tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
