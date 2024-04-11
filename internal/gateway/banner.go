package gateway

import (
	"BannerService/internal/models"
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

	createQuery := "INSERT INTO banners (version, feature_id, content, is_active) VALUES ($1, $2, $3, $4) RETURNING id"
	row := tx.QueryRow(createQuery, banner.Version, banner.FeatureId, banner.Content, banner.IsActive)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
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

func (p *BannerPostgres) GetBanner(tagId, featureId, limit, offset int32) ([]models.Banner, error) {
	var query string
	var args []interface{}

	query = `
        SELECT 
            b.id, b.version, b.feature_id, b.content, b.is_active, b.created_at, b.update_at, 
            ARRAY_AGG(tb.tag_id) AS tag_ids
        FROM 
            Banners b
        JOIN 
            tags_banners tb ON b.id = tb.banner_id AND b.version = tb.banner_version
        WHERE 
            1=1
    `

	if featureId != 0 {
		query += fmt.Sprintf(" AND b.feature_id = $%d", len(args)+1)
		args = append(args, featureId)
	}
	if tagId != 0 {
		query += fmt.Sprintf(" AND tb.tag_id = $%d", len(args)+1)
		args = append(args, tagId)
	}

	query += `
        GROUP BY 
            b.id, b.version, b.feature_id, b.content, b.is_active, b.created_at, b.update_at
    `

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}
	if offset != 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	// Создаем запрос для получения всех тегов для каждого баннера
	subQuery := `
        SELECT 
            banner_id, ARRAY_AGG(tag_id) AS tag_ids
        FROM 
            tags_banners
        GROUP BY 
            banner_id
    `

	// Добавляем подзапрос в основной запрос
	query = fmt.Sprintf(`
        SELECT 
            b.*, t.tag_ids
        FROM 
            (%s) AS b
        JOIN 
            (%s) AS t ON b.id = t.banner_id
    `, query, subQuery)

	res := make([]models.Banner, 0, limit)

	if err := p.db.Select(&res, query, args...); err != nil {
		return nil, err
	}

	return res, nil
}
