// Package repository provides tools for working with db
package repository

import (
	"context"
	"database/sql"
	"medods-jwt/internal/models"
	"medods-jwt/pkg/db/postgres"
	"medods-jwt/pkg/errors"
	"medods-jwt/pkg/utils"

	"github.com/lib/pq"
)

// AuthRepository is a repository for users
type AuthRepository struct {
	ctx *context.Context
	db  *postgres.DB
}

// NewUserRepository creates a new user repository
func NewAuthRepository(ctx *context.Context, db *postgres.DB) *AuthRepository {
	return &AuthRepository{ctx: ctx, db: db}
}

// GetUser returns a user from the database
func (u *AuthRepository) GetUser(guid, publicKey *string) (*models.User, error) {
	if guid == nil && publicKey == nil {
		return nil, errors.ErrNoArguments
	}
	if guid != nil && publicKey != nil {
		return nil, errors.ErrOnlyOneArgument
	}

	var user models.User
	query := `SELECT u.guid, u.email, u.ip, t.refresh_token, t.public_key FROM public.users as u LEFT JOIN public.tokens as t ON u.guid = t.guid WHERE`

	var arg string
	var err error
	if guid != nil {
		query += ` u.guid = $1;`
		arg = *guid
	} else {
		query += ` t.public_key = $1;`
		arg = *publicKey
	}
	err = u.db.QueryRow(query, arg).Scan(&user.GUID, &user.Email, &user.IP, &user.Refresh, &user.Public)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateRefreshToken updates the refresh token in the database
func (u *AuthRepository) UpdateRefreshToken(guid, refreshToken, publicKey string) error {
	query := `UPDATE public.tokens SET refresh_token = $1, public_key = $2 WHERE guid = $3;`

	hashedRefresh, err := utils.HashString(refreshToken)
	if err != nil {
		return err
	}
	res, err := u.db.Exec(query, hashedRefresh, publicKey, guid)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errors.ErrRefreshTokenNotUnique
		}
		return err
	}
	if value, _ := res.RowsAffected(); value == 0 {
		err = u.CreateRefreshToken(guid, hashedRefresh, publicKey)
		return err
	}
	return nil
}

// CreateRefreshToken creates a new refresh token in the database
func (u *AuthRepository) CreateRefreshToken(guid, refreshToken, publicKey string) error {
	query := `INSERT INTO public.tokens (guid, refresh_token, public_key) VALUES ($1, $2, $3);`
	hashedRefresh, err := utils.HashString(refreshToken)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(query, guid, hashedRefresh, publicKey)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return errors.ErrRefreshTokenNotUnique
			}
			if pgErr.Code == "23503" {
				return errors.ErrUserNotFound
			}
		}
		return err
	}
	return err
}
