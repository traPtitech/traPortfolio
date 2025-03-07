# portfolio

## Tables

| Name | Columns | Comment | Type |
| ---- | ------- | ------- | ---- |
| [accounts](accounts.md) | 7 | Twitterなどのアカウントテーブル | BASE TABLE |
| [contests](contests.md) | 8 | コンテストテーブル | BASE TABLE |
| [contest_teams](contest_teams.md) | 8 | コンテスト参加チームテーブル | BASE TABLE |
| [contest_team_user_belongings](contest_team_user_belongings.md) | 4 | コンテストチームとユーザー関係テーブル | BASE TABLE |
| [event_level_relations](event_level_relations.md) | 4 | knoQイベントと公開レベルの関係テーブル | BASE TABLE |
| [groups](groups.md) | 6 | グループテーブル | BASE TABLE |
| [group_user_admins](group_user_admins.md) | 4 | グループと管理者関係テーブル | BASE TABLE |
| [group_user_belongings](group_user_belongings.md) | 8 | グループとユーザー関係テーブル | BASE TABLE |
| [migrations](migrations.md) | 1 | gormigrate用のデータベースバージョンテーブル | BASE TABLE |
| [projects](projects.md) | 10 | プロジェクトテーブル | BASE TABLE |
| [project_members](project_members.md) | 8 | プロジェクト所属者テーブル | BASE TABLE |
| [users](users.md) | 7 | ユーザーテーブル | BASE TABLE |

## Relations

```mermaid
erDiagram

"accounts" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"contest_teams" }o--|| "contests" : "FOREIGN KEY (contest_id) REFERENCES contests (id)"
"contest_team_user_belongings" }o--|| "contest_teams" : "FOREIGN KEY (team_id) REFERENCES contest_teams (id)"
"contest_team_user_belongings" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"group_user_admins" }o--|| "groups" : "FOREIGN KEY (group_id) REFERENCES groups (group_id)"
"group_user_belongings" }o--|| "groups" : "FOREIGN KEY (group_id) REFERENCES groups (group_id)"
"group_user_belongings" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"project_members" }o--|| "projects" : "FOREIGN KEY (project_id) REFERENCES projects (id)"
"project_members" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"

"accounts" {
  char_36_ id PK
  tinyint_1_ type
  varchar_256_ name
  text url
  char_36_ user_id FK
  datetime_6_ created_at
  datetime_6_ updated_at
}
"contests" {
  char_36_ id PK
  varchar_128_ name
  text description
  text link
  datetime_6_ since
  datetime_6_ until
  datetime_6_ created_at
  datetime_6_ updated_at
}
"contest_teams" {
  char_36_ id PK
  char_36_ contest_id FK
  varchar_128_ name
  text description
  text result
  text link
  datetime_6_ created_at
  datetime_6_ updated_at
}
"contest_team_user_belongings" {
  char_36_ team_id PK
  char_36_ user_id PK
  datetime_6_ created_at
  datetime_6_ updated_at
}
"event_level_relations" {
  char_36_ id PK
  tinyint_3__unsigned level
  datetime_6_ created_at
  datetime_6_ updated_at
}
"groups" {
  char_36_ group_id PK
  varchar_32_ name
  text link
  text description
  datetime_6_ created_at
  datetime_6_ updated_at
}
"group_user_admins" {
  char_36_ user_id PK
  char_36_ group_id PK
  datetime_6_ created_at
  datetime_6_ updated_at
}
"group_user_belongings" {
  char_36_ user_id PK
  char_36_ group_id PK
  smallint_4_ since_year
  tinyint_1_ since_semester
  smallint_4_ until_year
  tinyint_1_ until_semester
  datetime_6_ created_at
  datetime_6_ updated_at
}
"migrations" {
  varchar_255_ id PK
}
"projects" {
  char_36_ id PK
  varchar_128_ name
  text description
  text link
  smallint_4_ since_year
  tinyint_1_ since_semester
  smallint_4_ until_year
  tinyint_1_ until_semester
  datetime_6_ created_at
  datetime_6_ updated_at
}
"project_members" {
  char_36_ project_id PK
  char_36_ user_id PK
  smallint_4_ since_year
  tinyint_1_ since_semester
  smallint_4_ until_year
  tinyint_1_ until_semester
  datetime_6_ created_at
  datetime_6_ updated_at
}
"users" {
  char_36_ id PK
  text description
  tinyint_1_ check
  varchar_32_ name
  tinyint_1_ state
  datetime_6_ created_at
  datetime_6_ updated_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
