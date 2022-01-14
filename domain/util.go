package domain

type YearWithSemesterDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     int
	Semester int
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
