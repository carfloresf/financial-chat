package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword_GenerateHash(t *testing.T) {
	type fields struct {
		pepper string
	}

	type args struct {
		password string
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		compareResult error
	}{
		{"success-generate-hash", fields{"pepper"}, args{"test"}, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &PasswordGenerator{
				pepper: tt.fields.pepper,
			}
			got, err := u.GenerateHash(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := u.CompareHashAndPassword(got, tt.args.password)
			assert.Equal(t, tt.compareResult, result)
		})
	}
}
