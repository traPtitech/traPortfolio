// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
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
	N0 UserAccountState = 0
	N1 UserAccountState = 1
	N2 UserAccountState = 2
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
type AddAccountRequest struct {
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
type AddContestTeamRequest struct {
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
type AddProjectMembersRequest struct {
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
	// コンテストの説明
	Description string `json:"description"`

	// イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// コンテストuuid
	Id uuid.UUID `json:"id"`

	// コンテストの詳細が載っているページへのリンク
	Link string `json:"link"`

	// コンテスト名
	Name string `json:"name"`

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
	// チーム情報
	Description string `json:"description"`

	// コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// コンテストチームの詳細が載っているページへのリンク
	Link string `json:"link"`

	// チームメンバーのUUID
	Members []User `json:"members"`

	// チーム名
	Name string `json:"name"`

	// 順位などの結果
	Result string `json:"result"`
}

// ContestTeamWithContestName defines model for ContestTeamWithContestName.
type ContestTeamWithContestName struct {
	// コンテスト名
	ContestName string `json:"contestName"`

	// コンテストチームuuid
	Id uuid.UUID `json:"id"`

	// チーム名
	Name string `json:"name"`

	// 順位などの結果
	Result string `json:"result"`
}

// 新規コンテストリクエスト
type CreateContestRequest struct {
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
type CreateProjectRequest struct {
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
type EditContestRequest struct {
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
type EditContestTeamRequest struct {
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
type EditEventRequest struct {
	// 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel *EventLevel `json:"eventLevel,omitempty"`
}

// プロジェクト変更リクエスト
type EditProjectRequest struct {
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

// アカウント変更リクエスト
type EditUserAccountRequest struct {
	// 外部アカウントの表示名
	DisplayName *string `json:"displayName,omitempty"`

	// 広報での利用が許可されているかどうか
	PrPermitted *PrPermitted `json:"prPermitted,omitempty"`

	// アカウントの種類
	Type *AccountType `json:"type,omitempty"`

	// アカウントurl
	Url *string `json:"url,omitempty" validate:"omitempty,url"`
}

// ユーザー情報変更リクエスト
type EditUserRequest struct {
	// 自己紹介(biography)
	Bio *string `json:"bio,omitempty"`

	// 本名を公開するかどうか
	// true: 公開
	// false: 非公開
	Check *bool `json:"check,omitempty"`
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
	// イベント説明
	Description string `json:"description"`

	// イベントやコンテストなどの存続期間
	Duration Duration `json:"duration"`

	// 公開範囲設定
	// 0 イベント企画者の名前を伏せて公開
	// 1 全て公開
	// 2 外部に非公開
	EventLevel EventLevel `json:"eventLevel"`

	// 主催者
	Hostname []User `json:"hostname"`

	// イベントuuid
	Id uuid.UUID `json:"id"`

	// イベント名
	Name string `json:"name"`

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
	// 班管理者
	Admin []User `json:"admin"`

	// 班説明
	Description string `json:"description"`

	// 班uuid
	Id uuid.UUID `json:"id"`

	// 班の詳細が載っているページへのリンク
	Link string `json:"link"`

	// 班メンバー
	Members []GroupMember `json:"members"`

	// 班名
	Name string `json:"name"`
}

// GroupMember defines model for GroupMember.
type GroupMember struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// ユーザーUUID
	Id uuid.UUID `json:"id"`

	// ユーザー名
	Name string `json:"name"`

	// 本名
	RealName string `json:"realName"`
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
type PrPermitted = bool

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
	// プロジェクト説明
	Description string `json:"description"`

	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// プロジェクトの詳細が載っているページへのリンク
	Link string `json:"link"`

	// プロジェクトメンバー
	Members []ProjectMember `json:"members"`

	// プロジェクト名
	Name string `json:"name"`
}

// ProjectMember defines model for ProjectMember.
type ProjectMember struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// ユーザーUUID
	Id uuid.UUID `json:"id"`

	// ユーザー名
	Name string `json:"name"`

	// 本名
	RealName string `json:"realName"`
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
	// 各種アカウントへのリンク
	Accounts []Account `json:"accounts"`

	// 自己紹介(biography)
	Bio string `json:"bio"`

	// ユーザーUUID
	Id uuid.UUID `json:"id"`

	// ユーザー名
	Name string `json:"name"`

	// 本名
	RealName string `json:"realName"`

	// ユーザーアカウント状態
	// 0: 凍結
	// 1: 有効
	// 2: 一時停止
	State UserAccountState `json:"state"`
}

// UserGroup defines model for UserGroup.
type UserGroup struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// 班uuid
	Id uuid.UUID `json:"id"`

	// 班名
	Name string `json:"name"`
}

// UserProject defines model for UserProject.
type UserProject struct {
	// 班やプロジェクトの期間
	// 年と前期/後期がある
	// untilがなかった場合存続中
	Duration YearWithSemesterDuration `json:"duration"`

	// プロジェクトuuid
	Id uuid.UUID `json:"id"`

	// プロジェクト名
	Name string `json:"name"`

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
type IncludeSuspendedInQuery = bool

// NameInQuery defines model for nameInQuery.
type NameInQuery = string

// ProjectIdInPath defines model for projectIdInPath.
type ProjectIdInPath uuid.UUID

// TeamIdInPath defines model for teamIdInPath.
type TeamIdInPath uuid.UUID

// UserIdInPath defines model for userIdInPath.
type UserIdInPath uuid.UUID

// CreateContestJSONBody defines parameters for CreateContest.
type CreateContestJSONBody = CreateContestRequest

// EditContestJSONBody defines parameters for EditContest.
type EditContestJSONBody = EditContestRequest

// AddContestTeamJSONBody defines parameters for AddContestTeam.
type AddContestTeamJSONBody = AddContestTeamRequest

// EditContestTeamJSONBody defines parameters for EditContestTeam.
type EditContestTeamJSONBody = EditContestTeamRequest

// AddContestTeamMembersJSONBody defines parameters for AddContestTeamMembers.
type AddContestTeamMembersJSONBody = MemberIDs

// EditContestTeamMembersJSONBody defines parameters for EditContestTeamMembers.
type EditContestTeamMembersJSONBody = MemberIDs

// EditEventJSONBody defines parameters for EditEvent.
type EditEventJSONBody = EditEventRequest

// CreateProjectJSONBody defines parameters for CreateProject.
type CreateProjectJSONBody = CreateProjectRequest

// EditProjectJSONBody defines parameters for EditProject.
type EditProjectJSONBody = EditProjectRequest

// DeleteProjectMembersJSONBody defines parameters for DeleteProjectMembers.
type DeleteProjectMembersJSONBody = MemberIDs

// AddProjectMembersJSONBody defines parameters for AddProjectMembers.
type AddProjectMembersJSONBody = AddProjectMembersRequest

// GetUsersParams defines parameters for GetUsers.
type GetUsersParams struct {
	// アカウントがアクティブでないユーザーを含めるかどうか
	IncludeSuspended *IncludeSuspendedInQuery `form:"includeSuspended,omitempty" json:"includeSuspended,omitempty" query:"includeSuspended"`

	// 指定した文字列がtraP IDに含まれているかどうか
	Name *NameInQuery `form:"name,omitempty" json:"name,omitempty" query:"name"`
}

// EditUserJSONBody defines parameters for EditUser.
type EditUserJSONBody = EditUserRequest

// AddUserAccountJSONBody defines parameters for AddUserAccount.
type AddUserAccountJSONBody = AddAccountRequest

// EditUserAccountJSONBody defines parameters for EditUserAccount.
type EditUserAccountJSONBody = EditUserAccountRequest

// CreateContestJSONRequestBody defines body for CreateContest for application/json ContentType.
type CreateContestJSONRequestBody = CreateContestJSONBody

// EditContestJSONRequestBody defines body for EditContest for application/json ContentType.
type EditContestJSONRequestBody = EditContestJSONBody

// AddContestTeamJSONRequestBody defines body for AddContestTeam for application/json ContentType.
type AddContestTeamJSONRequestBody = AddContestTeamJSONBody

// EditContestTeamJSONRequestBody defines body for EditContestTeam for application/json ContentType.
type EditContestTeamJSONRequestBody = EditContestTeamJSONBody

// AddContestTeamMembersJSONRequestBody defines body for AddContestTeamMembers for application/json ContentType.
type AddContestTeamMembersJSONRequestBody = AddContestTeamMembersJSONBody

// EditContestTeamMembersJSONRequestBody defines body for EditContestTeamMembers for application/json ContentType.
type EditContestTeamMembersJSONRequestBody = EditContestTeamMembersJSONBody

// EditEventJSONRequestBody defines body for EditEvent for application/json ContentType.
type EditEventJSONRequestBody = EditEventJSONBody

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = CreateProjectJSONBody

// EditProjectJSONRequestBody defines body for EditProject for application/json ContentType.
type EditProjectJSONRequestBody = EditProjectJSONBody

// DeleteProjectMembersJSONRequestBody defines body for DeleteProjectMembers for application/json ContentType.
type DeleteProjectMembersJSONRequestBody = DeleteProjectMembersJSONBody

// AddProjectMembersJSONRequestBody defines body for AddProjectMembers for application/json ContentType.
type AddProjectMembersJSONRequestBody = AddProjectMembersJSONBody

// EditUserJSONRequestBody defines body for EditUser for application/json ContentType.
type EditUserJSONRequestBody = EditUserJSONBody

// AddUserAccountJSONRequestBody defines body for AddUserAccount for application/json ContentType.
type AddUserAccountJSONRequestBody = AddUserAccountJSONBody

// EditUserAccountJSONRequestBody defines body for EditUserAccount for application/json ContentType.
type EditUserAccountJSONRequestBody = EditUserAccountJSONBody
