
## go / gin の学習用
#### はじめに
そのまま起動しても、Dockerで起動しても使えるようにしてます。
内容は徐々にブラッシュアアップ予定

#### 機能
 - ログイン認証
Cognitoを利用。
Cognitoで事前に作成したユーザーを使ってログインする。
 - データ登録/編集/削除
簡単なデータを扱う（そのうち拡張予定）

#### 今後の予定（優先度順）
 - 検索機能の実装
 - 入力値のバリデーション
 - DBをMySQLに変更（別コンテナとして起動させる予定）
 - セッション管理
 - ログアウト機能の実装

#### 事前準備
##### Dockerをインストール
https://www.docker.com/get-started/

##### AWS関連
 - アカウントの作成
 - Cognitoでユーザープールのアプリクライアントを作成（シークレットなし）
 - 「認証フローの設定」で「ALLOW_USER_PASSWORD_AUTH」にチェックを入れる
 - Cognitoユーザープールでユーザーを作成（IDとPASSWORDを設定）

##### Golangの環境が無い場合は、以下を参考にGolangをインストール
https://go.dev/doc/install

brewなどでインストールしたい場合は、ググってください。

環境変数（GOPATHとかGOROOT）の設定も必要になるかも。
→ 最初どうだったか忘れた

##### 最初からやるために、以下を実行
```
rm -f go.mod go.sum
```

##### 必要なパッケージのダウンロード
```
go mod init go_gin_practice
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider
go get github.com/aws/aws-sdk-go/aws
go get github.com/alexrudd/cognito-srp/v4
go get github.com/gin-gonic/gin
go get github.com/jinzhu/gorm
go get github.com/mattn/go-sqlite3
```

#### 環境変数の設定（ローカル）
AWSのリージョン（ap-northease-1など）、CognitoのユーザープールID（ap-northeast-1_xxxxxxxなど）、CognitoのアプリクライアントIDを設定
```
export AWS_REGION=xxxxxxxxxxxxxxx
export COGNITO_USERPOOL_ID=xxxxxxxxxxxxxxx
export COGNITO_APP_CLIENT_ID=xxxxxxxxxxxxxxx
```

#### 環境変数の設定（docker-compose.yml）
```
    environment:
      - AWS_REGION=xxxxxxxxxxxxxxx
      - COGNITO_USERPOOL_ID=xxxxxxxxxxxxxxx
      - COGNITO_APP_CLIENT_ID=xxxxxxxxxxxxxxx
```

#### 起動方法
```
go run cmd/main.go
```
または
```
docker-compose build & docker-compose up
```
→ そのうちMakefileを作る予定

