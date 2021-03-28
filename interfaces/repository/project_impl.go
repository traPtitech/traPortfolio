package repository

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

//TODO いつ？
var (
	semesterToMonth [2]time.Month = [2]time.Month{time.August, time.December}
)

type ProjectRepository struct {
	database.SQLHandler
}

func NewProjectRepository(sql database.SQLHandler) *ProjectRepository {
	return &ProjectRepository{SQLHandler: sql}
}

func (repo *ProjectRepository) PostProject(p *domain.ProjectDetail) (*model.Project, error) {
	project := model.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Start:       formatDuration(p.Duration.Since),
		End:         formatDuration(p.Duration.Until),
	}
	err := repo.Create(&project).Error()
	return &project, err
}

func formatDuration(date domain.YearWithSemester) time.Time {
	year := int(date.Year)
	month := semesterToMonth[date.Semester]
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}
