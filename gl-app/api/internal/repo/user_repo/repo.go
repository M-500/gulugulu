package user_repo

import "context"

type UserRepo interface {
	FindOneByEmail(ctx context.Context, email string) (*UserModel, error)
}

type userRepoImpl struct {
	dao UserDao
}

func (u *userRepoImpl) FindOneByEmail(ctx context.Context, email string) (*UserModel, error) {
	return u.dao.FindOneByEmail(ctx, email)
}
