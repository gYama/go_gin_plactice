package database

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProductData struct {
	gorm.Model
	Title string
	Url   string
	Memo  string
}

// 初期化
func ProductInit() {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Init）")
	}
	db.AutoMigrate(&ProductData{})
}

// 追加
func ProductInsert(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Insert)")
	}
	title := ctx.PostForm("title")
	url := ctx.PostForm("url")
	memo := ctx.PostForm("memo")
	db.Create(&ProductData{Title: title, Url: url, Memo: memo})
}

// 更新
func ProductUpdate(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Update)")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var product ProductData
	db.First(&product, id)

	product.Title = ctx.PostForm("title")
	product.Url = ctx.PostForm("url")
	product.Memo = ctx.PostForm("memo")
	db.Save(&product)
}

// 削除
func ProductDelete(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Delete)")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var product ProductData
	db.First(&product, id)
	db.Delete(&product)
}

// 全件取得
func ProductGetAll() []ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	var product []ProductData
	db.Order("created_at desc").Find(&product)
	return product
}

// レコード数取得
func ProductGetRecordCount() int64 {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetRecordCount())")
	}

	var count int64
	var product []ProductData

	db.Find(&product).Count(&count)
	fmt.Println(count)
	return count
}

// 1件取得
func ProductGetOne(ctx *gin.Context) ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetOne())")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var product ProductData
	db.First(&product, id)
	return product
}

func ProductSearch(ctx *gin.Context) []ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	title := ctx.PostForm("andor")
	url := ctx.PostForm("andor")
	memo := ctx.PostForm("andor")
	andor := ctx.PostForm("andor")

	// 各クエリーの生成用
	titleQuery := MakeQuery("title", title, db)
	urlQuery := MakeQuery("url", url, db)
	memoQuery := MakeQuery("memo", memo, db)

	var products []ProductData

	// 項目間のAND/OR
	if andor == "or" {
		// クエリーを出力する場合は、db.Debug().Where()とかやる
		db.Where(
			titleQuery,
		).Or(
			urlQuery,
		).Or(
			memoQuery,
		).Find(&products)
	} else {
		db.Where(
			titleQuery,
		).Where(
			urlQuery,
		).Where(
			memoQuery,
		).Find(&products)
	}

	return products
}
