// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.10.1 DO NOT EDIT.
package handler

import (
	"time"

	"github.com/gofrs/uuid"
)

// Defines values for EventLevel.
const (
	EventLevelN0 EventLevel = 0

	EventLevelN1 EventLevel = 1

	EventLevelN2 EventLevel = 2
)

// Defines values for Semester.
const (
	SemesterN0 Semester = 0

	SemesterN1 Semester = 1
)

// Defines values for UserAccountState.
const (
	UserAccountStateN0 UserAccountState = 0

	UserAccountStateN1 UserAccountState = 1

	UserAccountStateN2 UserAccountState = 2
)

// アカウントへのリンク
type Account struct {
	// 外部アカウントの表示名
	DisplayName string `json:"displayName"`

	// アカウントUUID
	Id uuid.UUID `json:"id"`

	// 広報での利用が許可されているかどうか
	PrPermitted PrPermitted `json:"prPermitted"`

	// アカウントの種類
	Type AccountType `json:"type"`

	// アカウントurl
	Url string `json:"url"`
}

// アカウントの種類
type AccountType int64

// 新規アカウントリクエスト
type AddAccount struct {
	// 外部アカウントの表示名
	DisplayName string `json:"displayName"`

	// 広報での利用が許可されているかどうか
	PrPermitted PrPermitted `json:"prPermitted"`

	// アカウントの種類
	Type AccountType `json:"type"`

	// アカウントurl
	Url string `json:"url" validate:"url"`
}

// 新規コンテストチームリクエスト
type AddContestTeam struct {
	// チーム情報
	Description string `json:"description"`

	// コンテストチームの説明が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// チーム名
	Name string `json:"name"`

	// 順位などの結果
	Result *string `json:"result,omitempty"`
}

// プロジェクトメンバー追加リクエスト
type AddProjectMembers struct {
	Members []MemberIDWithYearWithSemesterDuration `json:"members"`
}

// コンテスト情報
type Contest struct {
	// イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// コンテストuuid
	Id uuid.UUID `json:"id"`

	// コンテスト名
	Name string `json:"name"`
}

// ContestDetail defines model for ContestDetail.
type ContestDetail struct {
	// Embedded struct due to allOf(#/components/schemas/Contest)
	Contest `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// コンテストの説明
	Description string `json:"description"`

	// コンテストの詳細が載っているページへのリンク
	Link string `json:"link"`

	// コンテストチーム
	Teams []ContestTeam `json:"teams"`
}

// コンテストチーム情報
type ContestTeam struct {
	// コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// チーム名
	Name string `json:"name"`

	// 順位などの結果
	Result string `json:"result"`
}

// ContestTeamDetail defines model for ContestTeamDetail.
type ContestTeamDetail struct {
	// Embedded struct due to allOf(#/components/schemas/ContestTeam)
	ContestTeam `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// チーム情報
	Description string `json:"description"`

	// コンテストチームの詳細が載っているページへのリンク
	Link string `json:"link"`

	// チームメンバーのUUID
	Members []User `json:"members"`
}

// ContestTeamWithContestName defines model for ContestTeamWithContestName.
type ContestTeamWithContestName struct {
	// Embedded struct due to allOf(#/components/schemas/ContestTeam)
	ContestTeam `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// コンテスト名
	ContestName string `json:"contestName"`
}

// 新規コンテストリクエスト
type CreateContest struct {
	// コンテスト説明
	Description string `json:"description"`

	// イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// コンテストの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// コンテスト名
	Name string `json:"name"`
}

// 新規プロジェクトリクエスト
type CreateProject struct {
	// プロジェクト説明
	Description string `json:"description"`

	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// プロジェクトの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// プロジェクト名
	Name string `json:"name"`
}

// イベントやコンテストなどの存続期間
type Duration struct {
	// 期間始まり
	Since time.Time `json:"since"`

	// 期間終わり
	// untilがなかったらまだ存続している
	Until *time.Time `json:"until,omitempty"`
}

// コンテスト情報変更リクエスト
type EditContest struct {
	// コンテスト説明
	Description *string `json:"description,omitempty"`

	// イベントやコンテストなどの存続期間
	Duration *Duration `json:"duration,omitempty"`

	// コンテストの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// コンテスト名
	Name *string `json:"name,omitempty"`
}

// コンテストチーム情報修正リクエスト
type EditContestTeam struct {
	// チーム情報
	Description *string `json:"description,omitempty"`

	// コンテストチームの説明が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// チーム名
	Name *string `json:"name,omitempty"`

	// 順位などの結果
	Result *string `json:"result,omitempty"`
}

// イベント情報修正リクエスト
type EditEvent struct {
	// 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel *EventLevel `json:"eventLevel,omitempty"`
}

// プロジェクト変更リクエスト
type EditProject struct {
	// プロジェクト説明
	Description *string `json:"description,omitempty"`

	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration *YearWithSemesterDuration `json:"duration,omitempty"`

	// プロジェクトの詳細が載っているページへのリンク
	Link *string `json:"link,omitempty" validate:"url"`

	// プロジェクト名
	Name *string `json:"name,omitempty"`
}

// ユーザー情報変更リクエスト
type EditUser struct {
	// 自己紹介(biography)
	Bio *string `json:"bio,omitempty"`

	// 本名を公開するかどうか
	// true: 公開
	// false: 非公開
	Check *bool `json:"check,omitempty"`
}

// アカウント変更リクエスト
type EditUserAccount struct {
	// 外部アカウントの表示名
	DisplayName *string `json:"displayName,omitempty"`

	// 広報での利用が許可されているかどうか
	PrPermitted *PrPermitted `json:"prPermitted,omitempty"`

	// アカウントの種類
	Type *AccountType `json:"type,omitempty"`

	// アカウントurl
	Url *string `json:"url,omitempty" validate:"omitempty,url"`
}

// イベント情報
type Event struct {
	// イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// イベントuuid
	Id uuid.UUID `json:"id"`

	// イベント名
	Name string `json:"name"`
}

// EventDetail defines model for EventDetail.
type EventDetail struct {
	// Embedded struct due to allOf(#/components/schemas/Event)
	Event `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// イベント説明
	Description string `json:"description"`

	// 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel EventLevel `json:"eventLevel"`

	// 主催者
	Hostname []User `json:"hostname"`

	// 大学、オンラインなどの大まかな場所
	Place string `json:"place"`
}

