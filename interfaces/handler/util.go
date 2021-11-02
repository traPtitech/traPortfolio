package handler

import (
	"time"

	"github.com/traPtitech/traPortfolio/util/optional"
)

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

func convertToProjectDuration(since, until time.Time) ProjectDuration {
	s := timeToSem(since)
	u := timeToSem(until)
	return ProjectDuration{
		Since: s,
		Until: &u,
	}
}

func semToTime(date YearWithSemester) time.Time {
	year := int(date.Year)
	month := semesterToMonth[date.Semester]
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func timeToSem(t time.Time) YearWithSemester {
	year := t.Year()
	var semester Semester
	for i, v := range semesterToMonth {
		if v == t.Month() {
			semester = Semester(i)
		}
	}
	return YearWithSemester{
		Year:     year,
		Semester: semester,
	}
}

func optionalSemToTime(date YearWithSemester) optional.Time {
	t := optional.Time{}
	year := date.Year
	month := semesterToMonth[date.Semester]
	t.Time, t.Valid = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC), true

	return t
}
