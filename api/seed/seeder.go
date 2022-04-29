package seed

import (
	"log"

	"github.com/gob-mater/app/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Name: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Name: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

var comments = []models.Comment{
	{
		Body: "Comment body1",
	},
	{
		Body: "Comment body2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists("membership", &models.Comment{},  &models.Team{}, &models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID
		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		comments[i].AuthorID = users[i].ID
		comments[i].PostID = posts[i].ID
		err = db.Debug().Model(&models.Comment{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}