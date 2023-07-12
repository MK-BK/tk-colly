package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/MK-BK/tk-colly/api"
	"github.com/MK-BK/tk-colly/colly"
	"github.com/MK-BK/tk-colly/common"
	"github.com/MK-BK/tk-colly/models"
	"github.com/MK-BK/tk-colly/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var env = &models.GlobalEnvironment
var configuration *common.Configuration

var log = logrus.New()

func main() {
	initConfig()

	initEnvironment()

	initDB()

	// func() {
	// 	if err := readLockFile(); err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	if err := colly.Colly(); err != nil {
	// 		panic(fmt.Sprintf("Colly init err: %+v", err))
	// 	} else {
	// 		log.Infof("Database init Finished")
	// 	}

	// 	writeLockFile()
	// }()

	initRouter()
}

func initEnvironment() {
	env.MoviesManager = service.NewMovieManger()
	env.CategoryManager = service.NewCategoryManager()
}

func initConfig() {
	var err error
	configuration, err = common.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}
}

func initDB() {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configuration.DBName, configuration.DBPassword, configuration.DBHost, configuration.DBTable)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: false,
	})
	if err != nil {
		panic(err)
	}

	common.SetDatabase(db)

	db.Table("movies").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Movie{})
	db.Table("movie_views").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.MovieView{})
	db.Table("category").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Category{})
}

func readLockFile() error {
	latest, err := ioutil.ReadFile("./latest.lock")
	if err != nil {
		return err
	}

	t, err := time.Parse(models.TimeFormat, string(latest))
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

	f.Write([]byte(time.Now().Format(models.TimeFormat)))
}

func initRouter() error {
	engine := gin.Default()

	for _, router := range api.Routers {
		engine.Handle(router.Method, router.Path, router.Handler)
	}

	return engine.Run(":" + configuration.ListenPort)
}
