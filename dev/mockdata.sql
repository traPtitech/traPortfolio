CREATE DATABASE IF NOT EXISTS `portfolio`;

USE `portfolio`;

-- users --
INSERT INTO
  `users` (
    `id`,
    `description`,
    `check`,
    `name`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    "11111111-1111-1111-1111-111111111111",
    "I am user1.",
    1,
    "user1",
    NOW(),
    NOW()
  ),
  (
    "22222222-2222-2222-2222-222222222222",
    "I am user2.",
    1,
    "user2",
    NOW(),
    NOW()
  ),
  (
    "33333333-3333-3333-3333-333333333333",
    "I am lolico.",
    0,
    "lolico",
    NOW(),
    NOW()
  );

-- accounts --
INSERT INTO
  `accounts` (
    `id`,
    `type`,
    `name`,
    `url`,
    `user_id`,
    `check`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    'd834e180-2af9-4cfe-838a-8a3930666490',
    0,
    'nhwJy',
    '6fGaE',
    '11111111-1111-1111-1111-111111111111',
    true,
    NOW(),
    NOW()
  );

-- contests --
INSERT INTO
  `contests` (
    `id`,
    `name`,
    `description`,
    `link`,
    `since`,
    `until`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    '08eec963-0f29-48d1-929f-004cb67d8ce6',
    'iDHIMG0FpXzp',
    'HMYZL6fGaEFP1yynhwJyzAHyfjXUlr',
    'https://mBZH21eAer7qyNEUz',
    '2006-01-02 15:04:05',
    '2006-01-02 15:04:05',
    NOW(),
    NOW()
  );

-- contest_teams --
INSERT INTO
  `contest_teams` (
    `id`,
    `contest_id`,
    `name`,
    `description`,
    `result`,
    `link`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    '1543039b-424b-4d67-8c22-774cd9b7cd15',
    '08eec963-0f29-48d1-929f-004cb67d8ce6',
    'iDHIMG0FpXzp',
    'K3aa',
    'HMYZL6fGaEFP1yynhwJyzAHyfjXUlr',
    'https://mBZH21eAer7qyNEUz',
    NOW(),
    NOW()
  );

-- contest_team_user_belongings --
INSERT INTO
  `contest_team_user_belongings` (
    `team_id`,
    `user_id`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    '1543039b-424b-4d67-8c22-774cd9b7cd15',
    '11111111-1111-1111-1111-111111111111',
    '2021-11-17 18:27:04.061',
    '2021-11-17 18:27:04.061'
  );

-- event_level_relations --
INSERT INTO
  `event_level_relations` (
    `id`,
    `level`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    '1',
    NOW(),
    NOW()
  ),
  (
    "11111111-1111-1111-1111-111111111111",
    '0',
    NOW(),
    NOW()
  );

-- groups --
INSERT INTO
  `groups` (
    `group_id`,
    `name`,
    `link`,
    `leader`
  )
VALUES
  (
    'f86db5ec-dc02-4885-aa0a-732bb229a1b5',
    'SysAd班',
    'http://http://Jy',
    '052ace90-7a66-4770-9b86-95fc39e0f434'
  );

-- group_user_belongings --
INSERT INTO
  `group_user_belongings` (
    `user_id`,
    `group_id`,
    `since_year`,
    `since_semester`,
    `until_year`,
    `until_semester`
  )
VALUES
  (
    '11111111-1111-1111-1111-111111111111',
    'f86db5ec-dc02-4885-aa0a-732bb229a1b5',
    2021,
    0,
    2021,
    1
  );

-- projects --
INSERT INTO
  `projects` (
    `id`,
    `name`,
    `description`,
    `link`,
    `since`,
    `until`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    '061ea8ba-ca44-48ed-8ace-54752b5e88f0',
    'K3',
    'q7DOT9tGTOWeLrQKjLxzI3ivH9S71l',
    'http://cbW9wNbeHV6IkP252Z',
    '2021-08-01 00:00:00',
    '2021-12-01 00:00:00',
    NOW(),
    NOW()
  );

-- project_members --
INSERT INTO
  `project_members` (
    `id`,
    `project_id`,
    `user_id`,
    `since`,
    `until`,
    `created_at`,
    `updated_at`
  )
VALUES
  (
    '43f88b23-2d37-4e51-b825-e204c34a6e78',
    '061ea8ba-ca44-48ed-8ace-54752b5e88f0',
    '11111111-1111-1111-1111-111111111111',
    '2021-08-01 00:00:00',
    '2021-12-01 00:00:00',
    NOW(),
    NOW()
  );
