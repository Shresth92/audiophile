package database

import (
	"fmt"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDb() {
	err := utils.LoadEnv()
	host := utils.GetEnvValue("dbHost")
	user := utils.GetEnvValue("dbUser")
	password := utils.GetEnvValue("dbPassword")
	dbname := utils.GetEnvValue("dbName")
	dbPort := utils.GetEnvValue("dbPort")

	if err != nil {
		logrus.Errorf("Environment loading error; err: %s", err.Error())
	}
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, dbPort)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		logrus.Errorf("database connection failed; err: %s", err.Error())
	}
	Db = db
	if err := db.Exec("CREATE TYPE role_type AS ENUM ('admin','user')").Error; err != nil {
		logrus.Errorf("enum creation failed; err: %s", err)
	}
	if err := db.Exec("CREATE TYPE delivery_status AS ENUM ('onTheWay','delivered','canceled','return')").Error; err != nil {
		logrus.Errorf("enum creation failed; err: %s", err)
	}
	if err := db.AutoMigrate(&models.Users{}, &models.UserRole{}, &models.Address{}, &models.Session{}, &models.Category{}, &models.Brand{}, &models.Product{}, &models.Variants{}, &models.Offer{}, &models.Images{}, &models.VariantImages{}, &models.Orders{}, &models.ProductOrdered{}, &models.UserCart{}); err != nil {
		logrus.Errorf("automigration failed; err: %s", err.Error())
	}
}

func CloseDb() {
	DbInstance, _ := Db.DB()
	DbInstance.Close()
}
