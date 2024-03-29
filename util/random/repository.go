package random

import (
	"math/rand/v2"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
)

// UpdateContestArgs 全てのフィールドがvalidなUpdateContestArgsを生成します
func UpdateContestArgs() *repository.UpdateContestArgs {
	a := repository.UpdateContestArgs{
		Name:        optional.From(AlphaNumeric()),
		Description: optional.From(AlphaNumeric()),
		Link:        optional.From(RandURLString()),
		Since:       optional.From(Time()),
		Until:       optional.From(Time()),
	}
	return &a
}

// OptUpdateContestArgs validかどうかも含めてランダムなUpdateContestArgsを生成します
func OptUpdateContestArgs() *repository.UpdateContestArgs {
	a := repository.UpdateContestArgs{
		Name:        Optional(AlphaNumeric()),
		Description: Optional(AlphaNumeric()),
		Link:        Optional(RandURLString()),
		Since:       Optional(Time()),
		Until:       Optional(Time()),
	}
	return &a
}

// UpdateContestTeamArgs 全てのフィールドがvalidなUpdateContestTeamArgsを生成します
func UpdateContestTeamArgs() *repository.UpdateContestTeamArgs {
	a := repository.UpdateContestTeamArgs{
		Name:        optional.From(AlphaNumeric()),
		Result:      optional.From(AlphaNumeric()),
		Link:        optional.From(RandURLString()),
		Description: optional.From(AlphaNumeric()),
	}
	return &a
}

// OptUpdateContestTeamArgs validかどうかも含めてランダムなUpdateContestTeamArgsを生成します
func OptUpdateContestTeamArgs() *repository.UpdateContestTeamArgs {
	a := repository.UpdateContestTeamArgs{
		Name:        Optional(AlphaNumeric()),
		Result:      Optional(AlphaNumeric()),
		Link:        Optional(RandURLString()),
		Description: Optional(AlphaNumeric()),
	}
	return &a
}

// UpdateProjectArgs 全てのフィールドがvalidなUpdateProjectArgsを生成します
func UpdateProjectArgs() *repository.UpdateProjectArgs {
	a := repository.UpdateProjectArgs{
		Name:          optional.From(AlphaNumeric()),
		Description:   optional.From(AlphaNumeric()),
		Link:          optional.From(RandURLString()),
		SinceYear:     optional.From(int64(2100)), // TODO: intでよさそう
		SinceSemester: optional.From(int64(2)),
		UntilYear:     optional.From(int64(2100)),
		UntilSemester: optional.From(int64(2)),
	}
	return &a
}

// OptUpdateProjectArgs validかどうかも含めてランダムなUpdateProjectArgsを生成します
func OptUpdateProjectArgs() *repository.UpdateProjectArgs {
	a := repository.UpdateProjectArgs{
		Name:          Optional(AlphaNumeric()),
		Description:   Optional(AlphaNumeric()),
		Link:          Optional(AlphaNumeric()),
		SinceYear:     Optional(int64(2100)), // TODO: intでよさそう
		SinceSemester: Optional(int64(2)),
		UntilYear:     Optional(int64(2100)),
		UntilSemester: Optional(int64(2)),
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
