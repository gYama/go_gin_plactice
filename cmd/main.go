package main

import (
	"context"
	"fmt"
	"go_gin_practice/database"
	"net/http"
	"os"
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
	// store, _ := redis.NewStore(10, "tcp", "0.0.0.0:6379", "", []byte("secret"))

	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/**/*")

	database.Init()

	// 認証が必要なアクショングループを作成
	authGroup := router.Group("/")

	// ここの処理は、全て認証チェックが行われる
	// 認証チェックでエラーの場合は、ログイン画面が表示される
	authGroup.Use(sessionCheckMiddleware())
	{
		authGroup.GET("/", func(ctx *gin.Context) {
			count := database.ProductGetRecordCount()
			ctx.HTML(200, "dashboard.html", gin.H{
				"count": count,
			})
		})

	}

	// 認証が必要なアクショングループを作成（Productでまとめる）
	authProductGroup := router.Group("/")

	// ここの処理は、全て認証チェックが行われる
	// 認証チェックでエラーの場合は、ログイン画面が表示される
	authProductGroup.Use(sessionCheckMiddleware())
	{
		authProductGroup.GET("/product/regist", func(ctx *gin.Context) {
			products := database.ProductGetAll()
			ctx.HTML(200, "product_regist.html", gin.H{
				"products": products,
			})
		})

		authProductGroup.POST("/product/search", func(ctx *gin.Context) {
			products := database.ProductSearch(ctx)

			// 検索条件をセッションに保存
			// ちょっとベタだけど、とりあえず版
			session := sessions.Default(ctx)
			session.Set("product.title", ctx.PostForm("title"))
			session.Set("product.url", ctx.PostForm("url"))
			session.Set("product.memo", ctx.PostForm("memo"))
			session.Set("product.andor", ctx.PostForm("andor"))
			session.Save()

			ctx.HTML(200, "product_search.html", gin.H{
				"products": products,
				"title":    ctx.PostForm("title"),
				"url":      ctx.PostForm("url"),
				"memo":     ctx.PostForm("memo"),
				"andor":    ctx.PostForm("andor"),
				"count":    len(products),
			})
		})

		authProductGroup.GET("/product/search", func(ctx *gin.Context) {
			// セッションから検索条件を取り出し
			session := sessions.Default(ctx)

			// 再検索して表示する
			products := database.ProductSearch(ctx)

			title := ""
			url := ""
			memo := ""
			andor := ""

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
			if session.Get("andor") != nil {
				andor = session.Get("andor").(string)
			}

			ctx.HTML(200, "product_search.html", gin.H{
				"products": products,
				"title":    title,
				"url":      url,
				"memo":     memo,
				"andor":    andor,
				"count":    len(products),
			})
		})

		// Create
		authProductGroup.POST("/product/regist", func(ctx *gin.Context) {
			database.ProductInsert(ctx)
			ctx.Redirect(302, "/product/regist")
		})

		// Detail
		authProductGroup.GET("/product/detail/:id", func(ctx *gin.Context) {
			product := database.ProductGetOne(ctx)
			ctx.HTML(200, "product_detail.html", gin.H{"product": product})
		})

		// Update
		authProductGroup.POST("/product/update/:id", func(ctx *gin.Context) {
			database.ProductUpdate(ctx)
			ctx.Redirect(302, "/product/regist")
		})

		// 削除確認
		authProductGroup.GET("/product/delete_check/:id", func(ctx *gin.Context) {
			product := database.ProductGetOne(ctx)
			ctx.HTML(200, "product_delete.html", gin.H{"product": product})
		})

		// Delete
		authProductGroup.POST("/product/delete/:id", func(ctx *gin.Context) {
			database.ProductDelete(ctx)
			ctx.Redirect(302, "/product/regist")
		})
	}

	// 認証が必要なアクショングループを作成（Customerでまとめる）
	authCustomerGroup := router.Group("/")

	// ここの処理は、全て認証チェックが行われる
	// 認証チェックでエラーの場合は、ログイン画面が表示される
	authCustomerGroup.Use(sessionCheckMiddleware())
	{
		authCustomerGroup.GET("/customer/regist", func(ctx *gin.Context) {
			customers := database.CustomerGetAll()
			ctx.HTML(200, "customer_regist.html", gin.H{
				"customers": customers,
			})
		})

		authCustomerGroup.POST("/customer/search", func(ctx *gin.Context) {
			customers := database.CustomerSearch(ctx)

			// 検索条件をセッションに保存
			// ちょっとベタだけど、とりあえず版
			session := sessions.Default(ctx)
			session.Set("customer.first_name", ctx.PostForm("first_name"))
			session.Set("customer.second_name", ctx.PostForm("second_name"))
			session.Set("customer.phone", ctx.PostForm("phone"))
			session.Set("customer.mail_address", ctx.PostForm("mail_address"))
			session.Set("customer.zipcode", ctx.PostForm("zipcode"))
			session.Set("customer.address", ctx.PostForm("address"))
			session.Set("customer.memo", ctx.PostForm("memo"))
			session.Set("customer.andor", ctx.PostForm("andor"))
			session.Save()

			ctx.HTML(200, "customer_search.html", gin.H{
				"customers":    customers,
				"first_name":   ctx.PostForm("first_name"),
				"second_name":  ctx.PostForm("second_name"),
				"phone":        ctx.PostForm("phone"),
				"mail_address": ctx.PostForm("mail_address"),
				"zipcode":      ctx.PostForm("zipcode"),
				"address":      ctx.PostForm("address"),
				"memo":         ctx.PostForm("memo"),
				"andor":        ctx.PostForm("andor"),
				"count":        len(customers),
			})
		})

		authCustomerGroup.GET("/customer/search", func(ctx *gin.Context) {
			// セッションから検索条件を取り出し
			session := sessions.Default(ctx)

			// 再検索して表示する
			customers := database.CustomerSearch(ctx)

			first_name := ""
			second_name := ""
			phone := ""
			mail_address := ""
			zipcode := ""
			address := ""
			memo := ""
			andor := ""

			// Interface型で返されるので、stringで型変換してあげる
			if session.Get("first_name") != nil {
				first_name = session.Get("first_name").(string)
			}
			if session.Get("second_name") != nil {
				second_name = session.Get("second_name").(string)
			}
			if session.Get("phone") != nil {
				phone = session.Get("phone").(string)
			}
			if session.Get("mail_address") != nil {
				mail_address = session.Get("mail_address").(string)
			}
			if session.Get("zipcode") != nil {
				zipcode = session.Get("zipcode").(string)
			}
			if session.Get("address") != nil {
				address = session.Get("address").(string)
			}
			if session.Get("memo") != nil {
				memo = session.Get("memo").(string)
			}
			if session.Get("andor") != nil {
				andor = session.Get("andor").(string)
			}

			ctx.HTML(200, "customer_search.html", gin.H{
				"customers":    customers,
				"first_name":   first_name,
				"second_name":  second_name,
				"phone":        phone,
				"mail_address": mail_address,
				"zipcode":      zipcode,
				"address":      address,
				"memo":         memo,
				"andor":        andor,
				"count":        len(customers),
			})
		})

		// Create
		authCustomerGroup.POST("/customer/regist", func(ctx *gin.Context) {
			database.CustomerInsert(ctx)
			ctx.Redirect(302, "/customer/regist")
		})

		// Detail
		authCustomerGroup.GET("/customer/detail/:id", func(ctx *gin.Context) {
			customer := database.CustomerGetOne(ctx)
			ctx.HTML(200, "customer_detail.html", gin.H{"customer": customer})
		})

		// Update
		authCustomerGroup.POST("/customer/update/:id", func(ctx *gin.Context) {
			database.CustomerUpdate(ctx)
			ctx.Redirect(302, "/customer/regist")
		})

		// 削除確認
		authCustomerGroup.GET("/customer/delete_check/:id", func(ctx *gin.Context) {
			customer := database.CustomerGetOne(ctx)
			ctx.HTML(200, "customer_delete.html", gin.H{"customer": customer})
		})

		// Delete
		authCustomerGroup.POST("/customer/delete/:id", func(ctx *gin.Context) {
			database.CustomerDelete(ctx)
			ctx.Redirect(302, "/customer/regist")
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
	// fmt.Println(csrp)

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
