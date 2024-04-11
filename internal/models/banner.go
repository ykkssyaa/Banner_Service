package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ModelMap map[string]interface{}
type Tags []int32

type Banner struct {
	// Идентификатор баннера
	Id int32 `json:"banner_id,omitempty" db:"id"`
	// Версия баннера
	Version int32 `json:"version,omitempty" db:"version"`
	// Идентификаторы тэгов
	TagIds Tags `json:"tag_ids,omitempty" db:"tag_ids"`
	// Идентификатор фичи
	FeatureId int `json:"feature_id,omitempty" db:"feature_id"`
	// Содержимое баннера
	Content ModelMap `json:"content,omitempty" db:"content"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty" db:"is_active"`
	// Дата создания баннера
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	// Дата обновления баннера
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"update_at"`
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

func (m ModelMap) Scan(data interface{}) error {
	// Проверяем, что данные не nil
	if data == nil {
		return nil
	}

	// Проверяем, что данные являются []byte
	bytes, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert data to []byte")
	}

	// Декодируем JSON и присваиваем значение указателю на ModelMap
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	return nil
}

func (m ModelMap) Equal(other ModelMap) bool {
	if len(m) != len(other) {
		return false
	}
	for key, val := range m {
		if v, ok := other[key]; !ok {
			return false
		} else {
			if v != val {
				return false
			}
		}
	}

	return true
}

func (t *Tags) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	// Преобразование src в []byte
	b, ok := src.([]byte)
	if !ok {
		return errors.New("unexpected type for Tags")
	}

	// Преобразование []byte в строку
	s := string(b)

	// Удаление фигурных скобок из строки
	s = strings.Trim(s, "{}")

	// Разбивка строки по запятым
	parts := strings.Split(s, ",")

	// Создание слайса для тэгов
	tags := make([]int32, len(parts))

	// Преобразование строковых значений в int32
	for i, part := range parts {
		tag, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return err
		}
		tags[i] = int32(tag)
	}

	// Присвоение значения Tags
	*t = tags
	return nil
}

func (t *Tags) Value() (driver.Value, error) {
	if len(*t) == 0 {
		return nil, nil
	}
	return json.Marshal(*t)
}

func (t *Tags) Equal(other Tags) bool {
	if len(*t) != len(other) {
		return false
	}
	for i, val := range *t {
		if val != other[i] {
			return false
		}
	}
	return true
}
