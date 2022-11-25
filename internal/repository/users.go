package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

type UsersRepo struct {
	client *pgxpool.Pool
}

func NewUsersRepo(client *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{client: client}
}

func (u *UsersRepo) Insert(ctx context.Context, user *core.User) error {
	query := `
		INSERT INTO public.users
		    (id, email, firstname, lastname, display_name,
		     registration_date, hashed_password, roles)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8);
		`
	
	_, err := u.client.Exec(ctx, query, user.Id, user.Email, user.FirstName, user.LastName,
		user.DisplayName, user.RegistrationDate, user.HashedPassword, user.Roles)
	
	return err
}

func (u *UsersRepo) Update(ctx context.Context, id string, input *UpdateUserInput) error {
	query := `
		UPDATE public.users
		  SET (firstname, lastname, display_name) = ($1, $2, $3)
          WHERE id = $4
		`
	
	_, err := u.client.Exec(ctx, query, input.FirstName, input.LastName, input.DisplayName, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepo) GetById(ctx context.Context, id string) (*core.User, error) {
	query := `
		SELECT id, email, firstname, lastname, display_name,
		       registration_date, hashed_password, roles
		FROM public.users
		WHERE id = $1;
		`
	
	return u.getByField(ctx, query, id)
}

func (u *UsersRepo) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	query := `
		SELECT id, email, firstname, lastname, display_name,
		       registration_date, hashed_password, roles
		FROM public.users
		WHERE email = $1
		`
	
	return u.getByField(ctx, query, email)
}

func (u *UsersRepo) getByField(ctx context.Context, query string, field string) (*core.User, error) {
	var user core.User
	var r []uint8
	
	//TODO to think of a better way of scanning/storing []Role
	err := u.client.QueryRow(ctx, query, field).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DisplayName,
		&user.RegistrationDate,
		&user.HashedPassword,
		&r,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	user.Roles = security.ToRoles(r)
	
	return &user, nil
}
