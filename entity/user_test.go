package entity_test

import (
	"testing"

	"github.com/denniswon/gowiki/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u, err := entity.NewUser("sjobs@apple.com", "new_password", "Steve", "Jobs")
	assert.Nil(t, err)
	assert.Equal(t, u.FirstName, "Steve")
	assert.NotNil(t, u.ID)
	assert.NotEqual(t, u.Password, "new_password")
}

func TestValidatePassword(t *testing.T) {
	u, _ := entity.NewUser("sjobs@apple.com", "new_password", "Steve", "Jobs")
	err := u.ValidatePassword("new_password")
	assert.Nil(t, err)
	err = u.ValidatePassword("wrong_password")
	assert.NotNil(t, err)

}

func TestUserValidate(t *testing.T) {
	type test struct {
		email     string
		password  string
		firstName string
		lastName  string
		want      error
	}

	tests := []test{
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      nil,
		},
		{
			email:     "",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      entity.ErrInvalidEntity,
		},
		{
			email:     "sjobs@apple.com",
			password:  "",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      nil,
		},
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "",
			lastName:  "Jobs",
			want:      entity.ErrInvalidEntity,
		},
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "",
			want:      entity.ErrInvalidEntity,
		},
	}
	for _, tc := range tests {

		_, err := entity.NewUser(tc.email, tc.password, tc.firstName, tc.lastName)
		assert.Equal(t, err, tc.want)
	}

}
