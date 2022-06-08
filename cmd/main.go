package main

import (
	"context"
	"fmt"
	"go_gin_practice/database"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	// store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/*.html")

	database.Init()

	// 認証済みユーザー用のグループを作成
	authUserGroup := router.Group("/")

	// ここの処理は、全て認証チェックが行われる
	// 認証チェックでエラーの場合は、ログイン画面が表示される
	authUserGroup.Use(sessionCheckMiddleware())
	{
		authUserGroup.GET("/", func(ctx *gin.Context) {
			ctx.HTML(200, "dashboard.html", gin.H{})
		})

		authUserGroup.GET("/list", func(ctx *gin.Context) {
			products := database.GetAll()
			ctx.HTML(200, "index.html", gin.H{
				"products": products,
			})
		})

		authUserGroup.POST("/search", func(ctx *gin.Context) {
			title := ctx.PostForm("title")
			url := ctx.PostForm("url")
			memo := ctx.PostForm("memo")
			products := database.Search(title, url, memo)

			// 検索条件をセッションに保存
			// ちょっとベタだけど、とりあえず版
			session := sessions.Default(ctx)
			session.Set("title", title)
			session.Set("url", url)
			session.Set("memo", memo)
			session.Save()

			ctx.HTML(200, "search.html", gin.H{
				"products": products,
				"title":    title,
				"url":      url,
				"memo":     memo,
			})
		})

		authUserGroup.GET("/search", func(ctx *gin.Context) {
			// セッションから検索条件を取り出し
			session := sessions.Default(ctx)

			title := ""
			url := ""
			memo := ""

			// Interface型で返されるので、stringで型変換してあげる
			if session.Get("title") != nil {
				title = session.Get("title").(string)
			}
			if session.Get("url") != nil {
				url = session.Get("url").(string)
			}
			if session.Get("memo") != nil {
				memo = session.Get("memo").(string)
			}

			// 再検索して表示する
			products := database.Search(title, url, memo)

			ctx.HTML(200, "search.html", gin.H{
				"products": products,
				"title":    title,
				"url":      url,
				"memo":     memo,
			})
		})

		// Create
		authUserGroup.POST("/new", func(ctx *gin.Context) {
			title := ctx.PostForm("title")
			url := ctx.PostForm("url")
			memo := ctx.PostForm("memo")
			database.Insert(title, url, memo)
			ctx.Redirect(302, "/list")
		})

		// Detail
		authUserGroup.GET("/detail/:id", func(ctx *gin.Context) {
			n := ctx.Param("id")
			id, err := strconv.Atoi(n)
			if err != nil {
				panic(err)
			}
			product := database.GetOne(id)
			ctx.HTML(200, "detail.html", gin.H{"product": product})
		})

		// Update
		authUserGroup.POST("/update/:id", func(ctx *gin.Context) {
			n := ctx.Param("id")
			id, err := strconv.Atoi(n)
			if err != nil {
				panic("ERROR")
			}
			title := ctx.PostForm("title")
			url := ctx.PostForm("url")
			memo := ctx.PostForm("memo")
			database.Update(id, title, url, memo)
			ctx.Redirect(302, "/list")
		})

		// 削除確認
		authUserGroup.GET("/delete_check/:id", func(ctx *gin.Context) {
			n := ctx.Param("id")
			id, err := strconv.Atoi(n)
			if err != nil {
				panic("ERROR")
			}
			product := database.GetOne(id)
			ctx.HTML(200, "delete.html", gin.H{"product": product})
		})

		// Delete
		authUserGroup.POST("/delete/:id", func(ctx *gin.Context) {
			n := ctx.Param("id")
			id, err := strconv.Atoi(n)
			if err != nil {
				panic("ERROR")
			}
			database.Delete(id)
			ctx.Redirect(302, "/list")
		})
	}

	// Login
	router.POST("/login", func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		password := ctx.PostForm("password")

		err := login(id, password, ctx)
		if err != nil {
			fmt.Println(err)
			ctx.HTML(200, "login.html", gin.H{"message": "ログイン失敗"})
		}
		session := sessions.Default(ctx)
		session.Set("isAuthenticated", true)
		session.Save()

		ctx.Redirect(302, "/")
	})

	// Login
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	})

	// Logout
	router.GET("/logout", func(ctx *gin.Context) {
		logout(ctx)
		ctx.Redirect(302, "/login")
	})

	router.Run()
}

// Login（POST）
func login(id string, password string, ctx *gin.Context) error {

	// cognitosrpはあまりメンテナンスされてなさそう
	// いずれ自分で実装するが、一旦このまま利用する
	csrp, err := cognitosrp.NewCognitoSRP(
		id,
		password,
		os.Getenv("COGNITO_USERPOOL_ID"),
		os.Getenv("COGNITO_APP_CLIENT_ID"),
		aws.String(os.Getenv("COGNITO_APP_CLIENT_SECRET")),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(csrp)

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

// Logout
func logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}

// 認証チェック　エラーの場合は、ログイン画面を表示する
func sessionCheckMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		isAuthenticated := session.Get("isAuthenticated")

		if isAuthenticated == nil {
			ctx.Redirect(http.StatusMovedPermanently, "/login")
			ctx.Abort()
		} else {
			ctx.Set("isAuthenticated", isAuthenticated)
			ctx.Next()
		}
	}
}
