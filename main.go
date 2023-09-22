package main

import (
	"fmt"
	"log"
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
	"gorm.io/gorm/logger"
)

const TimeFormat = "2006-01-02 15:04:05"

var env = &models.GlobalEnvironment

var configuration *common.Configuration

var collector *colly.Collector

var tLog = logrus.New()

func main() {
	initConfig()

	initEnvironment()

	initDB()

	go initCollector()

	initRouter()
}

func initConfig() {
	var err error

	configuration, err = common.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}
}

func initEnvironment() {
	env.MoviesManager = service.NewMovieManger()
	env.CategoryManager = service.NewCategoryManager()
}

func initDB() {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configuration.DBName, configuration.DBPassword, configuration.DBHost, configuration.DBTable)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: false,
		Logger: logger.New(log.New(os.Stderr, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		panic(err)
	}

	common.SetDatabase(db)

	db.Table("category").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.Category{})
	db.Table("movie_players").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.MoviePlayer{})
	db.Table("movie_views").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&models.MovieView{})
}

func initCollector() {
	start := time.Now()

	collector = colly.NewCollector(nil)

	if err := readLockFile(); err != nil {
		tLog.Warn(err)
	}

	if err := collector.Colly(); err != nil {
		tLog.Errorf("Colly init err: %+v", err)
		return
	}

	tLog.Infof("Database init Finished, speed: %+v", time.Now().Sub(start))

	writeLockFile()
}

func readLockFile() error {
	latest, err := os.ReadFile("./latest.lock")
	if err != nil {
		return err
	}

	t, err := time.Parse(TimeFormat, string(latest))
	if err != nil {
		return err
	}

	collector.SetLatestUpdate(t)
	return nil
}

func writeLockFile() {
	f, err := os.Create("./latest.lock")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write([]byte(time.Now().Format(TimeFormat)))
}

func initRouter() error {
	engine := gin.Default()

	for _, router := range api.Routers {
		engine.Handle(router.Method, router.Path, router.Handler)
	}

	return engine.Run(":" + configuration.ListenPort)
}
