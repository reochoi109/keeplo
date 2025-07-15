package user

import "context"

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
}
