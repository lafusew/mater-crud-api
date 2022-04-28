package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model // extends Model that contains ID, CreatedAt, UpdatedAt, DeletedAt
	Body      string    `gorm:"size:24;not null;unique" json:"body"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	Post      Post      `json:"post"`
	PostID    uint32    `gorm:"not null" json:"post_id"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.Post = Post{}
	c.Author = User{}
	c.CreatedAt = time.Now()
}

func (c *Comment) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
	default:
		if c.Body == "" {
			return errors.New("required body")
		}
	}

	return nil
}

func (c *Comment) SaveTeam(db *gorm.DB) (*Comment, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}

	return c, nil
}

func (c *Comment) FindAll(db *gorm.DB) (*[]Comment, error) {
	comments := []Comment{}
	err := db.Debug().Model(&Comment{}).Limit(100).Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}

	if len(comments) > 0 {
		for i := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].AuthorID).Take(&c.Author).Error
			if err != nil {
				return &[]Comment{}, nil
			}

			err = db.Debug().Model(&Post{}).Where("id = ?", comments[i].PostID).Take(&c.Post).Error
			if err != nil {
				return &[]Comment{}, nil
			}
		}
	}
	return &comments, nil
}

func (c *Comment) FindByID(db *gorm.DB, cid uint32) (*Comment, error) {
	err := db.Debug().Model(Comment{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Comment{}, errors.New("comment not found")
	}
	return c, nil
}

func (c *Comment) DeleteComment(db *gorm.DB, cid uint32) (int64, error) {
	db = db.Debug().Model(&Comment{}).Where("id = ?", cid).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("comment not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
