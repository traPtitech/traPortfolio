package migration

import (
	"github.com/traPtitech/traPortfolio/domain"
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
		domain.User{},
		domain.EventLevelRelation{},
	}
}
