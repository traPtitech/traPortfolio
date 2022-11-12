//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config ../../.oapi.types.yml ../../docs/swagger/traPortfolio.v1.yaml

package handler

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type OptionalDuration struct {
	Since optional.Of[time.Time] `json:"since"`
	Until optional.Of[time.Time] `json:"until"`
}

type OptionalYearWithSemesterDuration struct {
	Since OptionalYearWithSemester
	Until OptionalYearWithSemester
}

type OptionalYearWithSemester struct {
	Year     optional.Of[int64]
	Semester optional.Of[int64]
}

func ConvertDuration(d domain.YearWithSemesterDuration) YearWithSemesterDuration {
	return newYearWithSemesterDuration(d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
}

func newYearWithSemesterDuration(sinceYear, sinceSemester, untilYear, untilSemester int) YearWithSemesterDuration {
	return YearWithSemesterDuration{
		Since: YearWithSemester{
			Year:     sinceYear,
			Semester: Semester(sinceSemester),
		},
		Until: &YearWithSemester{
			Year:     untilYear,
			Semester: Semester(untilSemester),
		},
	}
}
