package yiigo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DBDriver 数据库驱动
type DBDriver string

const (
	MySQL    DBDriver = "mysql"
	Postgres DBDriver = "pgx"
	SQLite   DBDriver = "sqlite3"
)

var dbMap = make(map[string]*sqlx.DB)

// DBConfig 数据库初始化配置
type DBConfig struct {
	// DSN 数据源名称
	// [-- MySQL] username:password@tcp(localhost:3306)/dbname?timeout=10s&charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
	// [Postgres] host=localhost port=5432 user=root password=secret dbname=test connect_timeout=10 sslmode=disable
	// [- SQLite] file::memory:?cache=shared
	DSN string `json:"dsn"`
	// Options 配置选项
	Options *DBOptions `json:"options"`
}

// DBOptions 数据库配置选项
type DBOptions struct {
	// MaxOpenConns 设置最大可打开的连接数
	MaxOpenConns int `json:"max_open_conns"`
	// MaxIdleConns 连接池最大闲置连接数
	MaxIdleConns int `json:"max_idle_conns"`
	// ConnMaxLifetime 连接的最大生命时长
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	// ConnMaxIdleTime 连接最大闲置时间
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

func initDB(name string, driver DBDriver, cfg *DBConfig) error {
	db, err := sql.Open(string(driver), cfg.DSN)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return err
	}

	if cfg.Options != nil {
		db.SetMaxOpenConns(cfg.Options.MaxOpenConns)
		db.SetMaxIdleConns(cfg.Options.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.Options.ConnMaxLifetime)
		db.SetConnMaxIdleTime(cfg.Options.ConnMaxIdleTime)
	}

	dbMap[name] = sqlx.NewDb(db, string(driver))

	return nil
}

// DB 返回一个sqlx数据库实例
func DB(name ...string) (*sqlx.DB, error) {
	key := Default
	if len(name) != 0 {
		key = name[0]
	}

	db, ok := dbMap[key]
	if !ok {
		return nil, fmt.Errorf("unknown db.%s (forgotten configure?)", key)
	}

	return db, nil
}

// MustDB 返回一个sqlx数据库实例，如果不存在，则Panic
func MustDB(name ...string) *sqlx.DB {
	db, err := DB(name...)
	if err != nil {
		logger.Panic(err.Error())
	}

	return db
}

// CloseDB 关闭数据库连接，如果未指定名称，则关闭全部
func CloseDB(name ...string) {
	if len(name) == 0 {
		for key, db := range dbMap {
			if err := db.Close(); err != nil {
				logger.Error(fmt.Sprintf("db.%s close error", key), zap.Error(err))
			}
		}

		return
	}

	for _, key := range name {
		if db, ok := dbMap[key]; ok {
			if err := db.Close(); err != nil {
				logger.Error(fmt.Sprintf("db.%s close error", key), zap.Error(err))
			}
		}
	}
}

// Transaction 执行数据库事物
func Transaction(ctx context.Context, db *sqlx.DB, f func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err = f(ctx, tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
