// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package handler

import (
	"time"

	"github.com/gofrs/uuid"
)

// Defines values for Semester.
const (
	First  Semester = 0
	Second Semester = 1
)

// Defines values for UserAccountState.
const (
	Active      UserAccountState = 1
	Deactivated UserAccountState = 0
	Suspended   UserAccountState = 2
)

// Account アカウントへのリンク
type Account struct {
	// DisplayName 外部アカウントの表示名
	DisplayName string `json:"displayName"`

	// Id アカウントUUID
	Id uuid.UUID `json:"id"`

	// PrPermitted 広報での利用が許可されているかどうか
	PrPermitted PrPermitted `json:"prPermitted"`

	// Type アカウントの種類
	Type AccountType `json:"type"`

	// Url アカウントurl
	Url string `json:"url"`
}

// AccountType アカウントの種類
type AccountType = uint8

// AddAccountRequest 新規アカウントリクエスト
type AddAccountRequest struct {
	// DisplayName 外部アカウントの表示名
	DisplayName string `json:"displayName"`

	// PrPermitted 広報での利用が許可されているかどうか
	PrPermitted PrPermitted `json:"prPermitted"`

	// Type アカウントの種類
	Type AccountType `json:"type"`

	// Url アカウントurl
	Url string `json:"url"`
}

// AddContestTeamRequest 新規コンテストチームリクエスト
type AddContestTeamRequest struct {
	// Description チーム情報
	Description string `json:"description"`

	// Link コンテストチームの説明が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name チーム名
	Name string `json:"name"`

	// Result 順位などの結果
	Result *string `json:"result,omitempty"`
}

// AddProjectMembersRequest プロジェクトメンバー追加リクエスト
type AddProjectMembersRequest struct {
	Members []MemberIDWithYearWithSemesterDuration `json:"members"`
}

// Contest コンテスト情報
type Contest struct {
	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// Id コンテストuuid
	Id uuid.UUID `json:"id"`

	// Name コンテスト名
	Name string `json:"name"`
}

// ContestDetail defines model for ContestDetail.
type ContestDetail struct {
	// Description コンテストの説明
	Description string `json:"description"`

	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// Id コンテストuuid
	Id uuid.UUID `json:"id"`

	// Link コンテストの詳細が載っているページへのリンク
	Link string `json:"link"`

	// Name コンテスト名
	Name string `json:"name"`

	// Teams コンテストチーム
	Teams []ContestTeam `json:"teams"`
}

// ContestTeam defines model for ContestTeam.
type ContestTeam struct {
	// Id コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// Members チームメンバーのユーザー情報
	Members []User `json:"members"`

	// Name チーム名
	Name string `json:"name"`

	// Result 順位などの結果
	Result string `json:"result"`
}

// ContestTeamDetail defines model for ContestTeamDetail.
type ContestTeamDetail struct {
	// Description チーム情報
	Description string `json:"description"`

	// Id コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// Link コンテストチームの詳細が載っているページへのリンク
	Link string `json:"link"`

	// Members チームメンバーのUUID
	Members []User `json:"members"`

	// Name チーム名
	Name string `json:"name"`

	// Result 順位などの結果
	Result string `json:"result"`
}

// ContestTeamWithoutMembers コンテストチーム情報(チームメンバーなし)
type ContestTeamWithoutMembers struct {
	// Id コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// Name チーム名
	Name string `json:"name"`

	// Result 順位などの結果
	Result string `json:"result"`
}

// CreateContestRequest 新規コンテストリクエスト
type CreateContestRequest struct {
	// Description コンテスト説明
	Description string `json:"description"`

	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// Link コンテストの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name コンテスト名
	Name string `json:"name"`
}

// CreateProjectRequest 新規プロジェクトリクエスト
type CreateProjectRequest struct {
	// Description プロジェクト説明
	Description string `json:"description"`

	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Link プロジェクトの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name プロジェクト名
	Name string `json:"name"`
}

// Duration イベントやコンテストなどの存続期間
type Duration struct {
	// Since 期間始まり
	Since time.Time `json:"since"`

	// Until 期間終わり
	// untilがなかったらまだ存続している
	Until *time.Time `json:"until,omitempty"`
}

// EditContestRequest コンテスト情報変更リクエスト
type EditContestRequest struct {
	// Description コンテスト説明
	Description *string `json:"description,omitempty"`

	// Duration イベントやコンテストなどの存続期間
	Duration *Duration `json:"duration,omitempty"`

	// Link コンテストの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name コンテスト名
	Name *string `json:"name,omitempty"`
}

// EditContestTeamRequest コンテストチーム情報修正リクエスト
type EditContestTeamRequest struct {
	// Description チーム情報
	Description *string `json:"description,omitempty"`

	// Link コンテストチームの説明が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name チーム名
	Name *string `json:"name,omitempty"`

	// Result 順位などの結果
	Result *string `json:"result,omitempty"`
}

// EditEventRequest イベント情報修正リクエスト
type EditEventRequest struct {
	// EventLevel 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel *EventLevel `json:"eventLevel,omitempty"`
}

