// Package migration migrate current struct
package migration

import (
	"fmt"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

// v1 unique_index:idx_room_uniqueの削除
func v2() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2",
		Migrate: func(db *gorm.DB) error {
			if err := db.AutoMigrate(&v2Project{}, &v2Contest{}, &v2ContestTeam{}); err != nil {
				return err
			}

			// プロジェクト名の重複禁止
			{
				projects := make([]*model.Project, 0)
				if err := db.Find(&projects).Error; err != nil {
					return err
				}

				projectMap := make(map[string][]uuid.UUID, len(projects))
				for _, p := range projects {
					projectMap[p.Name] = append(projectMap[p.Name], p.ID)
				}

				updates := make(map[uuid.UUID]string, len(projects))
				noDuplicate := false
				for !noDuplicate {
					noDuplicate = true
					for name, arr := range projectMap {
						if len(arr) <= 1 {
							continue
						}
						noDuplicate = false
						for i, pid := range arr {
							if i == 0 {
								projectMap[name] = []uuid.UUID{pid}
								continue
							}
							nameNew := fmt.Sprintf("%s (%d)", name, i)
							updates[pid] = nameNew
							projectMap[nameNew] = append(projectMap[nameNew], pid)
						}
					}
				}

				for id, nameNew := range updates {
					err := db.
						Model(&model.Project{}).
						Where(&model.Project{ID: id}).
						Update("name", nameNew).
						Error
					if err != nil {
						return err
					}
				}
			}

			// コンテスト名の重複禁止
			{
				contests := make([]*model.Contest, 0)
				if err := db.Find(&contests).Error; err != nil {
					return err
				}

				contestMap := make(map[string][]uuid.UUID, len(contests))
				for _, c := range contests {
					contestMap[c.Name] = append(contestMap[c.Name], c.ID)
				}

				updates := make(map[uuid.UUID]string, len(contests))
				noDuplicate := false
				for !noDuplicate {
					noDuplicate = true
					for name, arr := range contestMap {
						if len(arr) <= 1 {
							continue
						}
						noDuplicate = false
						for i, cid := range arr {
							if i == 0 {
								contestMap[name] = []uuid.UUID{cid}
								continue
							}
							nameNew := fmt.Sprintf("%s (%d)", name, i)
							updates[cid] = nameNew
							contestMap[nameNew] = append(contestMap[nameNew], cid)
						}
					}
				}

				for id, nameNew := range updates {
					err := db.
						Model(&model.Contest{}).
						Where(&model.Contest{ID: id}).
						Update("name", nameNew).
						Error
					if err != nil {
						return err
					}
				}
			}

			return db.
				Table("portfolio").
				Error
		},
	}
}

type v2Project struct {
	ID            uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name          string    `gorm:"type:varchar(128)"` // 制限増加 (32->128)
	Description   string    `gorm:"type:text"`
	Link          string    `gorm:"type:text"`
	SinceYear     int       `gorm:"type:smallint(4);not null"`
	SinceSemester int       `gorm:"type:tinyint(1);not null"`
	UntilYear     int       `gorm:"type:smallint(4);not null"`
	UntilSemester int       `gorm:"type:tinyint(1);not null"`
	CreatedAt     time.Time `gorm:"precision:6"`
	UpdatedAt     time.Time `gorm:"precision:6"`
}

func (*v2Project) TableName() string {
	return "projects"
}

type v2Contest struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name        string    `gorm:"type:varchar(128)"` // 制限増加 (32->128)
	Description string    `gorm:"type:text"`
	Link        string    `gorm:"type:text"`
	Since       time.Time `gorm:"precision:6"`
	Until       time.Time `gorm:"precision:6"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

func (*v2Contest) TableName() string {
	return "contests"
}

type v2ContestTeam struct {
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

func (*v2ContestTeam) TableName() string {
	return "contest_teams"
}
