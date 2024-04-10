package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ModelMap map[string]interface{}

type Banner struct {
	// Идентификатор баннера
	BannerId int32 `json:"banner_id,omitempty"`
	// Версия баннера
	Version int32 `json:"version,omitempty"`
	// Идентификаторы тэгов
	TagIds []int32 `json:"tag_ids,omitempty"`
	// Идентификатор фичи
	FeatureId int32 `json:"feature_id,omitempty"`
	// Содержимое баннера
	Content ModelMap `json:"content,omitempty"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
	// Дата создания баннера
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Дата обновления баннера
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (m ModelMap) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err.Error())
	}
	return string(data)
}

func (m ModelMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m ModelMap) Scan(data []byte) error {
	return json.Unmarshal(data, &m)
}
