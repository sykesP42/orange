package database

import (
	"database/sql"
	"fmt"
	"log"
	"server-go/internal/config"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err = initDefaultData(); err != nil {
		return fmt.Errorf("failed to init default data: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			real_name TEXT,
			role TEXT DEFAULT 'user',
			avatar TEXT,
			email TEXT,
			phone TEXT,
			bio TEXT,
			status TEXT DEFAULT 'active',
			team_id INTEGER,
			team_name TEXT,
			team_color TEXT,
			approved_by TEXT,
			approved_time TEXT,
			is_deleted INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS blogger (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT NOT NULL,
			category TEXT,
			product TEXT,
			contact TEXT,
			wechat TEXT,
			custom_contact TEXT,
			platform TEXT,
			platform_account TEXT,
			description TEXT,
			tags TEXT,
			avatar TEXT,
			user_name TEXT,
			real_name TEXT,
			status TEXT DEFAULT '初次联系',
			status_update_time TEXT,
			status_remark TEXT,
			is_deleted INTEGER DEFAULT 0,
			is_invalid INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS blogger_status_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER,
			old_status TEXT,
			new_status TEXT,
			operator TEXT,
			remark TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			color TEXT DEFAULT '#6b7280',
			icon TEXT DEFAULT '🏷️',
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS platforms (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			icon TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS operation_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT,
			target TEXT,
			operator TEXT,
			detail TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			logo TEXT,
			bg_image TEXT,
			color TEXT,
			leader_id INTEGER,
			leader_name TEXT,
			admin_ids TEXT,
			announcement TEXT,
			is_public INTEGER DEFAULT 1,
			need_approve INTEGER DEFAULT 0,
			member_count INTEGER DEFAULT 1,
			invite_code TEXT UNIQUE,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS team_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_id INTEGER,
			author_id INTEGER,
			author_name TEXT,
			author_username TEXT,
			author_avatar TEXT,
			title TEXT,
			content TEXT,
			category TEXT DEFAULT '综合',
			comment_count INTEGER DEFAULT 0,
			like_count INTEGER DEFAULT 0,
			view_count INTEGER DEFAULT 0,
			is_pinned INTEGER DEFAULT 0,
			is_featured INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS team_post_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			author_id INTEGER,
			author_name TEXT,
			author_username TEXT,
			author_avatar TEXT,
			content TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS team_post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS team_post_collects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS public_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER,
			title TEXT,
			content TEXT,
			category TEXT DEFAULT '讨论',
			is_pinned INTEGER DEFAULT 0,
			is_featured INTEGER DEFAULT 0,
			view_count INTEGER DEFAULT 0,
			like_count INTEGER DEFAULT 0,
			comment_count INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS public_post_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			author_id INTEGER,
			author_name TEXT,
			author_avatar TEXT,
			content TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS public_post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS public_post_collects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			type TEXT,
			title TEXT,
			content TEXT,
			priority TEXT DEFAULT 'normal',
			team_id INTEGER,
			from_user TEXT,
			blogger_id INTEGER,
			post_id INTEGER,
			is_read INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS private_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER,
			from_user_name TEXT,
			to_user_id INTEGER,
			content TEXT,
			is_read INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS team_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_id INTEGER,
			sender_id INTEGER,
			sender_name TEXT,
			sender_avatar TEXT,
			title TEXT,
			content TEXT,
			type TEXT DEFAULT 'normal',
			is_read INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS followup (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER,
			content TEXT,
			type TEXT DEFAULT '跟进',
			operator TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cooperation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER,
			title TEXT,
			date TEXT,
			status TEXT,
			product TEXT,
			amount REAL,
			note TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS blogger_transfer_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER,
			blogger_name TEXT,
			from_user_id INTEGER,
			from_user_name TEXT,
			to_user_id INTEGER,
			to_user_name TEXT,
			reason TEXT,
			status TEXT DEFAULT 'pending',
			from_confirmed INTEGER DEFAULT 1,
			to_confirmed INTEGER DEFAULT 0,
			admin_confirmed INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS system_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			value TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS announcements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			filename TEXT,
			size INTEGER,
			created DATETIME,
			mtime DATETIME,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func initDefaultData() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '管理员', 'admin', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "admin", string(hashedPassword))
		if err != nil {
			return err
		}

		hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '张三', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user1", string(hashedPassword1))

		hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '李四', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user2", string(hashedPassword2))

		hashedPassword3, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '王五', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user3", string(hashedPassword3))

		log.Println("Default users created (admin/admin123, user1/admin123, user2/admin123, user3/admin123)")
	}

	initCategories()
	initPlatforms()
	initProducts()

	return nil
}

func initCategories() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if count > 0 {
		return
	}

	categories := []struct {
		name  string
		color string
		icon  string
	}{
		{"美妆", "#f43f5e", "💄"},
		{"穿搭", "#a855f7", "👗"},
		{"美食", "#f97316", "🍔"},
		{"旅行", "#10b981", "✈️"},
		{"数码", "#3b82f6", "📱"},
		{"母婴", "#ec4899", "🍼"},
		{"家居", "#14b8a6", "🏠"},
		{"健身", "#ef4444", "💪"},
		{"情感", "#e11d48", "💕"},
		{"游戏", "#8b5cf6", "🎮"},
		{"影视", "#0ea5e9", "🎬"},
		{"音乐", "#a855f7", "🎵"},
		{"教育", "#22c55e", "📚"},
		{"职场", "#64748b", "💼"},
		{"科技", "#06b6d4", "🚀"},
		{"财经", "#eab308", "💰"},
		{"萌宠", "#d946ef", "🐱"},
		{"露营", "#059669", "🏕️"},
		{"户外", "#6366f1", "⛰️"},
		{"其他", "#6b7280", "🏷️"},
	}

	for _, cat := range categories {
		DB.Exec("INSERT INTO categories (name, color, icon) VALUES (?, ?, ?)", cat.name, cat.color, cat.icon)
	}
	log.Printf("Initialized %d categories", len(categories))
}

func initPlatforms() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM platforms").Scan(&count)
	if count > 0 {
		return
	}

	platforms := []string{
		"小红书", "抖音", "微博", "B站", "快手",
		"微信公众号", "视频号", "知乎", "淘宝直播", "其他",
	}

	for _, p := range platforms {
		DB.Exec("INSERT INTO platforms (name) VALUES (?)", p)
	}
	log.Printf("Initialized %d platforms", len(platforms))
}

func initProducts() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		return
	}

	products := []string{
		"品牌合作", "自营产品", "带货分成", "广告投放", "直播带货",
		"短视频", "图文种草", "知识付费", "会员付费", "活动赞助",
		"周边销售", "其他",
	}

	for _, p := range products {
		DB.Exec("INSERT INTO products (name) VALUES (?)", p)
	}
	log.Printf("Initialized %d products", len(products))
}

func AddLog(action, target, operator, detail string) {
	DB.Exec(`INSERT INTO operation_log (action, target, operator, detail) VALUES (?, ?, ?, ?)`,
		action, target, operator, detail)
}
