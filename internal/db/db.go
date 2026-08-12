package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

// DB 全局数据库连接
var DB *sql.DB

type sqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// InitDB 初始化数据库。新文件写入当前结构和默认数据；已有文件只补齐
// 当前结构，不覆盖用户直接修改的参考数据，也不升级历史数据库。
func InitDB(dbPath string) (*sql.DB, error) {
	databaseExists := true
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		databaseExists = false
	} else if err != nil {
		return nil, err
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 建议单连接
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := initializeCurrentDatabase(db, !databaseExists); err != nil {
		db.Close()
		if !databaseExists {
			_ = os.Remove(dbPath)
		}
		return nil, err
	}
	// 保存全局连接
	DB = db

	return db, nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		DB = nil
		return err
	}
	return nil
}

// GetDB 获取数据库连接
func GetDB() *sql.DB {
	return DB
}

// initializeCurrentDatabase 在一个事务中建立当前结构和索引，并仅在首次
// 创建数据库时写入参考数据。
func initializeCurrentDatabase(db *sql.DB, seedDefaults bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createTables(tx); err != nil {
		return err
	}
	if err := createIndexes(tx); err != nil {
		return err
	}
	if seedDefaults {
		if err := seedData(tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func createIndexes(db sqlExecutor) error {
	statements := []string{
		"CREATE INDEX IF NOT EXISTS idx_saves_created_at ON saves(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_items_region ON items(region)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_minigames_name ON minigames(name)",
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

// createTables 创建所有数据库表。
func createTables(db sqlExecutor) error {
	// 物资表（国内+国外）
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			price_min INTEGER NOT NULL,
			price_max INTEGER NOT NULL,
			effects TEXT NOT NULL,
			region TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 公司表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			risk INTEGER NOT NULL,
			profit INTEGER NOT NULL,
			time INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 古董表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS antiques (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			material INTEGER NOT NULL,
			img TEXT NOT NULL,
			desc TEXT NOT NULL,
			level INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 股票表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stocks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			risk INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 股票新闻表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stock_news (
			id INTEGER PRIMARY KEY,
			content TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 房产表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS houses (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			price_max INTEGER NOT NULL,
			price_min INTEGER NOT NULL,
			health INTEGER NOT NULL,
			fame INTEGER NOT NULL,
			img TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 车辆表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cars (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			price_max INTEGER NOT NULL,
			price_min INTEGER NOT NULL,
			health INTEGER NOT NULL,
			fame INTEGER NOT NULL,
			img TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 约会对象信息表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS datings (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			image TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender INTEGER NOT NULL CHECK (gender IN (0, 1)),
			nationality TEXT NOT NULL,
			occupation TEXT NOT NULL,
			description TEXT NOT NULL,
			cost INTEGER NOT NULL,
			meet_conditions TEXT NOT NULL,
			gifts TEXT NOT NULL,
			locations TEXT NOT NULL,
			meet_scene TEXT NOT NULL DEFAULT ''
		)
	`)
	if err != nil {
		return err
	}

	// 银行任务表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bank_tasks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			type TEXT NOT NULL,
			target_value INTEGER NOT NULL,
			reward INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 疾病数据表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS diseases (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			symptoms TEXT NOT NULL,
			healthimpact INTEGER NOT NULL,
			upgradedays INTEGER NOT NULL,
			treatment TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 药品数据表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS treats (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			price INTEGER NOT NULL,
			heal INTEGER NOT NULL,
			sideeffect INTEGER NOT NULL,
			source TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 医院数据表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS hospitals (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			icon TEXT NOT NULL,
			description TEXT NOT NULL,
			services TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 小游戏数据表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS minigames (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			cnname TEXT NOT NULL,
			icon TEXT NOT NULL,
			description TEXT NOT NULL,
			type TEXT NOT NULL,
			difficulty TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 存档元数据表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS saves (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			game_year INTEGER NOT NULL,
			save_version INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// 游戏状态存档表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS game_saves (
			save_id INTEGER PRIMARY KEY,
			game_data TEXT NOT NULL,
			FOREIGN KEY (save_id) REFERENCES saves(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// 用户状态存档表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_saves (
			save_id INTEGER PRIMARY KEY,
			user_data TEXT NOT NULL,
			FOREIGN KEY (save_id) REFERENCES saves(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// 公告存档表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS announce_saves (
			save_id INTEGER PRIMARY KEY,
			announce_data TEXT NOT NULL,
			FOREIGN KEY (save_id) REFERENCES saves(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
