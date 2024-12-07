// Package repository provides tools for working with db
package repository

import (
	"context"
	"medods-jwt/pkg/errors"
	"medods-jwt/pkg/logger"

	"go.uber.org/zap"
)

// EmailRepository is a repository for email
type EmailRepository struct {
	ctx *context.Context
}

// NewEmailRepository creates a new email repository
func NewEmailRepository(ctx *context.Context) *EmailRepository {
	return &EmailRepository{ctx: ctx}
}

// SendWarning sends warning
func (r *EmailRepository) SendWarning(email, userIP, currentIP string) error {
	logger.GetLoggerFromCtx(*r.ctx).Info("Send warning", zap.String("email", email), zap.String("userIP", userIP), zap.String("currentIP", currentIP))
	if false {
		return errors.ErrSendWarning
	}
	return nil
}
