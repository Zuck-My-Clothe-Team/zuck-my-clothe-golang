package platform

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
}

func InitDB(dsnConfigString string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsnConfigString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Users{})

	return &Postgres{db}, nil
}
