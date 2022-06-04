package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type ProductData struct {
	gorm.Model
	Title string
	Url   string
	Memo  string
}

// 初期化
func Init() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error（Init）")
	}
	db.AutoMigrate(&ProductData{})
	defer db.Close()
}

// 追加
func Insert(title string, url string, memo string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error（Insert)")
	}
	db.Create(&ProductData{Title: title, Url: url, Memo: memo})
	defer db.Close()
}

// 更新
func Update(id int, title string, url string, memo string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error（Update)")
	}
	var product ProductData
	db.First(&product, id)
	product.Title = title
	product.Url = url
	product.Memo = memo
	db.Save(&product)
	db.Close()
}

// 削除
func Delete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error（Delete)")
	}
	var product ProductData
	db.First(&product, id)
	db.Delete(&product)
	db.Close()
}

// 全件取得
func GetAll() []ProductData {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error(GetAll())")
	}

	var product []ProductData
	db.Order("created_at desc").Find(&product)
	db.Close()
	return product
}

// 1件取得
func GetOne(id int) ProductData {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("db open error(GetOne())")
	}
	var product ProductData
	db.First(&product, id)
	db.Close()
	return product
}
