package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/MK-BK/tk-colly/api"
	"github.com/MK-BK/tk-colly/colly"
	"github.com/MK-BK/tk-colly/database"
	"github.com/MK-BK/tk-colly/job"
	"github.com/MK-BK/tk-colly/models"
	"github.com/MK-BK/tk-colly/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var env = &models.GlobalEnvironment

var configuration *models.Configuration

func main() {
	env.MovieInterface = service.NewMovieManger()

	loadConfiguration()

	initDB()

	initTables()

	if err := readLockFile(); err != nil {
		fmt.Println(err)
	}

	now := time.Now()
	if err := colly.Colly(); err != nil {
		fmt.Printf("Colly init err: %+v", err)
	}

	fmt.Printf("Database init Finished, Speed %+v\n", time.Now().Sub(now))

	writeLockFile()

	initJob()

	initRouter()
}

func loadConfiguration() {
	config, err := models.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	configuration = config
}

func initDB() {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configuration.DBName, configuration.DBPassword, configuration.DBHost, configuration.DBTable)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.Set(db)
}

func initTables() {
	database.DB.Table("movies").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Movie{})
	database.DB.Table("movie_views").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.MovieView{})
}

func readLockFile() error {
	latest, err := ioutil.ReadFile("./latest.lock")
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02 15:04:05", string(latest))
	if err != nil {
		return err
	}

	colly.SetLatestUpdate(t)
	return nil
}

func writeLockFile() {
	f, err := os.Create("./latest.lock")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
}

func initJob() {
	job.Schedule(24 * time.Hour)
}

func initRouter() error {
	engine := gin.Default()

	for _, router := range api.Routers {
		engine.Handle(router.Method, router.Path, router.Handler)
	}

	return engine.Run(":" + configuration.ListenPort)
}
