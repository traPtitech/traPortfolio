package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

// Migrations is all db migrations
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v1(),
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
		model.GroupUser{},
		model.GroupUserBelonging{},
	}
}
