package database

import (
	"fmt"
	"regexp"
	"strings"

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
func Init() {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Init）")
	}
	db.AutoMigrate(&ProductData{})
}

// 追加
func Insert(title string, url string, memo string) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Insert)")
	}
	db.Create(&ProductData{Title: title, Url: url, Memo: memo})
}

// 更新
func Update(id int, title string, url string, memo string) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Update)")
	}
	var product ProductData
	db.First(&product, id)
	product.Title = title
	product.Url = url
	product.Memo = memo
	db.Save(&product)
}

// 削除
func Delete(id int) {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error（Delete)")
	}
	var product ProductData
	db.First(&product, id)
	db.Delete(&product)
}

// 全件取得
func GetAll() []ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	var product []ProductData
	db.Debug().Order("created_at desc").Find(&product)
	return product
}

// レコード数取得
func GetRecordCount() int64 {
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
func GetOne(id int) ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetOne())")
	}
	var product ProductData
	db.First(&product, id)
	return product
}

func Search(title string, url string, memo string, andor string) []ProductData {
	db, err := gorm.Open(sqlite.Open("test.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("db open error(GetAll())")
	}

	// 各クエリーの生成用
	titleQuery := db
	urlQuery := db
	memoQuery := db

	// 複数入力された場合は、スペース（半角/全角）区切りで分割
	reg := "[ 　]"

	// titleの検索クエリー
	titles := regexp.MustCompile(reg).Split(title, -1)

	for i := 0; i < len(titles); i++ {
		str := strings.TrimSpace(titles[i])
		if len(str) == 0 {
			continue
		}
		if i == 0 {
			titleQuery = titleQuery.Where("title like ?", "%"+str+"%")
		} else {
			titleQuery = titleQuery.Or("title like ?", "%"+str+"%")
		}
	}

	// urlの検索クエリー
	urls := regexp.MustCompile(reg).Split(url, -1)

	for i := 0; i < len(urls); i++ {
		str := strings.TrimSpace(urls[i])
		if len(str) == 0 {
			continue
		}
		if i == 0 {
			urlQuery = urlQuery.Where("title like ?", "%"+str+"%")
		} else {
			urlQuery = urlQuery.Or("title like ?", "%"+str+"%")
		}
	}

	// memoの検索クエリー
	memos := regexp.MustCompile(reg).Split(memo, -1)

	for i := 0; i < len(memos); i++ {
		str := strings.TrimSpace(memos[i])
		if len(str) == 0 {
			continue
		}
		if i == 0 {
			memoQuery = memoQuery.Where("title like ?", "%"+str+"%")
		} else {
			memoQuery = memoQuery.Or("title like ?", "%"+str+"%")
		}
	}

	var products []ProductData

	// 項目間のAND/OR
	if andor == "or" {
		// クエリーを出力する場合は、db.Debug()を使う
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
