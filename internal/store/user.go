package store

import (
	"context"

	"github.com/forChin/my-project/internal/model"
	"github.com/forChin/my-project/pkg/liberror"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db}
}

func (s *UserStore) Create(
	ctx context.Context, school model.User,
) (model.User, error) {
	// TODO: add logic
	return model.User{}, nil
}

func (s *UserStore) GetAll(
	ctx context.Context, limit, offset uint64, filter model.UserFilter,
) (schools []model.User, total uint64, err error) {
	// TODO: add logic
	return nil, 0, nil
}

func (s *UserStore) Delete(ctx context.Context, id int) error {
	// TODO: add logic

	return liberror.ErrNotFound // user not found
}

func (s *UserStore) Update(
	ctx context.Context, school model.User,
) error {
	// TODO: add logic
	return nil
}
