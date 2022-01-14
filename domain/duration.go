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
