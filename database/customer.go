package database

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CustomerData struct {
	gorm.Model
	FirstName   string
	SecondName  string
	Phone       string
	MailAddress string
	Zipcode     string
	Address     string
	Memo        string
}

// 初期化
func CustomerInit() {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Init）")
	}
	//
	db.AutoMigrate(&CustomerData{})
}

// 追加
func CustomerInsert(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Insert)")
	}
	db.Create(&CustomerData{
		FirstName:   strings.TrimSpace(ctx.PostForm("first_name")),
		SecondName:  strings.TrimSpace(ctx.PostForm("second_name")),
		Phone:       strings.TrimSpace(ctx.PostForm("phone")),
		MailAddress: strings.TrimSpace(ctx.PostForm("mail_address")),
		Zipcode:     strings.TrimSpace(ctx.PostForm("zipcode")),
		Address:     strings.TrimSpace(ctx.PostForm("address")),
		Memo:        strings.TrimSpace(ctx.PostForm("memo")),
	})
}

// 更新
func CustomerUpdate(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Update)")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var customer CustomerData
	db.First(&customer, id)
	customer.FirstName = strings.TrimSpace(ctx.PostForm("first_name"))
	customer.SecondName = strings.TrimSpace(ctx.PostForm("second_name"))
	customer.Phone = strings.TrimSpace(ctx.PostForm("phone"))
	customer.MailAddress = strings.TrimSpace(ctx.PostForm("mail_address"))
	customer.Zipcode = strings.TrimSpace(ctx.PostForm("zipcode"))
	customer.Address = strings.TrimSpace(ctx.PostForm("address"))
	customer.Memo = strings.TrimSpace(ctx.PostForm("memo"))
	db.Save(&customer)
}

// 削除
func CustomerDelete(ctx *gin.Context) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Delete)")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var customer CustomerData
	db.First(&customer, id)
	db.Delete(&customer)
}

// 全件取得
func CustomerGetAll() []CustomerData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	var customer []CustomerData
	db.Order("created_at desc").Find(&customer)
	return customer
}

// レコード数取得
func CustomerGetRecordCount() int64 {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetRecordCount())")
	}

	var count int64
	var customer []CustomerData

	db.Find(&customer).Count(&count)
	return count
}

// 1件取得
func CustomerGetOne(ctx *gin.Context) CustomerData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetOne())")
	}

	n := ctx.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	var customer CustomerData
	db.First(&customer, id)
	return customer
}

func CustomerSearch(ctx *gin.Context) []CustomerData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	andor := ctx.PostForm("andor")

	// 各クエリーの生成用
	firstNameQuery := MakeQuery("first_name", ctx.PostForm("first_name"), db)
	secondNameQuery := MakeQuery("second_name", ctx.PostForm("second_name"), db)
	phoneQuery := MakeQuery("phone", ctx.PostForm("phone"), db)
	mailAddressQuery := MakeQuery("mail_address", ctx.PostForm("mail_address"), db)
	zipcodeQuery := MakeQuery("zipcode", ctx.PostForm("zipcode"), db)
	addressQuery := MakeQuery("address", ctx.PostForm("address"), db)
	memoQuery := MakeQuery("memo", ctx.PostForm("memo"), db)

	var customer []CustomerData

	// 項目間のAND/OR
	if andor == "or" {
		// クエリーを出力する場合は、db.Debug().Where()とかやる
		db.Where(
			firstNameQuery,
		).Or(
			secondNameQuery,
		).Or(
			phoneQuery,
		).Or(
			mailAddressQuery,
		).Or(
			zipcodeQuery,
		).Or(
			addressQuery,
		).Or(
			memoQuery,
		).Find(&customer)
	} else {
		db.Where(
			firstNameQuery,
		).Where(
			secondNameQuery,
		).Where(
			phoneQuery,
		).Where(
			mailAddressQuery,
		).Where(
			zipcodeQuery,
		).Where(
			addressQuery,
		).Where(
			memoQuery,
		).Find(&customer)
	}

	return customer
}
