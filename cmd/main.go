package main

import (
	"context"
	"fmt"
	"go_gin_practice/database"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	database.Init()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	})

	router.GET("/add", func(ctx *gin.Context) {
		products := database.GetAll()
		ctx.HTML(200, "index.html", gin.H{
			"products": products,
		})
	})

	router.POST("/login", func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		password := ctx.PostForm("password")

		err := login(id, password)
		if err != nil {
			fmt.Println(err)
			ctx.HTML(200, "login.html", gin.H{"message": "ログイン失敗"})
		}
		ctx.Redirect(302, "/add")
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		url := ctx.PostForm("url")
		memo := ctx.PostForm("memo")
		database.Insert(title, url, memo)
		ctx.Redirect(302, "/add")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		product := database.GetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"product": product})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		title := ctx.PostForm("title")
		url := ctx.PostForm("url")
		memo := ctx.PostForm("memo")
		database.Update(id, title, url, memo)
		ctx.Redirect(302, "/add")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		product := database.GetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"product": product})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		database.Delete(id)
		ctx.Redirect(302, "/add")

	})

	router.Run()
}

// どこかでセッション管理しないといけないけど、とりあえず認証だけ通す
func login(id string, password string) error {

	// cognitosrpはあまりメンテナンスされてなさそうなので、自分で処理を書いた方がいいかも
	csrp, _ := cognitosrp.NewCognitoSRP(id, password, os.Getenv("COGNITO_USERPOOL_ID"), os.Getenv("COGNITO_APP_CLIENT_ID"), nil)

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	svc := cip.NewFromConfig(cfg)

	// initiate auth
	resp, err := svc.InitiateAuth(context.Background(), &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		return err
	}

	// respond to password verifier challenge
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			return err
		}

		fmt.Println(*resp)
		// 以下はエラーになったので、一旦コメントアウト
		// fmt.Printf("Access Token: %s\n", *resp.AuthenticationResult.AccessToken)
		// fmt.Printf("ID Token: %s\n", *resp.AuthenticationResult.IdToken)
		// fmt.Printf("Refresh Token: %s\n", *resp.AuthenticationResult.RefreshToken)
	}

	return nil

}
