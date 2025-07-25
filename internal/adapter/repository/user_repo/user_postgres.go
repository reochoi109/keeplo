package user_repo

import (
	"context"
	"keeplo/internal/domain/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGorm struct {
	ID           uuid.UUID `gorm:"column:id"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	IsActive     bool      `gorm:"column:is_active"`
	IsDeleted    bool      `gorm:"column:is_deleted"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type GormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepo(db *gorm.DB) user.Repository {
	return &GormUserRepo{db: db}
}

func (UserGorm) TableName() string {
	return "users"
}

func (r *GormUserRepo) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(toGorm(u)).Error
}

func (r *GormUserRepo) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).
		Model(&UserGorm{}).
		Where("id = ?", u.ID).
		Updates(toGorm(u)).Error
}

// 탈퇴한 사용자 포함하여 이메일 중복 여부 확인
func (r *GormUserRepo) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&UserGorm{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u UserGorm
	if err := r.db.WithContext(ctx).
		Where("email = ? AND is_deleted = false", email).
		First(&u).Error; err != nil {
		return nil, err
	}
	return toEntity(&u), nil
}

func (r *GormUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u UserGorm
	if err := r.db.WithContext(ctx).
		Where("id = ? AND is_deleted = false", id).
		First(&u).Error; err != nil {
		return nil, err
	}
	return toEntity(&u), nil
}

func (r *GormUserRepo) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&UserGorm{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_at": time.Now(),
		}).Error
}

func (r *GormUserRepo) HardDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&UserGorm{}).Error
}

func toEntity(u *UserGorm) *user.User {
	return &user.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		IsDeleted:    u.IsDeleted,
	}
}

func toGorm(u *user.User) *UserGorm {
	return &UserGorm{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		IsDeleted:    u.IsDeleted,
	}
}
