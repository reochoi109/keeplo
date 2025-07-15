package user

import (
	"context"
	"fmt"
	"keeplo/internal/domain/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 12

type Service interface {
	RegisterUser(ctx context.Context, email, password string) (*user.User, error)
	LoginUser(ctx context.Context, email, password string) (*user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	ResignUser(ctx context.Context, id string) error
	CheckPassword(ctx context.Context, id, password string) error
	UpdateNickname(ctx context.Context, id, nickname string) error
	UpdatePassword(ctx context.Context, id, currentPassword, newPassword string) error
	CheckDuplicateEmail(ctx context.Context, email string) (bool, error)
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, email, password string) (*user.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if len(email) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("input error, email or password is empty")
	}

	exist, err := s.repo.IsEmailExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("already exists email address")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newUser := &user.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		IsActive:     true,
		IsDeleted:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *service) LoginUser(ctx context.Context, email, password string) (*user.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if u.IsDeleted || !u.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("password or email is wrong")
	}
	return u, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ResignUser(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) CheckPassword(ctx context.Context, id, password string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (s *service) UpdateNickname(ctx context.Context, id, nickname string) error {
	if len(strings.TrimSpace(nickname)) == 0 {
		return fmt.Errorf("닉네임을 입력하세요")
	}
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u.IsDeleted || !u.IsActive {
		return fmt.Errorf("계정을 사용할 수 없습니다")
	}
	u.NickName = strings.TrimSpace(nickname)
	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *service) UpdatePassword(ctx context.Context, id, currentPassword, newPassword string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u.IsDeleted || !u.IsActive {
		return fmt.Errorf("계정을 사용할 수 없습니다")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("현재 비밀번호가 일치하지 않습니다")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), passwordCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if u.IsDeleted {
		return fmt.Errorf("이미 탈퇴된 계정입니다")
	}

	u.IsDeleted = true
	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *service) CheckDuplicateEmail(ctx context.Context, email string) (bool, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	return s.repo.IsEmailExists(ctx, email)
}
