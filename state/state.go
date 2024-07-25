package state;

import (
	"context"
	"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
)

var (
	Logger    *zap.Logger
	Context   = context.Background()
	Validator = validator.New()
)