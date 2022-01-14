package domain

type YearWithSemesterDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     uint
	Semester uint
}