// EditProjectRequest プロジェクト変更リクエスト
type EditProjectRequest struct {
	// Description プロジェクト説明
	Description *string `json:"description,omitempty"`

	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration *YearWithSemesterDuration `json:"duration,omitempty"`

	// Link プロジェクトの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty"`

	// Name プロジェクト名
	Name *string `json:"name,omitempty"`
}

// EditUserAccountRequest アカウント変更リクエスト
type EditUserAccountRequest struct {
	// DisplayName 外部アカウントの表示名
	DisplayName *string `json:"displayName,omitempty"`

	// PrPermitted 広報での利用が許可されているかどうか
	PrPermitted *PrPermitted `json:"prPermitted,omitempty"`

	// Type アカウントの種類
	Type *AccountType `json:"type,omitempty"`

	// Url アカウントurl
	Url *string `json:"url,omitempty"`
}

// EditUserRequest ユーザー情報変更リクエスト
type EditUserRequest struct {
	// Bio 自己紹介(biography)
	Bio *string `json:"bio,omitempty"`

	// Check 本名を公開するかどうか
	// true: 公開
	// false: 非公開
	Check *bool `json:"check,omitempty"`
}

// Event イベント情報
type Event struct {
	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// Id イベントuuid
	Id uuid.UUID `json:"id"`

	// Name イベント名
	Name string `json:"name"`
}

// EventDetail defines model for EventDetail.
type EventDetail struct {
	// Description イベント説明
	Description string `json:"description"`

	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// EventLevel 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel EventLevel `json:"eventLevel"`

	// Hostname 主催者
	Hostname []User `json:"hostname"`

	// Id イベントuuid
	Id uuid.UUID `json:"id"`

	// Name イベント名
	Name string `json:"name"`

	// Place 大学、オンラインなどの大まかな場所
	Place string `json:"place"`
}

// EventLevel 公開範囲設定
// 0 イベント企画者の名前を伏せて公開
// 1 全て公開
// 2 外部に非公開
type EventLevel = uint8

// Group 班情報
type Group struct {
	// Id 班uuid
	Id uuid.UUID `json:"id"`

	// Name 班名
	Name string `json:"name"`
}

// GroupDetail defines model for GroupDetail.
type GroupDetail struct {
	// Admin 班管理者
	Admin []User `json:"admin"`

	// Description 班説明
	Description string `json:"description"`

	// Id 班uuid
	Id uuid.UUID `json:"id"`

	// Link 班の詳細が載っているページへのリンク
	Link string `json:"link"`

	// Members 班メンバー
	Members []GroupMember `json:"members"`

	// Name 班名
	Name string `json:"name"`
}

// GroupMember defines model for GroupMember.
type GroupMember struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id ユーザーUUID
	Id uuid.UUID `json:"id"`

	// Name ユーザー名
	Name string `json:"name"`

	// RealName 本名
	RealName string `json:"realName"`
}

