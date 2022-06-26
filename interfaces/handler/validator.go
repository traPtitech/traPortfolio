package handler

import (
	"errors"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
)

type validator struct {
	logger echo.Logger
}

func newValidator(logger echo.Logger) echo.Validator {
	return &validator{logger}
}

func (v *validator) Validate(i interface{}) error {
	if vld, ok := i.(vd.Validatable); ok {
		if err := vld.Validate(); err != nil {
			if e, ok := err.(vd.InternalError); ok {
				return e.InternalError()
			}

			return err
		}
	} else {
		v.logger.Errorf("%T is not validatable", i)
	}

	return nil
}

var (
	vdRuleNameLength        = vd.Length(1, 32)
	vdRuleDescriptionLength = vd.Length(1, 256)
	vdRuleResultLength      = vd.Length(0, 32)
	vdRuleAccountTypeMin    = vd.Min(0) // TODO: handler.AccountTypeをuint型にしたら消す
	vdRuleAccountTypeMax    = vd.Max(int(domain.AccountLimit) - 1)
	vdRuleEventLevelMin     = vd.Min(0) // TODO: handler.EventLevelをuint型にしたら消す
	vdRuleEventLevelMax     = vd.Max(int(domain.EventLevelLimit) - 1)
)

// path parameter structs

func (p GetUsersParams) Validate() error {
	if p.IncludeSuspended != nil && p.Name != nil {
		return errors.New("include_suspended and name cannot be specified at the same time")
	}

	return vd.ValidateStruct(&p,
		vd.Field(&p.IncludeSuspended),
		vd.Field(&p.Name, vd.NilOrNotEmpty),
	)
}

// request body structs

func (r AddAccountRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.DisplayName, vd.Required),
		vd.Field(&r.PrPermitted, vd.NotNil),
		vd.Field(&r.Type, vd.NotNil, vdRuleAccountTypeMin, vdRuleAccountTypeMax),
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

func (r AddProjectMembersRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Members, vd.Required, vd.Each(vd.Required, is.UUID)),
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
		vd.Field(&r.DisplayName),
		vd.Field(&r.PrPermitted),
		vd.Field(&r.Type, vdRuleAccountTypeMin, vdRuleAccountTypeMax),
		vd.Field(&r.Url, is.URL),
	)
}

func (r EditContestRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vdRuleDescriptionLength),
		vd.Field(&r.Duration),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vdRuleNameLength),
	)
}

func (r EditContestTeamRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vdRuleDescriptionLength),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vdRuleNameLength),
		vd.Field(&r.Result, vdRuleResultLength),
	)
}

func (r EditEventRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.EventLevel, vdRuleEventLevelMin, vdRuleEventLevelMax),
	)
}

func (r EditProjectRequest) Validate() error {
	return vd.ValidateStruct(&r,
		vd.Field(&r.Description, vdRuleDescriptionLength),
		vd.Field(&r.Duration),
		vd.Field(&r.Link, is.URL),
		vd.Field(&r.Name, vdRuleNameLength),
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
		vd.Field(&r.UserId, vd.Required, is.UUID),
	)
}
