package domain

import "github.com/traPtitech/traPortfolio/util/optional"

type YearWithSemester struct {
	Year     int
	Semester int
}

func (ys YearWithSemester) IsValid() bool {
	return ys.Year >= 1970 && ys.Semester >= 0 && ys.Semester < 2
}

func (ys YearWithSemester) After(ys2 YearWithSemester) bool {
	return ys.Year > ys2.Year || (ys.Year == ys2.Year && ys.Semester > ys2.Semester)
}

type YearWithSemesterDuration struct {
	Since YearWithSemester
	Until optional.Of[YearWithSemester]
}

// TODO: !since.IsValid() || !until.IsValid()のときエラーを返す
func NewYearWithSemesterDuration(sinceYear, sinceSemester, untilYear, untilSemester int) YearWithSemesterDuration {
	since := YearWithSemester{
		Year:     sinceYear,
		Semester: sinceSemester,
	}
	until := YearWithSemester{
		Year:     untilYear,
		Semester: untilSemester,
	}

	return YearWithSemesterDuration{
		Since: since,
		Until: optional.New(until, until.IsValid()),
	}
}

func (d YearWithSemesterDuration) IsValid() bool {
	s := d.Since
	if !s.IsValid() {
		return false
	}

	u, ok := d.Until.V()
	if !ok {
		return true
	}

	return u.IsValid() && !s.After(u)
}

// out.Since <= in.Since <= in.Until <= out.Until
func (out YearWithSemesterDuration) Includes(in YearWithSemesterDuration) bool {
	if !in.IsValid() || !out.IsValid() || out.Since.After(in.Since) {
		return false
	}

	outUntil, outUntilOK := out.Until.V()
	inUntil, inUntilOK := in.Until.V()
	if !inUntilOK {
		if outUntilOK {
			return false
		}

		return !out.Since.After(in.Since)
	}

	if !outUntilOK {
		return true
	}

	return !inUntil.After(outUntil)
}
