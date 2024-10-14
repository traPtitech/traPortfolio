// Package migration migrate current struct
package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

// v3 contestTeam, group, projectの複数リンク対応
func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(db *gorm.DB) error {
			if err := db.AutoMigrate(&v3ContestLink{}, &v3ContestTeamLink{}, &v3GroupLink{}, &v3ProjectLink{}); err != nil {
				return err
			}

			// contestのlinkをcontest_linksに移動
			{
				contests := make([]*v3OldContest, 0)
				if err := db.Find(&contests).Error; err != nil {
					return err
				}

				for _, contest := range contests {
					contestLink := v3ContestLink{
						ID:    contest.ID,
						Order: 0,
						Link:  contest.Link,
					}
					if err := db.Create(&contestLink).Error; err != nil {
						return err
					}
				}

				if err := db.Migrator().DropColumn(v3OldContest{}, "link"); err != nil {
					return err
				}
			}

			// contest_teamsのlinkをcontest_team_linksに移動
			{
				contestTeams := make([]*v3OldContestTeam, 0)
				if err := db.Find(&contestTeams).Error; err != nil {
					return err
				}

				for _, contestTeam := range contestTeams {
					teamLink := v3ContestTeamLink{
						ID:    contestTeam.ID,
						Order: 0,
						Link:  contestTeam.Link,
					}
					if err := db.Create(&teamLink).Error; err != nil {
						return err
					}
				}

				if err := db.Migrator().DropColumn(v3OldContestTeam{}, "link"); err != nil {
					return err
				}
			}

			// groupのlinkをgroup_linksに移動
			{
				groups := make([]*v3OldGroup, 0)
				if err := db.Find(&groups).Error; err != nil {
					return err
				}

				for _, group := range groups {
					groupLink := v3GroupLink{
						ID:    group.GroupID,
						Order: 0,
						Link:  group.Link,
					}
					if err := db.Create(&groupLink).Error; err != nil {
						return err
					}
				}

				if err := db.Migrator().DropColumn(v3OldGroup{}, "link"); err != nil {
					return err
				}
			}

			// projectのlinkをgroup_linksに移動
			{
				projects := make([]*v3OldProject, 0)
				if err := db.Find(&projects).Error; err != nil {
					return err
				}

				for _, project := range projects {
					projectLink := v3ProjectLink{
						ID:    project.ID,
						Order: 0,
						Link:  project.Link,
					}
					if err := db.Create(&projectLink).Error; err != nil {
						return err
					}
				}

				if err := db.Migrator().DropColumn(v3OldProject{}, "link"); err != nil {
					return err
				}
			}

			return db.
				Table("portfolio").
				Error
		},
	}
}

type v3OldContest struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*v3OldContest) TableName() string {
	return "contests"
}

type v3ContestLink struct {
	ID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order int       `gorm:"type:int;not null;primaryKey"`
	Link  string    `gorm:"type:text;not null"`
}

func (*v3ContestLink) TableName() string {
	return "contest_links"
}

type v3OldContestTeam struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	ContestID   uuid.UUID `gorm:"type:char(36);not null"`
	Name        string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:text"`
	Result      string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`

	Contest model.Contest `gorm:"foreignKey:ContestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*v3OldContestTeam) TableName() string {
	return "contest_teams"
}

type v3ContestTeamLink struct {
	ID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order int       `gorm:"type:int;not null;primaryKey"`
	Link  string    `gorm:"type:text;not null"`
}

func (*v3ContestTeamLink) TableName() string {
	return "contest_team_links"
}

type v3OldGroup struct {
	GroupID     uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(32)"`
	Link        string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*v3OldGroup) TableName() string {
	return "groups"
}

type v3GroupLink struct {
	ID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order int       `gorm:"type:int;not null;primaryKey"`
	Link  string    `gorm:"type:text;not null"`
}

func (*v3GroupLink) TableName() string {
	return "group_links"
}

type v3OldProject struct {
	ID            uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name          string    `gorm:"type:varchar(128)"`
	Description   string    `gorm:"type:text"`
	Link          string    `gorm:"type:text"`
	SinceYear     int       `gorm:"type:smallint(4);not null"`
	SinceSemester int       `gorm:"type:tinyint(1);not null"`
	UntilYear     int       `gorm:"type:smallint(4);not null"`
	UntilSemester int       `gorm:"type:tinyint(1);not null"`
	CreatedAt     time.Time `gorm:"precision:6"`
	UpdatedAt     time.Time `gorm:"precision:6"`
}

func (*v3OldProject) TableName() string {
	return "projects"
}

type v3ProjectLink struct {
	ID    uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Order int       `gorm:"type:int;not null;primaryKey"`
	Link  string    `gorm:"type:text;not null"`
}

func (*v3ProjectLink) TableName() string {
	return "project_links"
}
