package user

import (
	"context"
	"errors"
	"keeplo/internal/domain/user"
	"keeplo/pkg/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const passwordCost = 12
const userTimeout = time.Second * 5

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
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	email = strings.TrimSpace(strings.ToLower(email))

	if len(email) == 0 || len(password) == 0 {
		log.Warn("RegisterUser - email or password is empty")
		return nil, user.ErrInvalidCredentials
	}

	exist, err := s.repo.IsEmailExists(ctx, email)
	if err != nil {
		log.Error("RegisterUser - failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exist {
		log.Warn("RegisterUser - email already exists", zap.String("email", email))
		return nil, user.ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		log.Error("RegisterUser - failed to hash password", zap.Error(err))
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
		log.Error("RegisterUser - failed to create user", zap.Error(err))
		return nil, err
	}

	log.Info("RegisterUser - user created", zap.String("user_id", newUser.ID.String()), zap.String("email", email))
	return newUser, nil
}

func (s *service) LoginUser(ctx context.Context, email, password string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		log.Error("LoginUser - user not found", zap.Error(err))
		return nil, err
	}

	if u.IsDeleted || !u.IsActive {
		log.Warn("LoginUser - user inactive or deleted", zap.String("email", email))
		return nil, user.ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		log.Warn("LoginUser - invalid password", zap.String("email", email))
		return nil, user.ErrInvalidCredentials
	}

	log.Info("LoginUser - success", zap.String("user_id", u.ID.String()), zap.String("email", email))
	return u, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("FindByID - user not found", zap.String("user_id", id))
			return nil, user.ErrUserNotFound
		}
		log.Error("FindByID - failed", zap.String("user_id", id), zap.Error(err))
		return nil, err
	}
	return u, nil
}

func (s *service) ResignUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("ResignUser - user not found", zap.String("user_id", id))
			return user.ErrUserNotFound
		}
		log.Error("ResignUser - failed", zap.String("user_id", id), zap.Error(err))
		return err
	}
	log.Info("ResignUser - success", zap.String("user_id", id))
	return nil
}

func (s *service) CheckPassword(ctx context.Context, id, password string) error {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("CheckPassword - user not found", zap.String("user_id", id))
			return user.ErrUserNotFound
		}
		log.Error("CheckPassword - failed to get user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		log.Warn("CheckPassword - mismatch", zap.String("user_id", id))
		return user.ErrPasswordMismatch
	}

	return nil
}

func (s *service) UpdateNickname(ctx context.Context, id, nickname string) error {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)

	nickname = strings.TrimSpace(nickname)
	if len(nickname) == 0 {
		log.Warn("UpdateNickname - nickname is empty", zap.String("user_id", id))
		return user.ErrNicknameRequired
	}

	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("UpdateNickname - user not found", zap.String("user_id", id))
			return user.ErrUserNotFound
		}
		log.Error("UpdateNickname - failed to get user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	if u.IsDeleted || !u.IsActive {
		log.Warn("UpdateNickname - user inactive", zap.String("user_id", id))
		return user.ErrInactiveAccount
	}

	u.NickName = nickname
	u.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, u); err != nil {
		log.Error("UpdateNickname - update failed", zap.String("user_id", id), zap.Error(err))
		return err
	}

	log.Info("UpdateNickname - success", zap.String("user_id", id))
	return nil
}

func (s *service) UpdatePassword(ctx context.Context, id, currentPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("UpdatePassword - user not found", zap.String("user_id", id))
			return user.ErrUserNotFound
		}
		log.Error("UpdatePassword - failed to get user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	if u.IsDeleted || !u.IsActive {
		log.Warn("UpdatePassword - user inactive", zap.String("user_id", id))
		return user.ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		log.Warn("UpdatePassword - current password mismatch", zap.String("user_id", id))
		return user.ErrPasswordMismatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), passwordCost)
	if err != nil {
		log.Error("UpdatePassword - hashing failed", zap.String("user_id", id), zap.Error(err))
		return err
	}

	u.PasswordHash = string(hash)
	u.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, u); err != nil {
		log.Error("UpdatePassword - update failed", zap.String("user_id", id), zap.Error(err))
		return err
	}

	log.Info("UpdatePassword - success", zap.String("user_id", id))
	return nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("DeleteUser - user not found", zap.String("user_id", id))
			return user.ErrUserNotFound
		}
		log.Error("DeleteUser - failed to get user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	if u.IsDeleted {
		log.Warn("DeleteUser - already deleted", zap.String("user_id", id))
		return user.ErrAlreadyDeleted
	}

	u.IsDeleted = true
	u.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, u); err != nil {
		log.Error("DeleteUser - update failed", zap.String("user_id", id), zap.Error(err))
		return err
	}

	log.Info("DeleteUser - success", zap.String("user_id", id))
	return nil
}

func (s *service) CheckDuplicateEmail(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, userTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	email = strings.TrimSpace(strings.ToLower(email))
	exists, err := s.repo.IsEmailExists(ctx, email)
	if err != nil {
		log.Error("CheckDuplicateEmail - failed", zap.String("email", email), zap.Error(err))
		return false, err
	}

	log.Debug("CheckDuplicateEmail - completed", zap.String("email", email), zap.Bool("exists", exists))
	return exists, nil
}
