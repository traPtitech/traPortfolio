# Development With Remote Containers

## Requirements

- Visial Studio Code
- [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) Extension
- Docker
- Docker Compose

## Usage

1. コマンドパレットを開き(`ctrl+shift+p`)、`Dev Container: Reopen in Container`を選択します
2. コンテナが起動するのを待ちます :coffee:
3. `The "gopls" command is not available. ...`のポップアップが表示されるので、`Install All`を選択します
4. 以降 <http://localhost:1323> と <http://localhost:3001> にアクセスできるようになります

アプリケーションの起動には以下のコマンドを使用してください

```bash
go run main.go -c ./dev/config.yaml
```

## Notes

- プロジェクトルート直下にある`docker-compose.yml`を使ってコンテナを立ち上げています
  - `backend`コンテナは`.devcontainer/docker-compose.yml`で上書きして、この中で開発することができるようになっています
