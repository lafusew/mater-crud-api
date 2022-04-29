package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Code   string   `gorm:"size:6;not null;unique" json:"code"`
	Name   string   `gorm:"size:16;not null" json:"name"`
	Desc   string   `json:"desc"`
	Teams  []*Team  `gorm:"many2many:team_tags"`
	Events []*Event `gorm:"many2many:event_tags"`
}

func (t *Tag) Prepare() {
	t.ID = 0
	t.Code = html.EscapeString(strings.TrimSpace(t.Code))
	t.CreatedAt = time.Now()
	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.Desc = html.EscapeString(strings.TrimSpace(t.Desc))
}

func (t *Tag) Validate(action string) error {
	return nil
}

func (t *Tag) SaveTag(db *gorm.DB) (*Tag, error) {
	err := db.Debug().Create(&t).Error
	if err != nil {
		return &Tag{}, err
	}

	return t, nil
}

func (t *Tag) FindAll(db *gorm.DB) (*[]Tag, error) {
	tags := []Tag{}

	err := db.Debug().Model(&Tag{}).Limit(100).Find(&tags).Error
	if err != nil {
		return &[]Tag{}, err
	}

	return &tags, nil
}

func (t *Tag) FindByID(db *gorm.DB, tid uint32) (*Tag, error) {
	err := db.Debug().Model(Tag{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Tag{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Tag{}, errors.New("comment not found")
	}
	return t, nil
}

func (c *Tag) DeleteComment(db *gorm.DB, tid uint32) (int64, error) {
	db = db.Debug().Model(&Tag{}).Where("id = ?", tid).Take(&Tag{}).Delete(&Tag{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("tag not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
