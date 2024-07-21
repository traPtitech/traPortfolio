package domain

import (
	"testing"

	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
)

func Test_YearWithSemester_IsValid(t *testing.T) {
	tests := map[string]struct {
		ys   YearWithSemester
		want bool
	}{
		"valid": {
			ys:   YearWithSemester{Year: 2021, Semester: 0},
			want: true,
		},
		"year must >=1970": {
			ys:   YearWithSemester{Year: 1969, Semester: 0},
			want: false,
		},
		"semester must >=0": {
			ys:   YearWithSemester{Year: 2021, Semester: -1},
			want: false,
		},
		"semester must <2": {
			ys:   YearWithSemester{Year: 2021, Semester: 2},
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := test.ys.IsValid(); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func Test_YearWithSemester_After(t *testing.T) {
	tests := map[string]struct {
		ys1  YearWithSemester
		ys2  YearWithSemester
		want bool
	}{
		"ys1.Year > ys2.Year": {
			ys1:  YearWithSemester{Year: 2021, Semester: 0},
			ys2:  YearWithSemester{Year: 2020, Semester: 0},
			want: true,
		},
		"ys1.Year < ys2.Year": {
			ys1:  YearWithSemester{Year: 2020, Semester: 0},
			ys2:  YearWithSemester{Year: 2021, Semester: 0},
			want: false,
		},
		"ys1.Year == ys2.Year, ys1.Semester > ys2.Semester": {
			ys1:  YearWithSemester{Year: 2021, Semester: 1},
			ys2:  YearWithSemester{Year: 2021, Semester: 0},
			want: true,
		},
		"ys1.Year == ys2.Year, ys1.Semester < ys2.Semester": {
			ys1:  YearWithSemester{Year: 2021, Semester: 0},
			ys2:  YearWithSemester{Year: 2021, Semester: 1},
			want: false,
		},
		"ys1.Year == ys2.Year, ys1.Semester == ys2.Semester": {
			ys1:  YearWithSemester{Year: 2021, Semester: 0},
			ys2:  YearWithSemester{Year: 2021, Semester: 0},
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := test.ys1.After(test.ys2); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func Test_YearWithSemesterDuration_IsValid(t *testing.T) {
	tests := map[string]struct {
		d    YearWithSemesterDuration
		want bool
	}{
		"since is invalid": {
			d: YearWithSemesterDuration{
				Since: YearWithSemester{Year: 1969, Semester: 0},
				Until: optional.New(
					YearWithSemester{Year: 2021, Semester: 0},
					true,
				),
			},
			want: false,
		},
		"until is invalid": {
			d: YearWithSemesterDuration{
				Since: YearWithSemester{Year: 2021, Semester: 0},
				Until: optional.New(
					YearWithSemester{Year: 0, Semester: 0},
					false,
				),
			},
			want: true,
		},
		"since > until": {
			d: YearWithSemesterDuration{
				Since: YearWithSemester{Year: 2021, Semester: 1},
				Until: optional.New(
					YearWithSemester{Year: 2021, Semester: 0},
					true,
				),
			},
			want: false,
		},
		"valid": {
			d: YearWithSemesterDuration{
				Since: YearWithSemester{Year: 2021, Semester: 0},
				Until: optional.New(
					YearWithSemester{Year: 2021, Semester: 1},
					true,
				),
			},
			want: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := test.d.IsValid(); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func Test_YearWithSemesterDuration_Includes(t *testing.T) {
	tests := map[string]struct {
		out  YearWithSemesterDuration
		in   YearWithSemesterDuration
		want bool
	}{
		"ok: out.Since <= in.Since <= in.Until <= out.Until": {
			out:  NewYearWithSemesterDuration(1970, 0, 2000, 1),
			in:   NewYearWithSemesterDuration(1980, 0, 1990, 1),
			want: true,
		},
		"ng: out.Since > in.Since": {
			out:  NewYearWithSemesterDuration(1980, 0, 2000, 1),
			in:   NewYearWithSemesterDuration(1970, 0, 1990, 1),
			want: false,
		},
		"ng: out.Until < in.Until": {
			out:  NewYearWithSemesterDuration(1970, 0, 2000, 1),
			in:   NewYearWithSemesterDuration(1980, 0, 2010, 1),
			want: false,
		},
		"ok: out.Until is nil": {
			out:  NewYearWithSemesterDuration(1970, 0, 0, 0),
			in:   NewYearWithSemesterDuration(1980, 0, 2000, 1),
			want: true,
		},
		"ng: in.Until is nil": {
			out:  NewYearWithSemesterDuration(1970, 0, 2000, 1),
			in:   NewYearWithSemesterDuration(1980, 0, 0, 0),
			want: false,
		},
		"ok: out.Until and in.Until are both nil": {
			out:  NewYearWithSemesterDuration(1970, 0, 0, 0),
			in:   NewYearWithSemesterDuration(1980, 0, 0, 0),
			want: true,
		},
		"ng: out is invalid": {
			out:  NewYearWithSemesterDuration(2000, 0, 1999, 0),
			in:   NewYearWithSemesterDuration(1980, 0, 1990, 1),
			want: false,
		},
		"ng: in is invalid": {
			out:  NewYearWithSemesterDuration(1970, 0, 2000, 1),
			in:   NewYearWithSemesterDuration(1980, 0, 1979, 0),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := test.out.Includes(test.in); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
