
# Golang / Gin / Docker / Docker-Compose の学習用
## はじめに
Golang,Gin,Dockerを使って、簡単なWebサービスを作ります。

そのまま起動しても、Dockerで起動しても使えるようにしてます。

画面は簡易版です。ロジック含め、内容は徐々にブラッシュアアップ予定。

Mac（Intelチップ）で作ってるので、他の環境だとコマンドが違うことがあります。

## 機能
 - ログイン認証

Cognitoを利用。

Cognitoで事前に作成したユーザーを使ってログインする。

　※一応セッション管理しているが、まだ暫定版

 - データ登録/編集/削除

簡単なデータを扱う（少しずつ拡張予定）

## 今後の予定（優先度順）
 - 検索機能の実装
 - 入力値のバリデーション
 - DBをMySQLに変更（別コンテナとして起動させる予定）
 - デザインを整える（Bootstrapを利用する予定）
 - セッション管理（改修）

これをそのままトレーニングの課題にしてもよさそう。

（しばらくmasterはこのままにしておくかも）

## 事前準備
#### Dockerをインストール
https://www.docker.com/get-started/

#### AWS関連
 - アカウントの作成
 - Cognitoでユーザープールのアプリクライアントを作成（シークレットあり）
 - 「認証フローの設定」で「ALLOW_USER_SRP_AUTH」にチェックを入れる
 - Cognitoユーザープールでユーザーを作成（IDとPASSWORDを設定）

#### Golangの環境が無い場合は、以下を参考にGolangをインストール
https://go.dev/doc/install

brewなどでインストールしたい場合は、各自で調査してください。

環境変数（GOPATHとかGOROOT）の設定も必要。

→ 各環境に合わせて設定してください。

#### 最初からやるために、以下を実行
```
rm -f go.mod go.sum
```

#### 必要なパッケージのダウンロード
```
go mod init go_gin_practice
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider
go get github.com/aws/aws-sdk-go/aws
go get github.com/alexrudd/cognito-srp/v4
go get github.com/gin-gonic/gin
go get github.com/jinzhu/gorm
go get github.com/mattn/go-sqlite3
go get github.com/gin-contrib/sessions
```

#### 環境変数の設定（ローカル）
AWSのリージョン（ap-northeast-1など）、CognitoのユーザープールID（ap-northeast-1_xxxxxxxなど）、CognitoのアプリクライアントIDを設定
```
export AWS_REGION=xxxxxxxxxxxxxxx
export COGNITO_USERPOOL_ID=xxxxxxxxxxxxxxx
export COGNITO_APP_CLIENT_ID=xxxxxxxxxxxxxxx
export COGNITO_APP_CLIENT_SECRET=xxxxxxxxxxxxxxx
```

#### envファイルの用意（docker-compose.ymlで読み込む用）
docker-compose.ymlと同じ階層に「variables.env」を作成し、以下を設定
```
AWS_REGION=xxxxxxxxxxxxxxx
COGNITO_USERPOOL_ID=xxxxxxxxxxxxxxx
COGNITO_APP_CLIENT_ID=xxxxxxxxxxxxxxx
COGNITO_APP_CLIENT_SECRET=xxxxxxxxxxxxxxx
```
※秘匿情報なので、「variables.env」は .gitignore に登録して、プッシュされないようにしておく

## 起動方法
```
go run cmd/main.go
```
または
```
docker-compose build & docker-compose up
```
→ そのうちMakefileを作る予定

## 備考
セッション管理方法は複数あるので、どれを採用するかは以下を参照。
https://github.com/gin-contrib/sessions
