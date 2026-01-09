package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	DBName   string `env:"DB_DBNAME,required"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	TimeZone string `env:"DB_TIMEZONE" envDefault:"Asia/Taipei"`

	MaxOpenConn     int           `env:"DB_MAX_OPEN" envDefault:"50"`
	MaxIdleConn     int           `env:"DB_MAX_IDLE" envDefault:"20"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE" envDefault:"5m"`
}

type DB struct {
	SQL     *sql.DB
	Builder sq.StatementBuilderType
}

func New(cfg Config) (*DB, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.TimeZone == "" {
		cfg.TimeZone = "Asia/Taipei"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.TimeZone,
	)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if cfg.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	} else {
		sqlDB.SetMaxOpenConns(20)
	}
	if cfg.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	} else {
		sqlDB.SetMaxIdleConns(10)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
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
