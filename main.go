package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/thecactusblue/hana-id/model"
	"github.com/thecactusblue/hana-id/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setup() error {
	db, err := gorm.Open(postgres.Open(viper.GetString("DATABASE")), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.RefreshToken{},
		&model.FlowState{},
	)
	if err != nil {
		return err
	}

	e := echo.New()
	web.Inject(e, db)

	return e.Start(viper.GetString("ADDRESS"))
}

func readConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	readConfig()
	err := setup()
	if err != nil {
		panic(err)
	}
}
