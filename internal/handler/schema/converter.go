//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config .oapi.types.yml ../../../docs/swagger/traPortfolio.v1.yaml

package schema

import (
	"github.com/traPtitech/traPortfolio/internal/domain"
)

func ConvertDuration(d domain.YearWithSemesterDuration) YearWithSemesterDuration {
	return YearWithSemesterDuration{
		Since: YearWithSemester{
			Year:     d.Since.Year,
			Semester: Semester(d.Since.Semester),
		},
		Until: &YearWithSemester{
			Year:     d.Until.Year,
			Semester: Semester(d.Until.Semester),
		},
	}
}
