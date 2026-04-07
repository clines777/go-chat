package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gochat/internal/infra"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var dbConn *DB

type DB struct {
	SQL     *sql.DB
	Builder sq.StatementBuilderType
	ctx     context.Context
}

// GetDBConn DB
func GetDBConn() (*DB, error) {

	if dbConn == nil {
		cfg := infra.GetEnvConfig()
		// 本地開發環境檢查並載入env參數
		envPath := "../.env"
		_, envFileErr := os.Stat(envPath)
		if !os.IsNotExist(envFileErr) {
			_ = godotenv.Load()
		}

		if cfg.DBSSLMode == "" {
			cfg.DBSSLMode = "disable"
		}
		if cfg.DBTimeZone == "" {
			cfg.DBTimeZone = "Asia/Taipei"
		}

		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBTimeZone,
		)

		sqlDB, dbErr := sql.Open("pgx", dsn)
		if dbErr != nil {
			return nil, fmt.Errorf("sql.Open: %w", dbErr)
		}

		if cfg.DBMaxConn > 0 {
			sqlDB.SetMaxOpenConns(cfg.DBMaxConn)
		} else {
			sqlDB.SetMaxOpenConns(20)
		}
		if cfg.DBMaxIdleConn > 0 {
			sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConn)
		} else {
			sqlDB.SetMaxIdleConns(10)
		}
		if cfg.DBConnMaxLifeTime > 0 {
			sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifeTime)
		} else {
			sqlDB.SetConnMaxLifetime(30 * time.Minute)
		}
		if cfg.DBConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)
		} else {
			sqlDB.SetConnMaxIdleTime(5 * time.Minute)
		}

		builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

		dbConn = &DB{
			SQL:     sqlDB,
			Builder: builder,
			ctx:     context.Background(),
		}
	}

	return dbConn, nil
}

func (d *DB) Ping() error {

	if err := d.SQL.PingContext(d.ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	var one int8
	if err := d.SQL.QueryRowContext(d.ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("select 1: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("unexpected ping result: %d", one)
	}
	return nil
}
