package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
)

// Migrations is all db migrations
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v1(),
		v2(), // プロジェクト名とコンテスト名の重複禁止と文字数制限増加(32->128)
		v3(), // ユーザーにDisplayNameを追加
	}
}

func AllTables() []interface{} {
	return []interface{}{
		model.User{},
		model.Account{},
		model.Project{},
		model.ProjectMember{},
		model.EventLevelRelation{},
		model.Contest{},
		model.ContestTeam{},
		model.ContestTeamUserBelonging{},
		model.Group{},
		model.GroupUserBelonging{},
		model.GroupUserAdmin{},
	}
}
