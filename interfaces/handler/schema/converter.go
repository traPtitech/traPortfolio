//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config .oapi.types.yml ../../../docs/swagger/traPortfolio.v1.yaml

package schema

import (
	"github.com/traPtitech/traPortfolio/domain"
)

func ConvertDuration(d domain.YearWithSemesterDuration) YearWithSemesterDuration {
	since := YearWithSemester{
		Year:     d.Since.Year,
		Semester: Semester(d.Since.Semester),
	}
	u, ok := d.Until.V()
	if !ok {
		return YearWithSemesterDuration{
			Since: since,
			Until: nil,
		}
	}

	until := YearWithSemester{
		Year:     u.Year,
		Semester: Semester(u.Semester),
	}

	return YearWithSemesterDuration{
		Since: since,
		Until: &until,
	}
}
