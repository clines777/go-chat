package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gochat/internal/infra"
)

type DB struct {
	DB      *sqlx.DB
	Builder sq.StatementBuilderType
	ctx     context.Context
}

var dbOnce = sync.OnceValues(buildDBConn)

func GetDBConn() (*DB, error) {
	return dbOnce()
}

func buildDBConn() (*DB, error) {
	cfg := infra.GetEnvConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBTimeZone,
	)

	sqlDB, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
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

	return &DB{
		DB:      sqlDB,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		ctx:     context.Background(),
	}, nil
}

func (d *DB) Ping() error {
	if err := d.DB.PingContext(d.ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	var one int8
	if err := d.DB.QueryRowContext(d.ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("select 1: %w", err)
	}
	if one != 1 {
		return fmt.Errorf("unexpected ping result: %d", one)
	}
	return nil
}
