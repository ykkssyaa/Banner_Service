package gateway

import (
	"BannerService/internal/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type BannerPostgres struct {
	db *sqlx.DB
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}

func (p *BannerPostgres) CreateBanner(banner models.Banner) (int, error) {

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	if banner.Id == 0 {

		createQuery := "INSERT INTO banners (version, feature_id, content, is_active) VALUES ($1, $2, $3, $4) RETURNING id"
		row := tx.QueryRow(createQuery, banner.Version, banner.FeatureId, banner.Content, banner.IsActive)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		createQuery := "INSERT INTO banners (id, version, feature_id, content, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err := tx.Exec(createQuery, banner.Id, banner.Version, banner.FeatureId, banner.Content, banner.IsActive, banner.CreatedAt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		id = int(banner.Id)
	}

	createTagsRelQuery := "INSERT INTO tags_banners (banner_id, banner_version, tag_id) VALUES ($1, $2, $3)"

	for _, tagId := range banner.TagIds {
		if _, err := tx.Exec(createTagsRelQuery, id, banner.Version, tagId); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return id, tx.Commit()
}

func (p *BannerPostgres) GetBanner(tagId, featureId, limit, offset int32, isActive *bool) ([]models.Banner, error) {
	var query string
	var args []interface{}

	query = `
        SELECT b.id, b.version, b.feature_id, b.content, b.is_active, b.created_at, b.update_at, 
            ARRAY_AGG(tb.tag_id) AS tag_ids
        FROM Banners b
        JOIN tags_banners tb ON b.id = tb.banner_id AND b.version = tb.banner_version
        WHERE ((b.version) IN (
             SELECT MAX(version)
             FROM Banners
             WHERE id = b.id) OR b.is_active)
    `

	if featureId != 0 {
		query += fmt.Sprintf(" AND b.feature_id = $%d", len(args)+1)
		args = append(args, featureId)
	}
	if tagId != 0 {
		query += fmt.Sprintf(" AND tb.tag_id = $%d", len(args)+1)
		args = append(args, tagId)
	}
	if isActive != nil {
		query += fmt.Sprintf(" AND b.is_active = $%d", len(args)+1)
		args = append(args, *isActive)
	}

	query += `
        GROUP BY b.id, b.version, b.feature_id, b.content, b.is_active, b.created_at, b.update_at
    `

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	if offset != 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	subQuery := `
        SELECT banner_id, banner_version, ARRAY_AGG(tag_id) AS tag_ids
        FROM tags_banners
        GROUP BY banner_id, banner_version
    `

	query = fmt.Sprintf(`
        SELECT  b.*, t.tag_ids
        FROM (%s) AS b
        JOIN (%s) AS t ON b.id = t.banner_id AND b.version = t.banner_version
        ORDER BY b.id , b.version desc 
    `, query, subQuery)

	res := make([]models.Banner, 0, limit)

	if err := p.db.Select(&res, query, args...); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *BannerPostgres) GetBannerById(id int32) (models.Banner, error) {
	var banner models.Banner

	err := p.db.Get(&banner, "SELECT * FROM banners WHERE id = $1 ORDER BY version desc LIMIT 1;", id)
	if err != nil {
		return models.Banner{}, err
	}

	return banner, nil
}

func (p *BannerPostgres) DeleteBanner(id int32) error {

	tx, err := p.db.Begin()

	if err != nil {
		return err
	}

	result, err := tx.Exec("DELETE FROM banners WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec("DELETE FROM tags_banners WHERE banner_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *BannerPostgres) SetActiveVersion(id, version int32, isActive bool) error {

	tx, err := p.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE banners SET is_active = $1 WHERE id = $2 AND version != $3",
		false, id, version)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE banners SET is_active = $1 WHERE id = $2 AND version = $3",
		isActive, id, version)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
