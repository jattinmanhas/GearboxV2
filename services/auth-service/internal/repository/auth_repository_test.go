package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	cleanup := func() {
		sqlxDB.Close()
	}

	return sqlxDB, mock, cleanup
}

func TestNewUserRepository(t *testing.T) {
	db, _, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)
	assert.NotNil(t, repo)
	assert.Implements(t, (*IUserRepository)(nil), repo)
}

func TestUserRepository_RegisterNewUser_Success(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username:    "testuser",
		Password:    "hashedpassword",
		Email:       "test@example.com",
		FirstName:   "John",
		MiddleName:  "M",
		LastName:    "Doe",
		Avatar:      "avatar.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedID := uint(1)

	// Mock the INSERT query with RETURNING clause
	mock.ExpectQuery(`
		INSERT INTO users \(
			username, password, email, first_name, middle_name, last_name, avatar, gender, date_of_birth
		\) VALUES \(
			\?, \?, \?, \?, \?, \?, \?, \?, \?
		\) RETURNING id;
	`).WithArgs(
		user.Username, user.Password, user.Email, user.FirstName,
		user.MiddleName, user.LastName, user.Avatar, user.Gender, user.DateOfBirth,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	err := repo.RegisterNewUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_RegisterNewUser_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username:    "testuser",
		Password:    "hashedpassword",
		Email:       "test@example.com",
		FirstName:   "John",
		MiddleName:  "M",
		LastName:    "Doe",
		Avatar:      "avatar.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Mock database error
	mock.ExpectQuery(`
		INSERT INTO users \(
			username, password, email, first_name, middle_name, last_name, avatar, gender, date_of_birth
		\) VALUES \(
			\?, \?, \?, \?, \?, \?, \?, \?, \?
		\) RETURNING id;
	`).WithArgs(
		user.Username, user.Password, user.Email, user.FirstName,
		user.MiddleName, user.LastName, user.Avatar, user.Gender, user.DateOfBirth,
	).WillReturnError(sql.ErrConnDone)

	err := repo.RegisterNewUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_RegisterNewUser_ScanError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user := &domain.User{
		Username:    "testuser",
		Password:    "hashedpassword",
		Email:       "test@example.com",
		FirstName:   "John",
		MiddleName:  "M",
		LastName:    "Doe",
		Avatar:      "avatar.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Mock the query but return invalid data that can't be scanned
	mock.ExpectQuery(`
		INSERT INTO users \(
			username, password, email, first_name, middle_name, last_name, avatar, gender, date_of_birth
		\) VALUES \(
			\?, \?, \?, \?, \?, \?, \?, \?, \?
		\) RETURNING id;
	`).WithArgs(
		user.Username, user.Password, user.Email, user.FirstName,
		user.MiddleName, user.LastName, user.Avatar, user.Gender, user.DateOfBirth,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("invalid_id"))

	err := repo.RegisterNewUser(context.Background(), user)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	expectedUser := &domain.User{
		ID:          1,
		Username:    "testuser",
		Password:    "hashedpassword",
		Email:       "test@example.com",
		FirstName:   "John",
		MiddleName:  "M",
		LastName:    "Doe",
		Avatar:      "avatar.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1;`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "password", "email", "first_name", "middle_name",
			"last_name", "avatar", "gender", "date_of_birth", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Email,
			expectedUser.FirstName, expectedUser.MiddleName, expectedUser.LastName, expectedUser.Avatar,
			expectedUser.Gender, expectedUser.DateOfBirth, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsDeleted,
		))

	user, err := repo.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Mock the SELECT query returning no rows
	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1;`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Mock database error
	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1;`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	user, err := repo.GetUserByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_Success(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	expectedUsers := []domain.User{
		{
			ID:          1,
			Username:    "user1",
			Password:    "hash1",
			Email:       "user1@example.com",
			FirstName:   "John",
			MiddleName:  "M",
			LastName:    "Doe",
			Avatar:      "avatar1.jpg",
			Gender:      "male",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsDeleted:   false,
		},
		{
			ID:          2,
			Username:    "user2",
			Password:    "hash2",
			Email:       "user2@example.com",
			FirstName:   "Jane",
			MiddleName:  "K",
			LastName:    "Smith",
			Avatar:      "avatar2.jpg",
			Gender:      "female",
			DateOfBirth: time.Date(1992, 5, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsDeleted:   false,
		},
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM users LIMIT \$1 OFFSET \$2;`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "password", "email", "first_name", "middle_name",
			"last_name", "avatar", "gender", "date_of_birth", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			expectedUsers[0].ID, expectedUsers[0].Username, expectedUsers[0].Password, expectedUsers[0].Email,
			expectedUsers[0].FirstName, expectedUsers[0].MiddleName, expectedUsers[0].LastName, expectedUsers[0].Avatar,
			expectedUsers[0].Gender, expectedUsers[0].DateOfBirth, expectedUsers[0].CreatedAt, expectedUsers[0].UpdatedAt, expectedUsers[0].IsDeleted,
		).AddRow(
			expectedUsers[1].ID, expectedUsers[1].Username, expectedUsers[1].Password, expectedUsers[1].Email,
			expectedUsers[1].FirstName, expectedUsers[1].MiddleName, expectedUsers[1].LastName, expectedUsers[1].Avatar,
			expectedUsers[1].Gender, expectedUsers[1].DateOfBirth, expectedUsers[1].CreatedAt, expectedUsers[1].UpdatedAt, expectedUsers[1].IsDeleted,
		))

	users, err := repo.GetAllUsers(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].Username, users[0].Username)
	assert.Equal(t, expectedUsers[1].Username, users[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_EmptyResult(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Mock the SELECT query returning no rows
	mock.ExpectQuery(`SELECT \* FROM users LIMIT \$1 OFFSET \$2;`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "password", "email", "first_name", "middle_name",
			"last_name", "avatar", "gender", "date_of_birth", "created_at", "updated_at", "is_deleted",
		}))

	users, err := repo.GetAllUsers(context.Background(), 10, 0)

	assert.NoError(t, err)
	// When no rows are returned, sqlx returns an empty slice, not nil
	assert.Len(t, users, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_DatabaseError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Mock database error
	mock.ExpectQuery(`SELECT \* FROM users LIMIT \$1 OFFSET \$2;`).
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	users, err := repo.GetAllUsers(context.Background(), 10, 0)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_WithPagination(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	expectedUser := domain.User{
		ID:          3,
		Username:    "user3",
		Password:    "hash3",
		Email:       "user3@example.com",
		FirstName:   "Bob",
		MiddleName:  "L",
		LastName:    "Johnson",
		Avatar:      "avatar3.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1988, 12, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Mock the SELECT query with pagination
	mock.ExpectQuery(`SELECT \* FROM users LIMIT \$1 OFFSET \$2;`).
		WithArgs(5, 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "username", "password", "email", "first_name", "middle_name",
			"last_name", "avatar", "gender", "date_of_birth", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Email,
			expectedUser.FirstName, expectedUser.MiddleName, expectedUser.LastName, expectedUser.Avatar,
			expectedUser.Gender, expectedUser.DateOfBirth, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.IsDeleted,
		))

	users, err := repo.GetAllUsers(context.Background(), 5, 10)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)
	assert.Equal(t, expectedUser.Username, users[0].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_ContextCancellation(t *testing.T) {
	db, _, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	user := &domain.User{
		Username:    "testuser",
		Password:    "hashedpassword",
		Email:       "test@example.com",
		FirstName:   "John",
		MiddleName:  "M",
		LastName:    "Doe",
		Avatar:      "avatar.jpg",
		Gender:      "male",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Test RegisterNewUser with cancelled context
	err := repo.RegisterNewUser(ctx, user)
	assert.Error(t, err)

	// Test GetUserByID with cancelled context
	userResult, err := repo.GetUserByID(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, userResult)

	// Test GetAllUsers with cancelled context
	users, err := repo.GetAllUsers(ctx, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, users)
}
