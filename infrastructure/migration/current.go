package migration

import (
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"gopkg.in/gormigrate.v1"
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
	}
}
