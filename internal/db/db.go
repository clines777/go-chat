package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gochat/internal/config"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL     *sql.DB
	Builder sq.StatementBuilderType
}

// Init DB
func Init(cfg *config.EnvConfig) (*DB, error) {

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

	return &DB{
		SQL:     sqlDB,
		Builder: builder,
	}, nil
}

func (d *DB) Ping(ctx context.Context) error {

	if err := d.SQL.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	var one int
	if err := d.SQL.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("select 1: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("unexpected ping result: %d", one)
	}
	return nil
}

func (d *DB) Close() error {
	return d.SQL.Close()
}
