package database

import (
	"log"
	"password_store/internal/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func (d *Database) Init() {
	pg_user, err := util.GetEnv("POSTGRES_USER")
	if err != nil {
		log.Fatalf(err.Error())
	}

	pg_password, err := util.GetEnv("POSTGRES_PASSWORD")
	if err != nil {
		log.Fatalf(err.Error())
	}

	pg_db, err := util.GetEnv("POSTGRES_DB")
	if err != nil {
		log.Fatalf(err.Error())
	}

	pg_port, err := util.GetEnv("POSTGRES_PORT")
	if err != nil {
		log.Fatalf(err.Error())
	}

	dsn := "host=db" + " user=" + pg_user + " password=" + pg_password + " dbname=" + pg_db + " port=" + pg_port + " sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	d.Db = db
}

func (d *Database) AutoMigrate() {
	d.Db.AutoMigrate(StoredUserCredentials{})
	d.Db.AutoMigrate(UserBalance{})
	d.Db.AutoMigrate(UserTransaction{})
}

var defaultPassword string = "password"
var defaultSalt string = "defaultSalt"

var seedUserCredentials = []StoredUserCredentials{
	{Username: "alice", Salt: defaultSalt, Hash: util.Hash([]byte(defaultPassword + defaultSalt))},
	{Username: "bob", Salt: defaultSalt, Hash: util.Hash([]byte(defaultPassword + defaultSalt))},
}

var seedUserBalance = []UserBalance{
	{Username: "alice", Balance: 100000},
	{Username: "bob", Balance: 100000},
}

func (d *Database) Seed() {
	for _, cred := range seedUserCredentials {
		d.Db.Create(&cred)
	}

	for _, balance := range seedUserBalance {
		d.Db.Create(&balance)
	}
}
