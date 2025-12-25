package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/techies/streamify/internal/database"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type BaseService struct {
	DB *database.Queries
}

func NewBaseService(db *database.Queries) BaseService {
	return BaseService{DB: db}
}
