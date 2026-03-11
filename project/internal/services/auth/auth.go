package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"project/internal/logger"
	"project/internal/models"
)

type sessionRepository interface {
	Create(ctx context.Context, s *models.Session) error
	FindByToken(ctx context.Context, token string) (*models.Session, error)
	FindByUserID(ctx context.Context, userID string) (*models.Session, error)
	FindByDeviceId(ctx context.Context, deviceID string) (*models.Session, error)
	DeleteByDeviceId(ctx context.Context, deviceID string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	FindAllByUserID(ctx context.Context, userID string) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	FindByUserIDAndFingerprint(ctx context.Context, userID, fingerprint string) (*models.Session, error)
}
type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, u *models.User) (*models.User, error)
}
type Service struct {
	users    userRepository
	sessions sessionRepository

	jwtSecret        string
	refreshSecret    string
	jwtExpiresIn     time.Duration
	refreshExpiresIn time.Duration
}

func NewService(users userRepository, sessions sessionRepository, jwtSecret string, jwtExpiresIn time.Duration, refreshExpiresIn time.Duration) *Service {
	return &Service{
		users:            users,
		sessions:         sessions,
		jwtSecret:        jwtSecret,
		jwtExpiresIn:     jwtExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*models.User, error) {
	_, err := s.users.FindByEmail(ctx, input.Email)
	if err == nil {
		return nil, models.ErrUserDuplicate
	}
	if !errors.Is(err, models.ErrUserNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &models.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Username:     input.Username,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := s.users.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return createdUser, nil
}

type LoginInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
}
type LoginOutput struct {
	AccessToken     string       `json:"access_token"`
	RefreshToken    string       `json:"refresh_token"`
	DeviceSessionID string       `json:"device_session_id"`
	User            *models.User `json:"user"`
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	usr, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("can't find user: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(input.Password)); err != nil {
		return nil, models.ErrInvalidPassword
	}

	access_token, err := s.CreateAuthToken(usr.ID)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refresh_token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	deviceSessionId, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("create device session id: %w", err)
	}

	fingerprint := generateFingerprint(input.UserAgent, input.IP)

	old, _ := s.sessions.FindByUserIDAndFingerprint(ctx, usr.ID, fingerprint)
	if old != nil {
		if err = s.sessions.DeleteByDeviceId(ctx, old.DeviceSessionID); err != nil {
			return nil, fmt.Errorf("delete old session: %w", err)
		}
	}

	now := time.Now()
	session := &models.Session{
		UserID:           usr.ID,
		UserAgent:        input.UserAgent,
		IP:               input.IP,
		RefreshTokenHash: hashToken(refresh_token),
		DeviceSessionID:  deviceSessionId,
		FingerPrint:      fingerprint,
		ExpiresAt:        now.Add(s.refreshExpiresIn),
		CreatedAt:        now,
		LastUsedAt:       now,
	}

	err = s.sessions.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &LoginOutput{
		AccessToken:     access_token,
		RefreshToken:    refresh_token,
		DeviceSessionID: deviceSessionId,
		User:            usr,
	}, nil
}

type RefreshInput struct {
	RefreshToken    string `json:"refresh_token"`
	UserAgent       string `json:"user_agent"`
	IP              string `json:"ip"`
	DeviceSessionID string `json:"device_session_id"`
}
type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	refreshHash := hashToken(input.RefreshToken)
	session, err := s.sessions.FindByToken(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, models.ErrSessionNotFound) {
			return nil, models.ErrSessionNotFound
		}
		return nil, fmt.Errorf("find session: %w", err)
	}
	if session.DeviceSessionID != input.DeviceSessionID {
		return nil, models.ErrInvalidDevice
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, models.ErrSessionExpired
	}

	fingerprint := generateFingerprint(input.UserAgent, input.IP)
	if session.FingerPrint != fingerprint {
		logger.LogWarnContext(ctx, "fingerprint mismatch — session invalidated",
			slog.String("user_id", session.UserID),
			slog.String("device_session_id", session.DeviceSessionID),
			slog.String("old_ip", session.IP),
			slog.String("new_ip", input.IP),
			slog.String("old_ua", session.UserAgent),
			slog.String("new_ua", input.UserAgent),
		)
		err = s.sessions.DeleteByDeviceId(ctx, input.DeviceSessionID)
		if err != nil {
			return nil, fmt.Errorf("delete session with suspicious activity: %w", err)
		}
		return nil, models.ErrSuspiciousActivity
	}

	access_token, err := s.CreateAuthToken(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	newRefreshToken, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	session.RefreshTokenHash = hashToken(newRefreshToken)
	session.LastUsedAt = time.Now()
	session.UserAgent = input.UserAgent
	session.IP = input.IP
	session.FingerPrint = fingerprint

	if err = s.sessions.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	return &RefreshOutput{
		AccessToken:  access_token,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, deviceSessionId string) error {
	err := s.sessions.DeleteByDeviceId(ctx, deviceSessionId)
	if err != nil && !errors.Is(err, models.ErrSessionNotFound) {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *Service) LogoutAll(ctx context.Context, userID string) error {
	err := s.sessions.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}
	return nil
}
