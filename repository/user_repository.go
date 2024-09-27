package repository

import (
	"context"
	"sysbitBroker/pkg"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	CreateOrUpdate(c context.Context, user CreateOrUpdateUserParams) (CreateOrUpdateUserRow, error)
	UpdateRefreshToken(c context.Context, token UpdateUserTokenParams) error
	GetUser(c context.Context, userID int64) (User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pg *pkg.Postgres) IUserRepository {
	return &userRepository{
		db: pg.Pool,
	}
}

type CreateOrUpdateUserRow struct {
	UserID       int64
	AppID        string
	Email        string
	Picture      pgtype.Text
	Role         []string
	IsActive     bool
	Name         pgtype.Text
	RefreshToken string
	IsDeleted    bool
	CreatedAt    int64
	UpdatedAt    int64
	Operation    string
}
type CreateOrUpdateUserParams struct {
	Userid       int64
	Appid        string
	Mail         string
	Pic          pgtype.Text
	Roles        []string
	Name         pgtype.Text
	Refreshtoken string
	Isdeleted    bool
	Isactive     bool
	Createdat    int64
	Updatedat    int64
	Picture      pgtype.Text
}

func (rp *userRepository) CreateOrUpdate(c context.Context, arg CreateOrUpdateUserParams) (CreateOrUpdateUserRow, error) {

	const createOrUpdateUser = `-- name: CreateOrUpdateUser :one
INSERT INTO users (
  user_id,
	app_id,
	email,
  picture,
	role,
	name,
	refresh_token,
	is_deleted,
	is_active,
	created_at,
	updated_at
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
) ON CONFLICT (email) DO UPDATE SET 
picture = $12,
name = $6
RETURNING user_id, app_id, email, picture, role, is_active, name, refresh_token, is_deleted, created_at, updated_at, CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END as operation
`
	row := rp.db.QueryRow(c, createOrUpdateUser,
		arg.Userid,
		arg.Appid,
		arg.Mail,
		arg.Pic,
		arg.Roles,
		arg.Name,
		arg.Refreshtoken,
		arg.Isdeleted,
		arg.Isactive,
		arg.Createdat,
		arg.Updatedat,
		arg.Picture,
	)
	var i CreateOrUpdateUserRow
	err := row.Scan(
		&i.UserID,
		&i.AppID,
		&i.Email,
		&i.Picture,
		&i.Role,
		&i.IsActive,
		&i.Name,
		&i.RefreshToken,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Operation,
	)
	return i, err
}

type UpdateUserTokenParams struct {
	UserID       int64
	RefreshToken string
}

func (rp *userRepository) UpdateRefreshToken(c context.Context, arg UpdateUserTokenParams) error {

	const updateUserToken = `-- name: UpdateUserToken :exec
	UPDATE users SET refresh_token = $2 WHERE user_id = $1
	`
	_, err := rp.db.Exec(c, updateUserToken, arg.UserID, arg.RefreshToken)
	return err
}

type User struct {
	UserID       int64
	AppID        string
	Email        string
	Picture      pgtype.Text
	Role         []string
	IsActive     bool
	Name         pgtype.Text
	RefreshToken string
	IsDeleted    bool
	CreatedAt    int64
	UpdatedAt    int64
}

func (rp *userRepository) GetUser(c context.Context, userID int64) (User, error) {
	const getUser = `-- name: GetUser :one
SELECT user_id, app_id, email, picture, role, is_active, name, refresh_token, is_deleted, created_at, updated_at FROM users
WHERE user_id = $1 LIMIT 1
`
	row := rp.db.QueryRow(c, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.AppID,
		&i.Email,
		&i.Picture,
		&i.Role,
		&i.IsActive,
		&i.Name,
		&i.RefreshToken,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
