package handler

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type Duration struct {
	Since time.Time `json:"since"`
	Until time.Time `json:"until"`
}

type OptionalDuration struct {
	Since optional.Time `json:"since"`
	Until optional.Time `json:"until"`
}

type OptionalProjectDuration struct {
	Since OptionalYearWithSemester
	Until OptionalYearWithSemester
}

type OptionalYearWithSemester struct {
	Year     optional.Int64
	Semester optional.Int64
}

func convertToProjectDuration(since, until time.Time) domain.ProjectDuration {
	return domain.ProjectDuration{
		Since: timeToSem(since),
		Until: timeToSem(until),
	}
}

func semToTime(date domain.YearWithSemester) time.Time {
	year := int(date.Year)
	month := semesterToMonth[date.Semester]
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func timeToSem(t time.Time) domain.YearWithSemester {
	year := uint(t.Year())
	var semester uint
	for i, v := range semesterToMonth {
		if v == t.Month() {
			semester = uint(i)
		}
	}
	return domain.YearWithSemester{
		Year:     year,
		Semester: semester,
	}
}

func optionalSemToTime(date OptionalYearWithSemester) optional.Time {
	t := optional.Time{}
	if date.Year.Valid && date.Semester.Valid {
		year := int(date.Year.Int64)
		month := semesterToMonth[date.Semester.Int64]
		t.Time, t.Valid = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC), true
	} else {
		t.Valid = false
	}
	return t
}
