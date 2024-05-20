package schema

import (
	"errors"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/traPtitech/traPortfolio/domain"
)

var (
	vdRuleNameLength        = vd.RuneLength(1, 32)
	vdRuleDisplayNameLength = vd.RuneLength(1, 256) // 外部アカウントのアカウント名文字数上限
	vdRuleDescriptionLength = vd.RuneLength(1, 256)
	vdRuleResultLength      = vd.RuneLength(0, 32)
	vdRuleAccountTypeMax    = vd.Max(domain.AccountLimit - 1)
	vdRuleEventLevelMax     = vd.Max(uint8(domain.EventLevelLimit) - 1)
)

// path parameter structs

func (p GetUsersParams) Validate() error {
	if p.IncludeSuspended != nil && p.Name != nil {
		return errors.New("include_suspended and name cannot be specified at the same time")
	}

	return vd.ValidateStruct(&p,
		vd.Field(&p.IncludeSuspended),
		vd.Field(&p.Name, vd.NilOrNotEmpty),
		vd.Field(&p.Limit, vd.Min(1), vd.NilOrNotEmpty),
	)
}

// request body structs

func (r AddAccountRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.DisplayName, vd.Required, vdRuleDisplayNameLength),
		vd.Field(&r.PrPermitted),
		vd.Field(&r.Type, vdRuleAccountTypeMax),
		vd.Field(&r.Url, vd.Required, is.URL),
	)
}

func (r AddContestTeamRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.Required, vdRuleDescriptionLength),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.Required, vdRuleNameLength),
		vd.Field(&r.Result, vdRuleResultLength),
	)
}

func (r EditProjectMembersRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Members, vd.Required),
	)
}

func (r CreateContestRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.Required, vdRuleDescriptionLength),
		vd.Field(&r.Duration, vd.Required),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.Required, vdRuleNameLength),
	)
}

func (r CreateProjectRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.Required, vdRuleDescriptionLength),
		vd.Field(&r.Duration, vd.Required),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.Required, vdRuleNameLength),
	)
}

func (r EditUserAccountRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.DisplayName, vd.NilOrNotEmpty, vdRuleDisplayNameLength),
		vd.Field(&r.PrPermitted),
		vd.Field(&r.Type, vdRuleAccountTypeMax),
		vd.Field(&r.Url, vd.NilOrNotEmpty, is.URL),
	)
}

func (r EditContestRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.NilOrNotEmpty, vdRuleDescriptionLength),
		vd.Field(&r.Duration, vd.NilOrNotEmpty),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.NilOrNotEmpty, vdRuleNameLength),
	)
}

func (r EditContestTeamRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.NilOrNotEmpty, vdRuleDescriptionLength),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.NilOrNotEmpty, vdRuleNameLength),
		vd.Field(&r.Result, vdRuleResultLength),
	)
}

func (r EditEventRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Level, vdRuleEventLevelMax),
	)
}

func (r EditProjectRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vd.NilOrNotEmpty, vdRuleDescriptionLength),
		vd.Field(&r.Duration, vd.NilOrNotEmpty),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vd.NilOrNotEmpty, vdRuleNameLength),
	)
}

func (r EditUserRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Bio, vdRuleDescriptionLength),
		vd.Field(&r.Check),
	)
}

// embedded structs

func (r Duration) Validate() error {
	if err := vd.ValidateStruct(&r,
		vd.Field(&r.Since, vd.Required),
		vd.Field(&r.Until),
	); err != nil {
		return err
	}

	if r.Until != nil {
		if r.Since.After(*r.Until) {
			return vd.ErrDateInvalid
		}
	}

	return nil
}

func (r MemberIDWithYearWithSemesterDuration) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Duration),
		vd.Field(&r.UserId, vd.Required, is.UUIDv4),
	)
}

// TODO: MemberIDsを埋め込んだリクエストボディを実装する
func (r MemberIDs) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Members, vd.Required, vd.Each(vd.Required, is.UUIDv4)),
	)
}
