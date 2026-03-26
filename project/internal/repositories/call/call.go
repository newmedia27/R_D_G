package call

import "project/internal/clients/livekit"

type Repository struct {
	repo *livekit.Client
}

func NewRepository(repo *livekit.Client) *Repository {
	return &Repository{repo: repo}
}

func (r *Repository) GetToken(userID, roomName string) (string, error) {
	token, err := r.repo.GenerateToken(userID, roomName)
	if err != nil {
		return "", err
	}
	return token, nil
}
