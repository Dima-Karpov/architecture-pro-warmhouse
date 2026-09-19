package db

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var errEmptyDatabase = errors.New("empty database name")

func Open(ctx context.Context, databaseURL string) (*gorm.DB, error) {
	name, err := databaseName(databaseURL)
	if err != nil {
		return nil, err
	}

	if err = ensureDatabase(ctx, databaseURL, name); err != nil {
		return nil, err
	}

	conn, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err = migrateIfNeeded(ctx, conn, &readingRow{}); err != nil {
		Close(conn)

		return nil, err
	}

	return conn, nil
}

func migrateIfNeeded(ctx context.Context, conn *gorm.DB, model any) error {
	if conn.WithContext(ctx).Migrator().HasTable(model) {
		return nil
	}

	if err := conn.WithContext(ctx).AutoMigrate(model); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func Close(conn *gorm.DB) {
	if conn == nil {
		return
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return
	}

	if err = sqlDB.Close(); err != nil {
		return
	}
}

func databaseName(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}

	name := strings.TrimPrefix(parsed.Path, "/")
	if name == "" {
		return "", errEmptyDatabase
	}

	return name, nil
}

func ensureDatabase(ctx context.Context, databaseURL, name string) error {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}

	parsed.Path = "/postgres"

	admin, err := gorm.Open(postgres.Open(parsed.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("connect admin database: %w", err)
	}
	defer Close(admin)

	var exists bool

	err = admin.WithContext(ctx).
		Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", name).
		Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("check database: %w", err)
	}

	if exists {
		return nil
	}

	quoted, err := quoteIdent(name)
	if err != nil {
		return err
	}

	if err = admin.Exec("CREATE DATABASE " + quoted).Error; err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func quoteIdent(name string) (string, error) {
	if name == "" || strings.ContainsAny(name, "\"\\; ") {
		return "", errEmptyDatabase
	}

	for _, r := range name {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return "", errEmptyDatabase
		}
	}

	return name, nil
}
