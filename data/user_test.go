package data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Find(t *testing.T) {
	LoadFixtures(t)

	repo := &UserRepository{DB: client}
	user, err := repo.Find(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "2024-03-28 08:00:00", user.CreatedAt.Format(time.DateTime))
}

func TestUserRepository_Select(t *testing.T) {
	LoadFixtures(t)

	repo := &UserRepository{DB: client}
	users, err := repo.Select(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestUserRepository_Create(t *testing.T) {
	LoadFixtures(t)

	repo := &UserRepository{DB: client}
	assert.NoError(t, repo.Create(context.Background(), User{
		Name:      "t",
		CreatedAt: time.Date(2024, 03, 28, 10, 10, 10, 0, time.Local),
	}))

	users, err := repo.Select(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 4)
}
