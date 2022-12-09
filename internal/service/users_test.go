package service

import (
	"context"
	"errors"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	repoMocks "github.com/zhuravlev-pe/course-watch/internal/repository/mocks"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

var errSomeDatabaseError = errors.New("some database error")
var someDate = time.Now()
var password = "uniquePassword"

func noError(t *testing.T, err error) {
	assert.NoError(t, err)
}

func TestUsersService_GetUserInfo(t *testing.T) {
	cases := map[string]struct {
		id         string
		setupMocks func(context.Context, *repoMocks.MockUsers)
		output     *GetUserInfoOutput
		checkError func(*testing.T, error)
	}{
		"success": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				u := &core.User{
					Id:               "1111111",
					Email:            "doe.h@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					Roles:            []security.Role{security.Student},
				}
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(u, nil).Times(1)
			},
			output: &GetUserInfoOutput{
				Id:               "1111111",
				Email:            "doe.h@example.com",
				FirstName:        "John",
				LastName:         "Doe",
				DisplayName:      "JohnnyD",
				RegistrationDate: someDate,
				Roles:            []security.Role{security.Student},
			},
			checkError: noError,
		},
		"not_found": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, repository.ErrNotFound).Times(1)
			},
			output: nil,
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		"db_error": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, errSomeDatabaseError).Times(1)
			},
			output: nil,
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errSomeDatabaseError)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsers := repoMocks.NewMockUsers(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := newUsersService(mockUsers, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsers)

			out, err := s.GetUserInfo(ctx, tc.id)

			assert.Equal(t, tc.output, out)
			tc.checkError(t, err)
		})
	}
}

func TestUsersService_UpdateUserInfo(t *testing.T) {
	cases := map[string]struct {
		id         string
		input      *UpdateUserInfoInput
		setupMocks func(context.Context, *repoMocks.MockUsers)
		checkError func(*testing.T, error)
	}{
		"success": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, nil).Times(1)
				var upd repository.UpdateUserInput
				upd.FirstName = "John"
				upd.LastName = "Doe"
				upd.DisplayName = "JohnnyD"
				mockUsers.EXPECT().Update(ctx, "1111111", &upd).Return(nil).Times(1)
			},
			checkError: noError,
		},
		"not_found": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, repository.ErrNotFound).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
		"db_error_on_get": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, errSomeDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errSomeDatabaseError)
			},
		},
		"db_error_on_update": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetById(ctx, "1111111").Return(nil, nil).Times(1)
				var upd repository.UpdateUserInput
				upd.FirstName = "John"
				upd.LastName = "Doe"
				upd.DisplayName = "JohnnyD"
				mockUsers.EXPECT().Update(ctx, "1111111", &upd).Return(errSomeDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errSomeDatabaseError)
			},
		},
		"validation_firstName_required": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				LastName: "Doe",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["first_name"]
				assert.True(t, ok)
			},
		},
		"validation_lastName_required": {
			id: "1111111",
			input: &UpdateUserInfoInput{
				FirstName: "John",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["last_name"]
				assert.True(t, ok)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsers := repoMocks.NewMockUsers(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := newUsersService(mockUsers, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsers)

			err = s.UpdateUserInfo(ctx, tc.id, tc.input)

			tc.checkError(t, err)
		})
	}
}

func TestUsersService_Signup(t *testing.T) {
	cases := map[string]struct {
		input *SignupUserInput
		// setupMocks func(context.Context, *repoMocks.MockUsers, *serviceMocks.MockIdGen)
		setupMocks func(context.Context, *repoMocks.MockUsers)
		checkError func(*testing.T, error)
	}{
		"success": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "John*123#",
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			// setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers, mockgenId *serviceMocks.MockIdGen) {
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(nil, nil).Times(1)
				mockUsers.EXPECT().Insert(ctx, gomock.Any()).Return(nil).Times(1)
			},
			checkError: noError,
		},
		"user_already_exist": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "John*123#",
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				u := &core.User{
					Id:               "1111111",
					Email:            "JognDoe@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					Roles:            []security.Role{security.Student},
				}
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(u, nil).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrUserAlreadyExist)
			},
		},
		"db_error_on_insert": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "John*123#",
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				u := &core.User{
					Id:               "1111111",
					Email:            "doe.h@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					Roles:            []security.Role{security.Student},
				}
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(u, nil).Times(1)
				mockUsers.EXPECT().Insert(ctx, gomock.Any()).Return(errSomeDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errSomeDatabaseError)
			},
		},
		"validation_firstName_required": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "John*123#",
				FirstName:   "",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["first_name"]
				assert.True(t, ok)
			},
		},
		"validation_lastName_required": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "John*123#",
				FirstName:   "John",
				LastName:    "",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["last_name"]
				assert.True(t, ok)
			},
		},
		"validation_weak_password": {
			input: &SignupUserInput{
				Email:       "JognDoe@example.com",
				Password:    "pw*123",
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["password"]
				assert.True(t, ok)
			},
		},
		"validation_invalid_email": {
			input: &SignupUserInput{
				Email:       "JognDoe@example",
				Password:    "John*123#",
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["email"]
				assert.True(t, ok)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsers := repoMocks.NewMockUsers(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := newUsersService(mockUsers, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsers)

			err = s.Signup(ctx, tc.input)

			tc.checkError(t, err)
		})
	}
}

func TestUsersService_Login(t *testing.T) {
	cases := map[string]struct {
		// email      string
		input      *LoginInput
		setupMocks func(context.Context, *repoMocks.MockUsers)
		checkError func(*testing.T, error)
	}{
		"success": {
			// email: "JognDoe@example.com",
			input: &LoginInput{
				Email:      "JognDoe@example.com",
				Password:   password,
				Persistent: true,
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				loggedUser := &core.User{
					Id:               "1111111",
					Email:            "doe.h@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					HashedPassword:   hashedPassword,
					Roles:            []security.Role{security.Student},
				}
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(loggedUser, nil).Times(1)
			},
			checkError: noError,
		},
		"invalid_password": {
			// email: "JognDoe@example.com",
			input: &LoginInput{
				Email:      "JognDoe@example.com",
				Password:   password,
				Persistent: true,
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+"qwerty"), bcrypt.DefaultCost)
				loggedUser := &core.User{
					Id:               "1111111",
					Email:            "doe.h@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					HashedPassword:   hashedPassword,
					Roles:            []security.Role{security.Student},
				}
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(loggedUser, nil).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		"invalid_credentials_no_user": {
			// email: "JognDoe@example.com",
			input: &LoginInput{
				Email:      "JognDoe@example.com",
				Password:   password,
				Persistent: true,
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(nil, repository.ErrNotFound).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		"db_error_on_select": {
			// email: "JognDoe@example.com",
			input: &LoginInput{
				Email:      "JognDoe@example.com",
				Password:   password,
				Persistent: true,
			},
			setupMocks: func(ctx context.Context, mockUsers *repoMocks.MockUsers) {
				mockUsers.EXPECT().GetByEmail(ctx, "JognDoe@example.com").Return(nil, errSomeDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errSomeDatabaseError)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsers := repoMocks.NewMockUsers(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := newUsersService(mockUsers, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsers)

			_, err = s.Login(ctx, tc.input)

			tc.checkError(t, err)
		})
	}
}
