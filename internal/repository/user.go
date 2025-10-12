package repository

import (
	"context"
	"errors"
	"fmt"
	"users-profile-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	NotFoundUserError = errors.New("users not found")
)

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *User {
	return &User{db: db}
}

type Field int

const (
	Id Field = iota
	Email
	Login
	ACCESS_EMAIL_STATUS
)

func (user *User) getBy(ctx context.Context, field Field, value interface{}) (*models.User, error) {
	var fieldName string
	switch field {
	case Id:
		fieldName = "id"
	case Email:
		fieldName = "email"
	case Login:
		fieldName = "login"
	default:
		return nil, fmt.Errorf("unknown field %v", field)
	}

	query := fmt.Sprintf(
		`
		SELECT id, login, email, access_email_status, role, is_deleted, delete_at FROM users_profile.users
		WHERE %s = $1;
	`, fieldName)

	var u models.User
	err := user.db.QueryRow(ctx, query, value).Scan(
		&u.Id,
		&u.Login,
		&u.Email,
		&u.AccessEmailStatus,
		&u.Role,
		&u.IsDeleted,
		&u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, NotFoundUserError
	}
	return &u, err
}

func (user *User) GetById(ctx context.Context, id int64) (*models.User, error) {
	return user.getBy(ctx, Id, id)
}

func (user *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return user.getBy(ctx, Email, email)
}

func (user *User) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	return user.getBy(ctx, Login, login)
}

func (user *User) Save(ctx context.Context, userToCreate *models.User) (int64, error) {
	query := `
		INSERT INTO users_profile.users (login, email, access_email_status, role, is_deleted, delete_at) VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id;
	`
	err := user.db.QueryRow(ctx, query,
		userToCreate.Login,
		userToCreate.Email,
		userToCreate.AccessEmailStatus,
		userToCreate.Role,
		userToCreate.IsDeleted,
		userToCreate.DeletedAt,
	).Scan(&userToCreate.Id)
	if err != nil {
		return 0, err
	}
	return userToCreate.Id, nil
}

func (user *User) updateField(ctx context.Context, id int64, setField Field, setValue interface{}) error {
	if setField == Id {
		return fmt.Errorf("cannot update id field")
	}

	var setFieldName string

	switch setField {
	case Email:
		setFieldName = "email"
	case Login:
		setFieldName = "login"
	case ACCESS_EMAIL_STATUS:
		setFieldName = "access_email_status"
	default:
		return fmt.Errorf("unknown set field %v", setField)
	}

	query := fmt.Sprintf(
		"UPDATE users_profile.users SET %s = $1, updated_at = NOW() WHERE id = $2",
		setFieldName,
	)

	_, err := user.db.Exec(ctx, query, setValue, id)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) UpdateEmail(ctx context.Context, id int64, newEmail string) error {
	return user.updateField(ctx, id, Email, newEmail)
}

func (user *User) UpdateLogin(ctx context.Context, id int64, newLogin string) error {
	return user.updateField(ctx, id, Login, newLogin)
}

func (user *User) UpdateAccessEmailStatus(ctx context.Context, id int64, newStatus string) error {
	return user.updateField(ctx, id, ACCESS_EMAIL_STATUS, newStatus)
}

func (user *User) Update(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users_profile.users  SET 
		                                login=$1, 
		                                email=$2, 
		                                access_email_status=$3 
		                                role=$4 
		                            is_deleted=$5 
		                            deleted_at=$6 
		                            WHERE id=$7;
	`

	_, err := user.db.Exec(ctx, query,
		u.Login,
		u.Email,
		u.AccessEmailStatus,
		u.Role,
		u.IsDeleted,
		u.DeletedAt,
		u.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update users for id=%d: %w", u.Id, err)
	}
	return nil
}