// 公開範囲設定
// 0 イベント企画者の名前を伏せて公開
// 1 全て公開
// 2 外部に非公開
type EventLevel int

// 班情報
type Group struct {
	// 班uuid
	Id uuid.UUID `json:"id"`

	// 班名
	Name string `json:"name"`
}

// GroupDetail defines model for GroupDetail.
type GroupDetail struct {
	// Embedded struct due to allOf(#/components/schemas/Group)
	Group `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 班説明
	Description string `json:"description"`

	// ユーザー情報
	Leader User `json:"leader"`

	// 班の詳細が載っているページへのリンク
	Link string `json:"link"`

	// 班メンバー
	Members []GroupMember `json:"members"`
}

// GroupMember defines model for GroupMember.
type GroupMember struct {
	// Embedded struct due to allOf(#/components/schemas/User)
	User `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`
}

// プロジェクトメンバーのユーザーUUID(期間含む)
type MemberIDWithYearWithSemesterDuration struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`
	UserId   uuid.UUID                `json:"userId"`
}

// ユーザーのUUIDの配列
type MemberIDs struct {
	// ユーザーのUUIDの配列
	Members []uuid.UUID `json:"members"`
}

// 広報での利用が許可されているかどうか
type PrPermitted bool

// プロジェクト情報
type Project struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// プロジェクト名
	Name string `json:"name"`
}

// ProjectDetail defines model for ProjectDetail.
type ProjectDetail struct {
	// Embedded struct due to allOf(#/components/schemas/Project)
	Project `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// プロジェクト説明
	Description string `json:"description"`

	// プロジェクトの詳細が載っているページへのリンク
	Link string `json:"link"`

	// プロジェクトメンバー
	Members []ProjectMember `json:"members"`
}

// ProjectMember defines model for ProjectMember.
type ProjectMember struct {
	// Embedded struct due to allOf(#/components/schemas/User)
	User `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`
}

// 0: 前期
// 1: 後期
type Semester int32

// ユーザー情報
type User struct {
	// ユーザーUUID
	Id uuid.UUID `json:"id"`

	// ユーザー名
	Name string `json:"name"`

	// 本名
	RealName string `json:"realName"`
}

// ユーザーアカウント状態
// 0: 凍結
// 1: 有効
// 2: 一時停止
type UserAccountState int32

// UserDetail defines model for UserDetail.
type UserDetail struct {
	// Embedded struct due to allOf(#/components/schemas/User)
	User `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 各種アカウントへのリンク
	Accounts []Account `json:"accounts"`

	// 自己紹介(biography)
	Bio string `json:"bio"`

	// ユーザーアカウント状態
	// 0: 凍結
	// 1: 有効
	// 2: 一時停止
	State UserAccountState `json:"state"`
}

// UserGroup defines model for UserGroup.
type UserGroup struct {
	// Embedded struct due to allOf(#/components/schemas/Group)
	Group `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`
}

