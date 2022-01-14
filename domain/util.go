package domain

type YearWithSemesterDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     int
	Semester int
}
