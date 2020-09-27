package user

import (
	"context"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// TODO: 設定読み込み
const traQToken = ""
const portalToken = ""

// traQUser traQ上のユーザー情報
type traQUser struct {
	State       uint8  `json:"state"` // TODO: 特別な型にする
	Bot         bool   `json:"bot"`
	DisplayName string `json:"displayName"`
	Name        string `json:"name"` // これはportalのIDと一致
}

// portalUser Portal上のユーザー情報
type portalUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

func fetchTraQUser(ctx context.Context, name string) (traQUser, error) {
	// TODO
	return traQUser{}, nil
}

func fetchPortalUser(ctx context.Context, name string) (portalUser, error) {
	// TODO
	return portalUser{}, nil
}

// User Portfolioのレスポンスで使うユーザー情報
type User struct {}

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) getUser(ctx context.Context, name string) User {
	_, _ = fetchTraQUser(ctx, name)
	_, _ = fetchPortalUser(ctx, name)
	_, _ = s.repo.Get(name)
	return User{}
}
