package postgres_test

import (
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/mirasildev/note_project/storage/repo"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName:   faker.FirstName(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Email:       faker.Email(),
		ImageURL:    faker.Sentence(),
		CreatedAt:   time.Now(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func deleteUser(id int64, t *testing.T) {
	err := strg.User().Delete(id)
	require.NoError(t, err)
}

func updateUser(t *testing.T) *repo.User {
	u := createUser(t)
	user, err := strg.User().Update(&repo.User{
		ID:          u.ID,
		FirstName:   faker.FirstName(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Email:       faker.Email(),
		ImageURL:    faker.Sentence(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}
func TestCreateUser(t *testing.T) {
	user := createUser(t)
	deleteUser(user.ID, t)
	require.NotEmpty(t, user)
}

func TestGetUser(t *testing.T) {
	c := createUser(t)

	user, err := strg.User().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestUpdateUser(t *testing.T) {
	user := updateUser(t)
	require.NotEmpty(t, user)
}

func TestGetAllUser(t *testing.T) {

	users, err := strg.User().GetAllUsers(&repo.GetAllUsersParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, users)

}
