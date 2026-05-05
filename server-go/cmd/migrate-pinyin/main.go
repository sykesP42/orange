package main

import (
	"database/sql"
	"fmt"
	"log"
	"server-go/internal/utils"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "./data/orange.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("打开数据库失败:", err)
	}
	defer db.Close()

	fmt.Println("开始为现有博主数据填充拼音字段...")

	rows, err := db.Query(`SELECT id, nickname FROM blogger WHERE (nickname_pinyin IS NULL OR nickname_pinyin = '' OR nickname_pinyin IS NOT NULL) AND is_deleted = 0`)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var nickname string
		err := rows.Scan(&id, &nickname)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		pinyin := utils.ConvertToPinyin(nickname)
		_, err = db.Exec(`UPDATE blogger SET nickname_pinyin = ? WHERE id = ?`, pinyin, id)
		if err != nil {
			log.Printf("更新ID=%d失败: %v", id, err)
			continue
		}

		count++
		if count%100 == 0 {
			fmt.Printf("已处理 %d 条记录...\n", count)
		}
	}

	fmt.Printf("✅ 迁移完成！共更新 %d 条博主记录的拼音字段\n", count)
}
