// internal/db/db.go
package db

import (
	"SomeProject/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"strings"
)

type DBPool struct {
	Pool *pgxpool.Pool
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func NewConfig() *Config {
	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}

func CreateDatabase(cfg *Config) error {
	systemDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)

	conn, err := pgx.Connect(context.Background(), systemDSN)
	if err != nil {
		return fmt.Errorf("unable to connect to system database: %v", err)
	}
	defer conn.Close(context.Background())

	createDBQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", cfg.DBName)
	if _, err := conn.Exec(context.Background(), createDBQuery); err != nil {
		return fmt.Errorf("failed to create database \"%s\": %v", cfg.DBName, err)
	}

	return nil
}

func ConnectDB(cfg *Config) (*DBPool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	dbPool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &DBPool{Pool: dbPool}, nil
}

func LoadBannerByParams(dbPool *pgxpool.Pool, tagID, featureID string, useLastRevision bool) (models.Banner, error) {
	var banner models.Banner

	query := `
    SELECT Title, Text, URL
    FROM banners 
    WHERE tag_id=$1 AND feature_id=$2 
    ORDER BY created_at DESC`

	var err error
	if useLastRevision {
		query += " LIMIT 1"
		err = dbPool.QueryRow(context.Background(), query, tagID, featureID).Scan(&banner.Title, &banner.Text, &banner.URL)
	} else {
		err = dbPool.QueryRow(context.Background(), query, tagID, featureID).Scan(&banner.Title, &banner.Text, &banner.URL)
	}

	if err != nil {
		return models.Banner{}, fmt.Errorf("failed to load banner: %v", err)
	}

	return banner, nil
}

func InsertBanner(dbPool *pgxpool.Pool, banner models.Banner) (int, error) {
	var bannerID int
	query := `INSERT INTO banners1 (Title, Text, URL, TagIDs, FeatureID, IsActive) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := dbPool.QueryRow(context.Background(), query, banner.Title, banner.Text, banner.URL, banner.TagIDs, banner.FeatureID, banner.IsActive).Scan(&bannerID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert banner: %v", err)
	}
	return bannerID, nil
}

func UpdateBanner(dbPool *pgxpool.Pool, bannerID int, update models.BannerUpdate) error {
	var setClauses []string
	var args []interface{}
	args = append(args, bannerID)

	if update.TagIDs != nil {
		setClauses = append(setClauses, fmt.Sprintf("TagIDs = $%d", len(args)+1))
		args = append(args, update.TagIDs)
	}
	if update.FeatureID != nil {
		setClauses = append(setClauses, fmt.Sprintf("FeatureID = $%d", len(args)+1))
		args = append(args, update.FeatureID)
	}
	if update.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("Title = $%d", len(args)+1))
		args = append(args, *update.Title)
	}
	if update.Text != nil {
		setClauses = append(setClauses, fmt.Sprintf("Text = $%d", len(args)+1))
		args = append(args, *update.Text)
	}
	if update.URL != nil {
		setClauses = append(setClauses, fmt.Sprintf("URL = $%d", len(args)+1))
		args = append(args, *update.URL)
	}
	if update.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("IsActive = $%d", len(args)+1))
		args = append(args, update.IsActive)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields specified for update")
	}

	query := fmt.Sprintf("UPDATE banners1 SET %s WHERE id = $1", strings.Join(setClauses, ", "))
	cmdTag, err := dbPool.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no banner found with ID %d", bannerID)
	}

	return nil
}

func DeleteBanner(dbPool *pgxpool.Pool, bannerID int) error {

	query := `DELETE FROM banners WHERE id = $1`

	cmdTag, err := dbPool.Exec(context.Background(), query, bannerID)
	if err != nil {
		return fmt.Errorf("error deleting banner with ID %d: %v", bannerID, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no banner found with ID %d", bannerID)
	}

	return nil
}