// UserProject defines model for UserProject.
type UserProject struct {
	// Embedded struct due to allOf(#/components/schemas/Project)
	Project `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	UserDuration YearWithSemesterDuration `json:"userDuration"`
}

// 年度と前期/後期
type YearWithSemester struct {
	// 0: 前期
	// 1: 後期
	Semester Semester `json:"semester"`
	Year     int      `json:"year"`
}

// 班やプロジェクトの期間
// 年と前期/後期がある
// untilがなかった場合存続中
type YearWithSemesterDuration struct {
	// 年度と前期/後期
	Since YearWithSemester `json:"since"`

	// 年度と前期/後期
	Until *YearWithSemester `json:"until,omitempty"`
}

// AccountIdInPath defines model for accountIdInPath.
type AccountIdInPath uuid.UUID

// ContestIdInPath defines model for contestIdInPath.
type ContestIdInPath uuid.UUID

// EventIdInPath defines model for eventIdInPath.
type EventIdInPath uuid.UUID

// GroupIdInPath defines model for groupIdInPath.
type GroupIdInPath uuid.UUID

// IncludeSuspendedInQuery defines model for includeSuspendedInQuery.
type IncludeSuspendedInQuery bool

// NameInQuery defines model for nameInQuery.
type NameInQuery string

// ProjectIdInPath defines model for projectIdInPath.
type ProjectIdInPath uuid.UUID

// TeamIdInPath defines model for teamIdInPath.
type TeamIdInPath uuid.UUID

// UserIdInPath defines model for userIdInPath.
type UserIdInPath uuid.UUID

// CreateContestJSONBody defines parameters for CreateContest.
type CreateContestJSONBody CreateContest

// EditContestJSONBody defines parameters for EditContest.
type EditContestJSONBody EditContest

// AddContestTeamJSONBody defines parameters for AddContestTeam.
type AddContestTeamJSONBody AddContestTeam

// EditContestTeamJSONBody defines parameters for EditContestTeam.
type EditContestTeamJSONBody EditContestTeam

// DeleteContestTeamMembersJSONBody defines parameters for DeleteContestTeamMembers.
type DeleteContestTeamMembersJSONBody MemberIDs

// AddContestTeamMembersJSONBody defines parameters for AddContestTeamMembers.
type AddContestTeamMembersJSONBody MemberIDs

// EditEventJSONBody defines parameters for EditEvent.
type EditEventJSONBody EditEvent

// CreateProjectJSONBody defines parameters for CreateProject.
type CreateProjectJSONBody CreateProject

// EditProjectJSONBody defines parameters for EditProject.
type EditProjectJSONBody EditProject

// DeleteProjectMembersJSONBody defines parameters for DeleteProjectMembers.
type DeleteProjectMembersJSONBody MemberIDs

// AddProjectMembersJSONBody defines parameters for AddProjectMembers.
type AddProjectMembersJSONBody AddProjectMembers

// GetUsersParams defines parameters for GetUsers.
type GetUsersParams struct {
	// アカウントがアクティブでないユーザーを含めるかどうか
	IncludeSuspended *IncludeSuspendedInQuery `json:"includeSuspended,omitempty" query:"includeSuspended"`

	// 指定した文字列がtraP IDに含まれているかどうか
	Name *NameInQuery `json:"name,omitempty" query:"name"`
}

// EditUserJSONBody defines parameters for EditUser.
type EditUserJSONBody EditUser

// AddUserAccountJSONBody defines parameters for AddUserAccount.
type AddUserAccountJSONBody AddAccount

// EditUserAccountJSONBody defines parameters for EditUserAccount.
type EditUserAccountJSONBody EditUserAccount

// CreateContestJSONRequestBody defines body for CreateContest for application/json ContentType.
type CreateContestJSONRequestBody CreateContestJSONBody

// EditContestJSONRequestBody defines body for EditContest for application/json ContentType.
type EditContestJSONRequestBody EditContestJSONBody

// AddContestTeamJSONRequestBody defines body for AddContestTeam for application/json ContentType.
type AddContestTeamJSONRequestBody AddContestTeamJSONBody

// EditContestTeamJSONRequestBody defines body for EditContestTeam for application/json ContentType.
type EditContestTeamJSONRequestBody EditContestTeamJSONBody

// DeleteContestTeamMembersJSONRequestBody defines body for DeleteContestTeamMembers for application/json ContentType.
type DeleteContestTeamMembersJSONRequestBody DeleteContestTeamMembersJSONBody

// AddContestTeamMembersJSONRequestBody defines body for AddContestTeamMembers for application/json ContentType.
type AddContestTeamMembersJSONRequestBody AddContestTeamMembersJSONBody

// EditEventJSONRequestBody defines body for EditEvent for application/json ContentType.
type EditEventJSONRequestBody EditEventJSONBody

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody CreateProjectJSONBody

// EditProjectJSONRequestBody defines body for EditProject for application/json ContentType.
type EditProjectJSONRequestBody EditProjectJSONBody

// DeleteProjectMembersJSONRequestBody defines body for DeleteProjectMembers for application/json ContentType.
type DeleteProjectMembersJSONRequestBody DeleteProjectMembersJSONBody

// AddProjectMembersJSONRequestBody defines body for AddProjectMembers for application/json ContentType.
type AddProjectMembersJSONRequestBody AddProjectMembersJSONBody

// EditUserJSONRequestBody defines body for EditUser for application/json ContentType.
type EditUserJSONRequestBody EditUserJSONBody

// AddUserAccountJSONRequestBody defines body for AddUserAccount for application/json ContentType.
type AddUserAccountJSONRequestBody AddUserAccountJSONBody

// EditUserAccountJSONRequestBody defines body for EditUserAccount for application/json ContentType.
type EditUserAccountJSONRequestBody EditUserAccountJSONBody
