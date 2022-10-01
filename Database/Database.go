package Database

import (
	"backend-d-embung/Entities"
	"fmt"
	"log"
	"os"

	_"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() *gorm.DB {
	var db *gorm.DB
	var err error

	// Buka Koneksi (using mariadb/mysql)
	// db, err = gorm.Open(
	// 	mysql.Open(
	// 		fmt.Sprintf(
	// 			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
	// 			os.Getenv("DB_USER"),
	// 			os.Getenv("DB_PASS"),
	// 			os.Getenv("DB_HOST"),
	// 			os.Getenv("DB_NAME"),
	// 		),
	// 	),
	// 	&gorm.Config{})
	dsn := fmt.Sprintf("user=%s "+
			"password=%s "+
			"host=%s "+
			"TimeZone=Asia/Singapore "+
			"port=%s "+
			"dbname=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Model
	if err = db.AutoMigrate(
		Entities.Artikel{}, 
		Entities.Operasional{}, 
		Entities.Testimoni{},
	); 
	err != nil {
		log.Fatal(err.Error())
	}

	return db
}
