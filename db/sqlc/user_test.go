package db

import (
	"context"
	"testing"
	"time"

	"github.com/jotabf/simplebank/util"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hash, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hash,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)   // Check expected error
	require.NotEmpty(t, user) // Check if empty
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt) // Check if timestamp was generated

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(),
		GetUserParams{
			Username: user1.Username,
			Email:    user1.Email,
		})

	require.NoError(t, err)    // Check expected error
	require.NotEmpty(t, user2) // Check if empty
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second) // Difference can not be greater than one second
}

// func TestUpdateUser(t *testing.T) {
// 	user1 := createRandomUser(t)

// 	arg := UpdateUserParams{
// 		ID:      user1.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	user2, err := testQueries.UpdateUser(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, user2)
// 	require.Equal(t, user1.ID, user2.ID)
// 	require.Equal(t, user1.Owner, user2.Owner)
// 	require.Equal(t, arg.Balance, user2.Balance)
// 	require.Equal(t, user1.Currency, user2.Currency)
// 	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second) // Difference can not be greater than one second
// }

// func TestDeleteUser(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	err := testQueries.DeleteUser(context.Background(), user1.ID)
// 	require.NoError(t, err)

// 	user2, err := testQueries.GetUser(context.Background(), user1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, user2)
// }

// func TestListUsers(t *testing.T) {
// 	var users1 []User
// 	for i := 0; i < 10; i++ {
// 		users1 = append(users1, createRandomUser(t))
// 	}

// 	arg := ListUsersParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	users2, err := testQueries.ListUsers(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.Len(t, users2, 5)

// 	for _, user := range users2 {
// 		require.NotEmpty(t, user)
// 	}
// }
