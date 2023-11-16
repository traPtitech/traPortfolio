package domain

type YearWithSemesterDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     int
	Semester int
}

func (ys YearWithSemester) After(ys2 YearWithSemester) bool {
	return ys.Year > ys2.Year || (ys.Year == ys2.Year && ys.Semester > ys2.Semester)
}

func NewYearWithSemesterDuration(sinceYear, sinceSemester, untilYear, untilSemester int) YearWithSemesterDuration {
	return YearWithSemesterDuration{
		Since: YearWithSemester{
			Year:     sinceYear,
			Semester: sinceSemester,
		},
		Until: YearWithSemester{
			Year:     untilYear,
			Semester: untilSemester,
		},
	}
}

func (d YearWithSemesterDuration) IsValid() bool {
	s := d.Since
	u := d.Until

	return u.After(s) && s.Year >= 1970 && u.Year < 2070 && s.Semester >= 0 && s.Semester < 2 && u.Semester >= 0 && u.Semester < 2
}

func (out YearWithSemesterDuration) Includes(in YearWithSemesterDuration) bool {
	return !out.Since.After(in.Since) && !in.Until.After(out.Until)
}
