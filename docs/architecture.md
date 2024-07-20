# アーキテクチャ

クリーンアーキテクチャを参考に構成

## パッケージ構成

```bash
.
├── main.go # エントリーポイント
├── injector.go # 依存性注入
└── internal
   ├── domain # 他層に依存しないドメインオブジェクトを格納する
   ├── usecases # アプリケーションの具体的な操作を表現する (domain層に依存)
   │  └── repository # リポジトリ操作に関するインターフェイスの定義
   ├── handler # Echoによるハンドラー&ルーティング (domain層、usecases層に依存)
   │  └── schema # OpenAPIを基に自動生成されたAPIスキーマ
   ├── infrastructure # 外部APIやDBへのアクセス (domain層、usecases層に依存)
   │  ├── external # 外部APIへのアクセス
   │  ├── migration # DBのマイグレーション
   │  └── repository # DBへのアクセス (usecases/repository の実装)
   │     └── model # GORMのDBスキーマ
   └── pkgs # 汎用パッケージ
      ├── config # 設定ファイル、環境変数を読み込み管理する
      ├── mockdata # E2Eテスト、結合テストで用いるサンプルデータを格納 (消したい)
      ├── optional # optionalな値を扱うためのパッケージ
      ├── random # テストで用いる乱数生成パッケージ
      └── testutils # テストで用いるユーティリティ
```

## 依存関係 (WIP)

ユーザー周りに限って紹介

```mermaid
classDiagram
    class User["domain.User"] {
        ID       uuid.UUID
        Name     string
        realName string
        Check    bool
    }

    class IUS["usecases.service.UserService"] {
        <<interface>>
        GetUsers(...) ([]*domain.User, error)
        ...
    }
    class US["usecases.service.userService"] {
        user repository.UserRepository
        event repository.EventRepository
    }

    class IUR["usecases.repository.UserRepository"] {
        <<interface>>
        GetUsers(...) ([]*domain.User, error)
        ...
    }

    IUS --> User
    IUR --> User
    US --|> IUS
    US --> IUR

    class API["interfaces.handler.API"] {
        Ping *PingHandler
        User *UserHandler
        ...
    }
    class UH["interfaces.handler.UserHandler"] {
        s service.UserService
    }

    API --> UH
    UH --> IUS

    class UR["infrastructure.repository.UserRepository"] {
        h      *gorm.DB
        portal external.PortalAPI
        traQ   external.TraQAPI
    }

    class IPortalAPI["infrastructure.external.PortalAPI"] {
        <<interface>>
        ...
    }
    class PortalAPI["infrastructure.external.portalAPI"]

    class ITraQAPI["infrastructure.external.TraQAPI"] {
        <<interface>>
        ...
    }
    class TraQAPI["infrastructure.external.traQAPI"]

    UR --|> IUR
    PortalAPI --|> IPortalAPI
    TraQAPI --|> ITraQAPI
    UR --> IPortalAPI
    UR --> ITraQAPI
```
