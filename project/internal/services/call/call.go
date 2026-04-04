package call

import (
	"project/internal/repositories/call"
)

type callRepository interface {
	GetToken(userId, roomName string) (string, error)
}
type Service struct {
	repo callRepository
}

func NewService(repo *call.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetToken(userId, roomName string) (string, error) {
	return s.repo.GetToken(userId, roomName)
}
