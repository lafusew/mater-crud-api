package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model             // extends Model that contains ID, CreatedAt, UpdatedAt, DeletedAt
	Name           string  `gorm:"size:24;not null;unique" json:"name"`
	Desc           string  `gorm:"size:255;" json:"desc"`
	ProfilePicture string  `gorm:"size:255;" json:"profile_picture"`
	IsPrivate      bool    `gorm:"default:false" json:"is_private"`
	Members        []*User `gorm:"many2many:membership"`
	Tags           []*Tag  `gorm:"many2many:team_tags"`
}

func (t *Team) Prepare() {
	t.ID = 0
	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.Desc = html.EscapeString(strings.TrimSpace(t.Desc))
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Team) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
	default:
		if t.Name == "" {
			return errors.New("required name")
		}
	}

	return nil
}

func (t *Team) SaveTeam(db *gorm.DB) (*Team, error) {
	err := db.Debug().Create(&t).Error
	if err != nil {
		return &Team{}, err
	}

	return t, nil
}

func (t *Team) FindAll(db *gorm.DB) (*[]Team, error) {
	teams := []Team{}
	err := db.Debug().Model(&Team{}).Limit(100).Find(&teams).Error
	if err != nil {
		return &[]Team{}, err
	}

	return &teams, nil
}

func (t *Team) FindByID(db *gorm.DB, tid uint32) (*Team, error) {
	err := db.Debug().Model(Team{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Team{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Team{}, errors.New("Team not found")
	}
	return t, nil
}

func (t *Team) UpdateTeam(db *gorm.DB, tid uint32) (*Team, error) {
	err := db.Debug().Model(&Team{}).Where("id = ?", tid).Take(&Team{}).Updates(Team{
		Name:           t.Name,
		Desc:           t.Desc,
		ProfilePicture: t.ProfilePicture,
		IsPrivate:      t.IsPrivate,
	}).Error
	if err != nil {
		return &Team{}, err
	}

	return t, nil
}

func (t *Team) DeleteTeam(db *gorm.DB, tid uint32) (int64, error) {
	db = db.Debug().Model(&Team{}).Where("id = ?", tid).Take(&Team{}).Delete(&Team{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Team not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
