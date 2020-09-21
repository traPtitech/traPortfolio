# パッケージ構成

クリーンアーキテクチャ っぽくしてるけど一部そうなってないところもある(interactor とか controller がなくて直接 usecase の handler になってたりする)

```
.
├── dev/    開発用のなにかのつもり
│  └── bin/
├── docs    ドキュメント. swaggerとかtblsのやつとか
│  ├── architecture.md
│  ├── dbschema/
│  └── swagger/
├── domain/     ドメインオブジェクト. tagとかgorm依存になるかもしれない
├── infrastructure/     dbとかフレームワークに関するもの
│  ├── migration/       gomigrateをつかったdbのマイグレーション.
│  │  ├── current.go
│  │  ├── migrate.go
│  │  └── v1.go
│  ├── router.go    echo routerの定義
│  ├── sqlhandler.go    dbの初期化とinterface/database/sqlhandlerの実装
│  ├── wire.go
│  └── wire_gen.go
├── interfaces/     interface層
│  ├── database/    databaseについてのinterface
│  │  └── sqlhandler.go
│  └── repository/      usecases/repositoryの実装
│     └── user_impl.go
├── main.go
└── usecases/       ユースケース層
   ├── handler/     echo ハンドラー
   │  └── api.go
   ├── repository/      レポジトリインターフェス
   │  └── error.go
   ├── service/     ビジネスロジックとかを書くつもり
   │  └── user_service/
   │     └── user.go
   └── usecase/     ユースケースのインターフェース
      ├── ping.go
      └── user.go
```
