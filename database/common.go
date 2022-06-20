package database

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

func Init() {
	ProductInit()
	CustomerInit()
}

func MakeQuery(column string, input string, db *gorm.DB) *gorm.DB {
	// 複数入力された場合は、スペース（半角/全角）区切りで分割
	reg := "[ 　]"

	// titleをスペースで分割
	words := regexp.MustCompile(reg).Split(input, -1)

	// titleの検索クエリー生成
	for i := 0; i < len(words); i++ {
		str := strings.TrimSpace(words[i])
		if len(str) == 0 {
			continue
		}

		likeStr := fmt.Sprintf("%s like ?", column)

		if i == 0 {
			db = db.Where(likeStr, "%"+str+"%")
		} else {
			db = db.Or(likeStr, "%"+str+"%")
		}
	}

	return db

}
