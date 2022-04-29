package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model        // extends Model that contains ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string `gorm:"size:255;not null;unique" json:"title"`
	Content    string `gorm:"size:255;not null;" json:"content"`
	AuthorID   uint   `gorm:"not null" json:"author_id"`
	Comments   []Comment
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {

	if p.Title == "" {
		return errors.New("required title")
	}
	if p.Content == "" {
		return errors.New("required content")
	}
	if p.AuthorID < 1 {
		return errors.New("required author")
	}
	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	err := db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	posts := []Post{}

	err := db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{
		Title:   p.Title,
		Content: p.Content,
	}).Error
	if err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint) (int64, error) {
	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
