package random

import (
	"math/rand/v2"
	"time"

	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
)

// CreateContestArgs
func CreateContestArgs() *repository.CreateContestArgs {
	return &repository.CreateContestArgs{
		Name:        AlphaNumeric(),
		Description: AlphaNumeric(),
		Links:       Array(RandURLString, 1, 3),
		Since:       time.Now(),
		Until:       Optional(time.Now().Add(time.Hour)),
	}
}

// UpdateContestArgs 全てのフィールドがvalidなUpdateContestArgsを生成します
func UpdateContestArgs() *repository.UpdateContestArgs {
	a := repository.UpdateContestArgs{
		Name:        optional.From(AlphaNumeric()),
		Description: optional.From(AlphaNumeric()),
		Links:       optional.From(Array(RandURLString, 1, 3)),
		Since:       optional.From(Time()),
		Until:       optional.From(Time()),
	}
	return &a
}

// CreateContestTeamArgs
func CreateContestTeamArgs() *repository.CreateContestTeamArgs {
	return &repository.CreateContestTeamArgs{
		Name:        AlphaNumeric(),
		Result:      Optional(AlphaNumeric()),
		Links:       Array(RandURLString, 1, 3),
		Description: AlphaNumeric(),
	}
}

// UpdateContestTeamArgs 全てのフィールドがvalidなUpdateContestTeamArgsを生成します
func UpdateContestTeamArgs() *repository.UpdateContestTeamArgs {
	a := repository.UpdateContestTeamArgs{
		Name:        optional.From(AlphaNumeric()),
		Result:      optional.From(AlphaNumeric()),
		Links:       optional.From(Array(RandURLString, 1, 3)),
		Description: optional.From(AlphaNumeric()),
	}
	return &a
}

// OptUpdateContestTeamArgs validかどうかも含めてランダムなUpdateContestTeamArgsを生成します
func OptUpdateContestTeamArgs() *repository.UpdateContestTeamArgs {
	a := repository.UpdateContestTeamArgs{
		Name:        Optional(AlphaNumeric()),
		Result:      Optional(AlphaNumeric()),
		Links:       Optional(Array(RandURLString, 1, 3)),
		Description: Optional(AlphaNumeric()),
	}
	return &a
}

// CreateProjectArgs 全てのフィールドがvalidなCreateProjectArgsを生成します
func CreateProjectArgs() *repository.CreateProjectArgs {
	return &repository.CreateProjectArgs{
		Name:          AlphaNumeric(),
		Description:   AlphaNumeric(),
		Links:         Array(RandURLString, 1, 3),
		SinceYear:     2100,
		SinceSemester: 0,
		UntilYear:     2100,
		UntilSemester: 1,
	}
}

// UpdateProjectArgs 全てのフィールドがvalidなUpdateProjectArgsを生成します
func UpdateProjectArgs() *repository.UpdateProjectArgs {
	a := repository.UpdateProjectArgs{
		Name:          optional.From(AlphaNumeric()),
		Description:   optional.From(AlphaNumeric()),
		Links:         optional.From(Array(RandURLString, 1, 3)),
		SinceYear:     optional.From(int64(2100)), // TODO: intでよさそう
		SinceSemester: optional.From(int64(0)),
		UntilYear:     optional.From(int64(2100)),
		UntilSemester: optional.From(int64(1)),
	}
	return &a
}

// OptUpdateProjectArgs validかどうかも含めてランダムなUpdateProjectArgsを生成します
func OptUpdateProjectArgs() *repository.UpdateProjectArgs {
	a := repository.UpdateProjectArgs{
		Name:          Optional(AlphaNumeric()),
		Description:   Optional(AlphaNumeric()),
		Links:         Optional(Array(RandURLString, 1, 3)),
		SinceYear:     Optional(int64(2100)), // TODO: intでよさそう
		SinceSemester: Optional(int64(0)),
		UntilYear:     Optional(int64(2100)),
		UntilSemester: Optional(int64(1)),
	}
	return &a
}

// UpdateUserArgs 全てのフィールドがvalidなUpdateUserArgsを生成します
func UpdateUserArgs() *repository.UpdateUserArgs {
	a := repository.UpdateUserArgs{
		Description: optional.From(AlphaNumeric()),
		Check:       optional.From(Bool()),
	}
	return &a
}

// OptUpdateUserArgs validかどうかも含めてランダムなUpdateUserArgsを生成します
func OptUpdateUserArgs() *repository.UpdateUserArgs {
	a := repository.UpdateUserArgs{
		Description: Optional(AlphaNumeric()),
		Check:       Optional(Bool()),
	}
	return &a
}

// UpdateAccountArgs 全てのフィールドがvalidなUpdateAccountArgsを生成します
func UpdateAccountArgs() *repository.UpdateAccountArgs {
	a := repository.UpdateAccountArgs{
		DisplayName: optional.From(AlphaNumeric()),
		Type:        optional.From(rand.N(domain.AccountLimit)),
		URL:         optional.From(RandURLString()),
		PrPermitted: optional.From(Bool()),
	}
	return &a
}

// OptUpdateAccountArgs validかどうかも含めてランダムなUpdateAccountArgsを生成します
func OptUpdateAccountArgs() *repository.UpdateAccountArgs {
	a := repository.UpdateAccountArgs{
		DisplayName: Optional(AlphaNumeric()),
		Type:        Optional(rand.N(domain.AccountLimit)),
		URL:         Optional(RandURLString()),
		PrPermitted: Optional(Bool()),
	}
	return &a
}
