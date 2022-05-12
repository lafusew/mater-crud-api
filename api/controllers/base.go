package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/gob-mater/app/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database \n", Dbdriver)
		}
	}

	if true {
		err = server.DB.Debug().DropTableIfExists(
			"team_tags",
			"event_tags",
			&models.Event{},
			"membership",
			&models.Comment{},
			&models.Team{},
			&models.Post{},
			&models.Tag{},
			&models.User{},
		).Error
		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}
	}

	server.DB.Debug().AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Team{},
		&models.Comment{},
		&models.Tag{},
		&models.Event{},
	)

	InitializeForeignKeys(server.DB)

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func InitializeForeignKeys(db *gorm.DB) error {
	defer fmt.Println("fkeys initialized")

	// POST
	err := db.Debug().Model(models.Post{}).AddForeignKey("author_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to posts %v", err)
	}

	// COMMENT
	err = db.Debug().Model(models.Comment{}).AddForeignKey("author_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to comments %v", err)
	}

	err = db.Debug().Model(models.Comment{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to comments: %v", err)
	}

	// MEMBERSHIP JOIN TABLE RELATIONS
	err = db.Debug().Table("membership").AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}

	err = db.Debug().Table("membership").AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}

	// TEAM_TAGS JOIN TABLE RELATIONS
	err = db.Debug().Table("team_tags").AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}

	err = db.Debug().Table("team_tags").AddForeignKey("tag_id", "tags(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}

	err = db.Debug().Table("event_tags").AddForeignKey("event_id", "events(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}

	err = db.Debug().Table("event_tags").AddForeignKey("tag_id", "tags(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return fmt.Errorf("couldn't add foreign key to membership: %v", err)
	}
	
	return nil
}