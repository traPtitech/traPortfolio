# project_members

## Description

プロジェクト所属者テーブル

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `project_members` (
  `id` char(36) NOT NULL,
  `project_id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL,
  `since` datetime(6) DEFAULT NULL,
  `until` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_project_members_project` (`project_id`),
  KEY `fk_project_members_user` (`user_id`),
  CONSTRAINT `fk_project_members_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`),
  CONSTRAINT `fk_project_members_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | char(36) |  | false |  |  |  |
| project_id | char(36) |  | false |  | [projects](projects.md) | プロジェクトUUID |
| user_id | char(36) |  | false |  | [users](users.md) | ユーザーUUID |
| since | datetime(6) |  | true |  |  | プロジェクト所属開始時期 |
| until | datetime(6) |  | true |  |  | プロジェクト所属終了時期 |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| fk_project_members_project | FOREIGN KEY | FOREIGN KEY (project_id) REFERENCES projects (id) |
| fk_project_members_user | FOREIGN KEY | FOREIGN KEY (user_id) REFERENCES users (id) |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| fk_project_members_project | KEY fk_project_members_project (project_id) USING BTREE |
| fk_project_members_user | KEY fk_project_members_user (user_id) USING BTREE |
| PRIMARY | PRIMARY KEY (id) USING BTREE |

## Relations

![er](project_members.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
