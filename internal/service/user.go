package service

import (
	"context"

	"github.com/forChin/my-project/internal/model"
	"github.com/forChin/my-project/internal/store"
)

type UserService struct {
	userStore *store.UserStore
}

func NewUserService(
	userStore *store.UserStore,
) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) Create(ctx context.Context, school model.User) (model.User, error) {
	return s.userStore.Create(ctx, school)
}

func (s *UserService) GetAll(
	ctx context.Context, limit, offset uint64, filter model.UserFilter,
) (schools []model.User, total uint64, err error) {
	return s.userStore.GetAll(ctx, limit, offset, filter)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.userStore.Delete(ctx, id)
}

func (s *UserService) Update(ctx context.Context, school model.User) error {
	return s.userStore.Update(ctx, school)
}
