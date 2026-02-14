package main

import (
	migrations "analog-be/migration"
	"context"
	"crypto/tls"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5437")
	user := GetEnv("DB_USER", "test")
	password := GetEnv("DB_PASSWORD", "test")
	database := GetEnv("DB_NAME", "test")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: os.Getenv("DB_TLS_SKIP_VERIFY") == "true",
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(host+":"+port),
		pgdriver.WithTLSConfig(tlsConfig),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
		pgdriver.WithApplicationName("analog"),
		pgdriver.WithTimeout(10*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(10*time.Second),
		pgdriver.WithWriteTimeout(10*time.Second),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(pgconn)

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	ctx := context.Background()

	if os.Args[1] == "init" {
		err = migrator.Init(ctx)
		if err != nil {
			panic(err)
		}
	}

	if os.Args[1] == "up" {
		group, err := migrator.Migrate(ctx)
		if err != nil {
			panic(err)
		}
		if group.IsZero() {
			println("변경 사항이 없습니다.")
		} else {
			println("마이그레이션 완료:", group.String())
		}
	}

}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