// MemberIDWithYearWithSemesterDuration プロジェクトメンバーのユーザーUUID(期間含む)
type MemberIDWithYearWithSemesterDuration struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`
	UserId   uuid.UUID                `json:"userId"`
}

// MemberIDs ユーザーのUUIDの配列
type MemberIDs struct {
	// Members ユーザーのUUIDの配列
	Members []uuid.UUID `json:"members"`
}

// PrPermitted 広報での利用が許可されているかどうか
type PrPermitted = bool

// Project プロジェクト情報
type Project struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// Name プロジェクト名
	Name string `json:"name"`
}

// ProjectDetail defines model for ProjectDetail.
type ProjectDetail struct {
	// Description プロジェクト説明
	Description string `json:"description"`

	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// Link プロジェクトの詳細が載っているページへのリンク
	Link string `json:"link"`

	// Members プロジェクトメンバー
	Members []ProjectMember `json:"members"`

	// Name プロジェクト名
	Name string `json:"name"`
}

// ProjectMember defines model for ProjectMember.
type ProjectMember struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id ユーザーUUID
	Id uuid.UUID `json:"id"`

	// Name ユーザー名
	Name string `json:"name"`

	// RealName 本名
	RealName string `json:"realName"`
}

// Semester 0: 前期
// 1: 後期
type Semester int32

// User ユーザー情報
type User struct {
	// Id ユーザーUUID
	Id uuid.UUID `json:"id"`

	// Name ユーザー名
	Name string `json:"name"`

	// RealName 本名
	RealName string `json:"realName"`
}

// UserAccountState ユーザーアカウント状態
// 0: 凍結
// 1: 有効
// 2: 一時停止
type UserAccountState int32

// UserContest defines model for UserContest.
type UserContest struct {
	// Duration イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// Id コンテストuuid
	Id uuid.UUID `json:"id"`

	// Name コンテスト名
	Name string `json:"name"`

	// Teams コンテストチーム
	Teams []ContestTeamWithoutMembers `json:"teams"`
}

// UserDetail defines model for UserDetail.
type UserDetail struct {
	// Accounts 各種アカウントへのリンク
	Accounts []Account `json:"accounts"`

	// Bio 自己紹介(biography)
	Bio string `json:"bio"`

	// Id ユーザーUUID
	Id uuid.UUID `json:"id"`

	// Name ユーザー名
	Name string `json:"name"`

	// RealName 本名
	RealName string `json:"realName"`

	// State ユーザーアカウント状態
	// 0: 凍結
	// 1: 有効
	// 2: 一時停止
	State UserAccountState `json:"state"`
}

// UserGroup defines model for UserGroup.
type UserGroup struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id 班uuid
	Id uuid.UUID `json:"id"`

	// Name 班名
	Name string `json:"name"`
}

// UserProject defines model for UserProject.
type UserProject struct {
	// Duration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// Id プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// Name プロジェクト名
	Name string `json:"name"`

	// UserDuration 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	UserDuration YearWithSemesterDuration `json:"userDuration"`
}

// YearWithSemester 年度と前期/後期
type YearWithSemester struct {
	// Semester 0: 前期
	// 1: 後期
	Semester Semester `json:"semester"`
	Year     int      `json:"year"`
}

// YearWithSemesterDuration 班やプロジェクトの期間
// 年と前期/後期がある
// untilがなかった場合存続中
type YearWithSemesterDuration struct {
	// Since 年度と前期/後期
	Since YearWithSemester `json:"since"`

	// Until 年度と前期/後期
	Until *YearWithSemester `json:"until,omitempty"`
}

// AccountIdInPath defines model for accountIdInPath.
type AccountIdInPath = uuid.UUID

// ContestIdInPath defines model for contestIdInPath.
type ContestIdInPath = uuid.UUID

// EventIdInPath defines model for eventIdInPath.
type EventIdInPath = uuid.UUID

// GroupIdInPath defines model for groupIdInPath.
type GroupIdInPath = uuid.UUID

// IncludeSuspendedInQuery defines model for includeSuspendedInQuery.
type IncludeSuspendedInQuery = bool

// LimitInQuery defines model for limitInQuery.
type LimitInQuery = int

// NameInQuery defines model for nameInQuery.
type NameInQuery = string

// ProjectIdInPath defines model for projectIdInPath.
type ProjectIdInPath = uuid.UUID

// TeamIdInPath defines model for teamIdInPath.
type TeamIdInPath = uuid.UUID

// UserIdInPath defines model for userIdInPath.
type UserIdInPath = uuid.UUID

// GetUsersParams defines parameters for GetUsers.
type GetUsersParams struct {
	// IncludeSuspended アカウントがアクティブでないユーザーを含めるかどうか
	IncludeSuspended *IncludeSuspendedInQuery `form:"includeSuspended,omitempty" json:"includeSuspended,omitempty" query:"includeSuspended"`

	// Name 指定した文字列がtraP IDに含まれているかどうか
	Name *NameInQuery `form:"name,omitempty" json:"name,omitempty" query:"name"`

	// Limit 取得数の上限
	Limit *LimitInQuery `form:"limit,omitempty" json:"limit,omitempty" query:"limit"`
}

// CreateContestJSONRequestBody defines body for CreateContest for application/json ContentType.
type CreateContestJSONRequestBody = CreateContestRequest

// EditContestJSONRequestBody defines body for EditContest for application/json ContentType.
type EditContestJSONRequestBody = EditContestRequest

// AddContestTeamJSONRequestBody defines body for AddContestTeam for application/json ContentType.
type AddContestTeamJSONRequestBody = AddContestTeamRequest

// EditContestTeamJSONRequestBody defines body for EditContestTeam for application/json ContentType.
type EditContestTeamJSONRequestBody = EditContestTeamRequest

// AddContestTeamMembersJSONRequestBody defines body for AddContestTeamMembers for application/json ContentType.
type AddContestTeamMembersJSONRequestBody = MemberIDs

// EditContestTeamMembersJSONRequestBody defines body for EditContestTeamMembers for application/json ContentType.
type EditContestTeamMembersJSONRequestBody = MemberIDs

// EditEventJSONRequestBody defines body for EditEvent for application/json ContentType.
type EditEventJSONRequestBody = EditEventRequest

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = CreateProjectRequest

// EditProjectJSONRequestBody defines body for EditProject for application/json ContentType.
type EditProjectJSONRequestBody = EditProjectRequest

// DeleteProjectMembersJSONRequestBody defines body for DeleteProjectMembers for application/json ContentType.
type DeleteProjectMembersJSONRequestBody = MemberIDs

// AddProjectMembersJSONRequestBody defines body for AddProjectMembers for application/json ContentType.
type AddProjectMembersJSONRequestBody = AddProjectMembersRequest

// EditUserJSONRequestBody defines body for EditUser for application/json ContentType.
type EditUserJSONRequestBody = EditUserRequest

// AddUserAccountJSONRequestBody defines body for AddUserAccount for application/json ContentType.
type AddUserAccountJSONRequestBody = AddAccountRequest

// EditUserAccountJSONRequestBody defines body for EditUserAccount for application/json ContentType.
type EditUserAccountJSONRequestBody = EditUserAccountRequest
