package database

import (
	"database/sql"
	"fmt"
	"log"
	"server-go/internal/config"
	"server-go/pkg/auth"

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

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(0)

	if _, err = DB.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return fmt.Errorf("failed to enable WAL mode: %w", err)
	}
	if _, err = DB.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		return fmt.Errorf("failed to set synchronous mode: %w", err)
	}
	if _, err = DB.Exec("PRAGMA cache_size=-64000"); err != nil {
		return fmt.Errorf("failed to set cache size: %w", err)
	}
	if _, err = DB.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return fmt.Errorf("failed to set busy timeout: %w", err)
	}

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
		`CREATE TABLE IF NOT EXISTS blogger_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			color TEXT DEFAULT '#6b7280',
			is_system INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS blogger_tag_relations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(blogger_id, tag_id)
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
		`CREATE TABLE IF NOT EXISTS followup_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			name TEXT NOT NULL,
			content TEXT NOT NULL,
			category TEXT,
			is_default INTEGER DEFAULT 0,
			use_count INTEGER DEFAULT 0,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP
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
			tags TEXT,
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
		`CREATE TABLE IF NOT EXISTS blogger_evaluations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			blogger_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			rating INTEGER DEFAULT 5,
			comment TEXT,
			tags TEXT,
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(blogger_id, user_id)
		)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return err
		}
	}

	migrations := []string{
		`ALTER TABLE team_posts ADD COLUMN tags TEXT`,
	}
	for _, m := range migrations {
		DB.Exec(m)
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`,
		`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`,
		`CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_user_name ON blogger(user_name)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_status ON blogger(status)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_is_deleted ON blogger(is_deleted)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_category ON blogger(category)`,
		`CREATE INDEX IF NOT EXISTS idx_followup_blogger_id ON followup(blogger_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cooperation_blogger_id ON cooperation(blogger_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_posts_team_id ON team_posts(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_posts_author_id ON team_posts(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_post_comments_post_id ON team_post_comments(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_post_likes_post_id ON team_post_likes(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_post_collects_post_id ON team_post_collects(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_public_posts_author_id ON public_posts(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read)`,
		`CREATE INDEX IF NOT EXISTS idx_private_messages_from_user ON private_messages(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_private_messages_to_user ON private_messages(to_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_team_messages_team_id ON team_messages(team_id)`,
		`CREATE INDEX IF NOT EXISTS idx_operation_log_create_time ON operation_log(create_time)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_status_history_blogger_id ON blogger_status_history(blogger_id)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_transfer_requests_status ON blogger_transfer_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_blogger_evaluations_blogger_id ON blogger_evaluations(blogger_id)`,
	}
	for _, idx := range indexes {
		if _, err := DB.Exec(idx); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
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
		hashedPassword, err := auth.HashPassword("admin123")
		if err != nil {
			return err
		}
		_, err = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '管理员', 'admin', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "admin", hashedPassword)
		if err != nil {
			return err
		}

		hashedPassword1, _ := auth.HashPassword("admin123")
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '张三', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user1", hashedPassword1)

		hashedPassword2, _ := auth.HashPassword("admin123")
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '李四', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user2", hashedPassword2)

		hashedPassword3, _ := auth.HashPassword("admin123")
		_, _ = DB.Exec(`
			INSERT INTO users (username, password, real_name, role, status, create_time, update_time)
			VALUES (?, ?, '王五', 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, "user3", hashedPassword3)

		log.Println("Default users created (admin/admin123, user1/admin123, user2/admin123, user3/admin123)")
	}

	initCategories()
	initPlatforms()
	initProducts()
	initBloggerTags()
	initFollowupTemplates()

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

func initBloggerTags() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM blogger_tags").Scan(&count)
	if count > 0 {
		return
	}

	tags := []struct {
		name     string
		color    string
		isSystem int
	}{
		{"优质", "#10b981", 1},
		{"潜力", "#3b82f6", 1},
		{"待跟进", "#f59e0b", 1},
		{"重点", "#ef4444", 1},
		{"黑名单", "#6b7280", 1},
		{"长期合作", "#8b5cf6", 1},
		{"新博主", "#06b6d4", 1},
		{"高转化", "#f97316", 1},
	}

	for _, tag := range tags {
		DB.Exec("INSERT INTO blogger_tags (name, color, is_system) VALUES (?, ?, ?)", tag.name, tag.color, tag.isSystem)
	}
	log.Printf("Initialized %d blogger tags", len(tags))
}

func initFollowupTemplates() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM followup_templates").Scan(&count)
	if count > 0 {
		return
	}

	templates := []struct {
		name     string
		content  string
		category string
	}{
		{"初次联系模板", "您好！我是XXX公司的商务合作负责人，关注您很久了。\n\n您的内容非常有特色，我们希望能够与您探讨合作机会。\n\n方便的话可以添加我的微信：XXXX", "初次联系"},
		{"合作洽谈模板", "您好！\n\n我们品牌正在寻找优质博主合作，您的账号非常符合我们的需求。\n\n请问您有兴趣了解一下合作详情吗？期待您的回复~", "洽谈中"},
		{"合作确认模板", "您好！很高兴您对我们品牌的认可。\n\n关于合作细节，我来跟您详细说明：\n1. 合作形式：\n2. 内容要求：\n3. 报价：\n\n请确认是否可以合作，谢谢！", "合作中"},
		{"感谢回复模板", "非常感谢您的回复！\n\n我们会尽快安排专人与您对接详细合作事宜。\n\n祝您创作愉快！", "通用"},
		{"跟进未回复模板", "您好！之前给您发的合作邀请不知道您是否看到了？\n\n如果您对我们的合作有任何问题或建议，欢迎随时联系我。\n\n期待您的回复~", "跟进"},
	}

	for _, t := range templates {
		DB.Exec("INSERT INTO followup_templates (name, content, category, is_default) VALUES (?, ?, ?, 1)", t.name, t.content, t.category)
	}
	log.Printf("Initialized %d followup templates", len(templates))
}

func AddLog(action, target, operator, detail string) {
	DB.Exec(`INSERT INTO operation_log (action, target, operator, detail) VALUES (?, ?, ?, ?)`,
		action, target, operator, detail)
}
