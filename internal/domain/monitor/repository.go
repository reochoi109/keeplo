package monitor

import "context"

type Repository interface {
	Create(ctx context.Context, m *Monitor) error
	Update(ctx context.Context, m *Monitor) error
	FindByUserID(ctx context.Context, userID string) ([]*Monitor, error)
	FindByID(ctx context.Context, id string) (*Monitor, error)
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error

	WithTx(ctx context.Context, fn func(txRepo Repository) error) error
}
