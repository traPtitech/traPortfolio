//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config ../../.oapi.types.yml ../../docs/swagger/traPortfolio.v1.yaml

package handler

import (
	"time"

	"github.com/traPtitech/traPortfolio/util/optional"
)

type OptionalDuration struct {
	Since optional.Time `json:"since"`
	Until optional.Time `json:"until"`
}

type OptionalYearWithSemesterDuration struct {
	Since OptionalYearWithSemester
	Until OptionalYearWithSemester
}

type OptionalYearWithSemester struct {
	Year     optional.Int64
	Semester optional.Int64
}

func convertToYearWithSemesterDuration(since, until time.Time) YearWithSemesterDuration {
	s := timeToSem(since)
	u := timeToSem(until)
	return YearWithSemesterDuration{
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
