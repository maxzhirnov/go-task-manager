package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVerificationToken(t *testing.T) {
	token1, err1 := GenerateVerificationToken()
	token2, err2 := GenerateVerificationToken()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2) // Tokens should be unique
	assert.Len(t, token1, 64)          // 32 bytes in hex = 64 characters
}

func TestCreateVerificationToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	now := time.Now()

	mock.ExpectQuery("INSERT INTO verification_tokens").
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	token, err := CreateVerificationToken(db, userID)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, userID, token.UserID)
	assert.NotEmpty(t, token.Token)
	assert.True(t, token.ExpiresAt.After(now))
}

func TestGetVerificationTokenForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		userID  int
		mockSQL func()
		want    string
		wantErr bool
	}{
		{
			name:   "Valid token exists",
			userID: 1,
			mockSQL: func() {
				mock.ExpectQuery("SELECT token FROM verification_tokens").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"token"}).AddRow("valid-token"))
			},
			want:    "valid-token",
			wantErr: false,
		},
		{
			name:   "No token found",
			userID: 2,
			mockSQL: func() {
				mock.ExpectQuery("SELECT token FROM verification_tokens").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			got, err := GetVerificationTokenForUser(db, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		token   string
		mockSQL func()
		wantErr bool
	}{
		{
			name:  "Successful verification",
			token: "valid-token",
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE verification_tokens").
					WithArgs("valid-token").
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				mock.ExpectExec("UPDATE users SET is_verified").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Invalid token",
			token: "invalid-token",
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE verification_tokens").
					WithArgs("invalid-token").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			err := VerifyEmail(db, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	validUser := &User{
		Email:    "test@example.com",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
	}

	tests := []struct {
		name    string
		user    *User
		mockSQL func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful creation",
			user: validUser,
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						false, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("INSERT INTO verification_tokens").
					WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Missing email",
			user:    &User{Password: "$2a$10$hashedpassword"},
			mockSQL: func() {},
			wantErr: true,
			errMsg:  "email and password are required",
		},
		{
			name: "Unhashed password",
			user: &User{
				Email:    "test@example.com",
				Password: "plaintext",
			},
			mockSQL: func() {},
			wantErr: true,
			errMsg:  "password must be hashed before saving",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			err := tt.user.CreateUser(db)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResendVerificationToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		userID  int
		mockSQL func()
		wantErr bool
		errMsg  string
	}{
		{
			name:   "Successful resend",
			userID: 1,
			mockSQL: func() {
				// Check if user is verified
				mock.ExpectQuery("SELECT is_verified FROM users").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"is_verified"}).AddRow(false))

				// Deactivate old tokens
				mock.ExpectExec("UPDATE verification_tokens").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Create new token
				mock.ExpectQuery("INSERT INTO verification_tokens").
					WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name:   "Already verified user",
			userID: 2,
			mockSQL: func() {
				mock.ExpectQuery("SELECT is_verified FROM users").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"is_verified"}).AddRow(true))
			},
			wantErr: true,
			errMsg:  "user is already verified",
		},
		{
			name:   "User not found",
			userID: 3,
			mockSQL: func() {
				mock.ExpectQuery("SELECT is_verified FROM users").
					WithArgs(3).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			token, err := ResendVerificationToken(db, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, tt.userID, token.UserID)
				assert.NotEmpty(t, token.Token)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	tests := []struct {
		name    string
		email   string
		mockSQL func()
		want    User
		wantErr bool
	}{
		{
			name:  "User exists",
			email: "test@example.com",
			mockSQL: func() {
				rows := sqlmock.NewRows([]string{
					"id", "email", "username", "password", "is_verified", "created_at", "updated_at",
				}).AddRow(
					1, "test@example.com", "testuser", "hashedpass", true, now, now,
				)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: User{
				ID:         1,
				Email:      "test@example.com",
				Username:   "testuser",
				Password:   "hashedpass",
				IsVerified: true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			wantErr: false,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			mockSQL: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			got, err := GetUserByEmail(db, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// func TestGetUserStatistics(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	tests := []struct {
// 		name    string
// 		userID  int
// 		mockSQL func()
// 		want    *UserStatistics
// 		wantErr bool
// 	}{
// 		{
// 			name:   "Statistics exist",
// 			userID: 1,
// 			mockSQL: func() {
// 				rows := sqlmock.NewRows([]string{
// 					"user_id", "username", "total_tasks", "completed_tasks",
// 					"pending_tasks", "in_progress_tasks", "deleted_tasks",
// 					"tasks_created_today",
// 				}).AddRow(1, "testuser", 10, 5, 3, 2, 0, 1)
// 				mock.ExpectQuery("SELECT (.+) FROM user_statistics").
// 					WithArgs(1).
// 					WillReturnRows(rows)
// 			},
// 			want: &UserStatistics{
// 				UserID:            1,
// 				Username:          "testuser",
// 				TotalTasks:        10,
// 				CompletedTasks:    5,
// 				PendingTasks:      3,
// 				InProgressTasks:   2,
// 				DeletedTasks:      0,
// 				TasksCreatedToday: 1,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "Statistics not found",
// 			userID: 2,
// 			mockSQL: func() {
// 				mock.ExpectQuery("SELECT (.+) FROM user_statistics").
// 					WithArgs(2).
// 					WillReturnError(sql.ErrNoRows)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSQL()
// 			got, err := GetUserStatistics(db, tt.userID)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want, got)
// 			}
// 		})
// 	}
// }

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.password}
			err := u.HashPassword()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, len(u.Password) > 0)
				assert.NotEqual(t, tt.password, u.Password)
				assert.True(t, len(u.Password) > 50)
			}
		})
	}
}
