package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"my_app/internal/cache"
	"my_app/internal/models"
	"os"
	"strings"

	pq "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

func InitDB() {
	var err error
	param := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err = sql.Open("postgres", param)
	if err != nil {
		log.Fatal(err)
	}

	createBannersTable()
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing PostgreSQL connection: %v", err)
	}
}

func createBannersTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS banners (
		id SERIAL PRIMARY KEY,
		tag_ids INT[] NOT NULL,
		feature_id INT NOT NULL,
		content JSON NOT NULL,
		is_active BOOLEAN NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetBannerForUser(featureId *int, tagId *int, use_last_revision bool, isAdmin bool) (*models.ModelMap, error) {
	var banner *models.BannerExpanded
	var err error
	if !use_last_revision {
		banner, err = cache.GetBannerFromCache(featureId, tagId)
	}
	if use_last_revision || err != nil && strings.Contains(err.Error(), "no banner found") {
		banner, err = getBannerFromDB(featureId, tagId)
		if err != nil {
			return nil, err
		}
		cache.SaveBannerToCacheAsync(featureId, tagId, banner)
	} else if err != nil {
		return nil, err
	}
	if !banner.IsActive && !isAdmin {
		return nil, fmt.Errorf("no banner found")
	}

	return &banner.Content, nil
}

func getBannerFromDB(featureId *int, tagId *int) (*models.BannerExpanded, error) {
	var banner models.BannerExpanded
	query := `SELECT * FROM banners WHERE 1=1`
	if featureId != nil {
		query += fmt.Sprintf(" AND feature_id = %d", *featureId)
	}
	if tagId != nil {
		query += fmt.Sprintf(" AND %d = ANY(tag_ids)", *tagId)
	}
	query += " LIMIT 1"

	row := db.QueryRow(query)
	var contentJSON []byte
	err := row.Scan(&banner.ID, pq.Array(&banner.TagIds), &banner.FeatureId, &contentJSON, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no banner found")
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

func GetBanners(featureId *int, tagId *int, limit *int, offset *int) ([]models.BannerExpanded, error) {
	var banners []models.BannerExpanded

	query := `SELECT * FROM banners WHERE 1=1`
	if featureId != nil {
		query += fmt.Sprintf(" AND feature_id = %d", *featureId)
	}
	if tagId != nil {
		query += fmt.Sprintf(" AND %d = ANY(tag_ids)", *tagId)
	}
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *offset)
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banner models.BannerExpanded
		var contentJSON []byte
		err := rows.Scan(&banner.ID, pq.Array(&banner.TagIds), &banner.FeatureId, &contentJSON, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(contentJSON, &banner.Content)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return banners, nil
}

func CreateBanner(banner models.BannerNoId) (int32, error) {
	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO banners (tag_ids, feature_id, content, is_active) VALUES ($1, $2, $3, $4) RETURNING id`
	var bannerId int32
	err = db.QueryRow(query, pq.Array(banner.TagIds), banner.FeatureId, contentJSON, banner.IsActive).Scan(&bannerId)
	if err != nil {
		return 0, err
	}
	return bannerId, nil
}

func UpdateBanner(id int, banner models.BannerNoId) error {
	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return err
	}
	query := `UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4, updated_at = NOW() WHERE id = $5`
	_, err = db.Exec(query, pq.Array(banner.TagIds), banner.FeatureId, contentJSON, banner.IsActive, id)
	return err
}

func DeleteBanner(id int) error {
	query := `DELETE FROM banners WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func BannerExists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM banners WHERE id = $1)`
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
