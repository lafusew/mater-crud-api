package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Name           string `gorm:"size:24;not null;unique" json:"name"`
	Desc           string `gorm:"size:255;" json:"desc"`
	IsPrivate      bool   `gorm:"default:false" json:"is_private"`
	ProfilePicture string `gorm:"size:255;" json:"profile_picture"`
	CreatorID      uint   `gorm:"not null" json:"creator_id"`
	TeamID         uint   `gorm:"not null" json:"team_id"`
	BeginsAt       time.Time
	EndsAt         time.Time
	Tags           []*Tag `gorm:"many2many:event_tags"`
}

func (e *Event) Prepare() {
	e.ID = 0
	e.Name = html.EscapeString(strings.TrimSpace(e.Name))
	e.Desc = html.EscapeString(strings.TrimSpace(e.Desc))
	e.CreatedAt = time.Now()
}

func (e *Event) Validate(action string) error {
	return nil
}

func (e *Event) Save(db *gorm.DB) (*Event, error) {
	err := db.Debug().Create(&e).Error
	if err != nil {
		return &Event{}, err
	}

	return e, nil
}

func (e *Event) FindAll(db *gorm.DB) (*[]Event, error) {
	tags := []Event{}

	err := db.Debug().Model(&Event{}).Limit(100).Find(&tags).Error
	if err != nil {
		return &[]Event{}, err
	}

	return &tags, nil
}

func (e *Event) FindByID(db *gorm.DB, tid uint32) (*Event, error) {
	err := db.Debug().Model(Event{}).Where("id = ?", tid).Take(&e).Error
	if err != nil {
		return &Event{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Event{}, errors.New("comment not found")
	}
	return e, nil
}

func (c *Event) Delete(db *gorm.DB, tid uint32) (int64, error) {
	db = db.Debug().Model(&Event{}).Where("id = ?", tid).Take(&Event{}).Delete(&Event{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("event not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
