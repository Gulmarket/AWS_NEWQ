package repository

import (
    "github.com/gulmarket/model"
    "github.com/jmoiron/sqlx"
)

type pGAwsRepository struct {
    DB *sqlx.DB
}

func NewAwsRepository(db *sqlx.DB) model.AWSRepository {
    return &pGAwsRepository{
        DB: db,
    }
}
