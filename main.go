package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/phantom-atom/file-explorer/repository/simple"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phantom-atom/file-explorer/cache"
	"github.com/phantom-atom/file-explorer/cache/redis"
	"github.com/phantom-atom/file-explorer/config"
	"github.com/phantom-atom/file-explorer/internal/locker"
	"github.com/phantom-atom/file-explorer/internal/log"
	"github.com/phantom-atom/file-explorer/internal/stats"
	"github.com/phantom-atom/file-explorer/internal/utils/executor"
	"github.com/phantom-atom/file-explorer/mailer"
	"github.com/phantom-atom/file-explorer/middleware"
	"github.com/phantom-atom/file-explorer/models"
	"github.com/phantom-atom/file-explorer/repository"
	"github.com/phantom-atom/file-explorer/services"
	v1 "github.com/phantom-atom/file-explorer/web/api/v1"
	"github.com/phantom-atom/file-explorer/web/api/v1/register"
	"golang.org/x/crypto/acme/autocert"
)

var (
	globalConfig  *config.Config
	database      *gorm.DB
	workDirectory string
	fileService   *services.FileService
	userService   *services.UserService
	closeExecutor = &executor.Executor{}
	namedLocker   locker.NamedLocker
	asyncMailer   *mailer.Mailer
	dataCache     cache.Cache
	dataContext   repository.DataContext
)

func main() {
	//退出前关闭所有添加的closer
	defer closeExecutor.Execute(false, func(a executor.Action, err error) {
		log.Warn("msg", "occur an error when "+a.Tag()+" close", "error", err.Error())
	})

	//初始化配置文件
	loadConfig()

	//设置全局日志
	log.SetLogger(log.NewZapLogger(globalConfig))

	//初始化邮件模板
	initMailTemplate()

	//初始化全局数据库
	initDatabase()

	//初始化邮件发送器
	initMailer()

	//初始化缓存
	initCache()

	namedLocker = locker.NewGLocker()

	//初始化仓库
	initRepository()

	//初始化全局路径
	initWD()

	//初始化全局服务
	initServices()

	//运行http
	runHTTPServer()
}

func configFunc() *config.Config {
	return globalConfig
}

func initMailer() {
	conf := configFunc()
	asyncMailer = mailer.NewMailer(context.Background(), configFunc, conf.Email.MaxQueueSize)
	closeExecutor.AddFuncWithTag("mailer", asyncMailer.Close)
}

func initMailTemplate() {

	mailTempl := &configFunc().EmailTemplate
	items := mailTempl.Items

	dir, _ := os.Getwd()
	basePath := filepath.Join(dir, mailTempl.BasePath)

	for _, item := range items {
		filename := filepath.Join(basePath, item.Filename)
		if err := mailer.RegisterMailTemplate(
			item.Alias,
			filename,
			item.Subject,
			item.ContentType,
		); err != nil {
			log.Panic("msg", "load mail template file failed", "filename", filename, "error", err.Error())
		}
	}
}

func initCache() {
	conf := &configFunc().Cache
	var err error
	switch conf.Engine {
	case "redis":
		dataCache, err = redis.NewCache(configFunc)
		if err != nil {
			log.Panic("msg", "occur an error when initialize cache", "error", err.Error())
		}
	default:
		log.Panic("msg", "occur an error when initialize cache", "error", "cache engine is unsupported", "engine", conf.Engine)
	}
}

func initRepository() {
	var err error
	dataContext, err = simple.NewContext(
		database, dataCache, configFunc, time.Now,
	)

	if err != nil {
		log.Panic("msg", "occur an error when initializa repository", "error", err.Error())
	}
}

func initDatabase() {
	dbConf := &globalConfig.Database

	dbConnStr := ""
	switch dbConf.Engine {
	case "postgres":
		dbConnStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbConf.Host, dbConf.Port,
			dbConf.User, dbConf.Password,
			dbConf.DBName, dbConf.SSLMode)
	default:
		log.Panic("msg", "occur an error when initialize database", "error", "database engine is unsupported")
	}

	db, err := gorm.Open(dbConf.Engine, dbConnStr)
	if err != nil {
		log.Panic("msg", "occur an error when initialize database", "error", err.Error())
	}

	err = db.AutoMigrate(&models.File{}, &models.User{}).Error
	if err != nil {
		log.Panic("msg", "occur an error when initialize database", "error", err.Error())
	}

	database = db
	closeExecutor.AddFuncWithTag("Database", database.Close)
}

func initWD() {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("msg", "occur an error when get execute directory", "error", err.Error())
		return
	}
	workDirectory = dir
}

func initServices() {
	initFileService()
	initUserService()
}

func initFileService() {
	fs := services.NewFileService(
		configFunc,
		func() string {
			return uuid.New().String()
		},
		dataContext,
		namedLocker,
	)

	fileService = fs
	closeExecutor.AddFuncWithTag("FileService", fs.Close)
}

func initUserService() {
	us := services.NewUserService(
		configFunc,
		func() string {
			return uuid.New().String()
		},
		time.Now,
		dataContext,
		asyncMailer,
		namedLocker,
	)

	userService = us
}

func runHTTPServer() {
	httpConf := &globalConfig.HTTP
	promConf := &globalConfig.Prometheus

	if globalConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	if promConf.Enable {
		promParamFn := middleware.DefaultMetricsParamsFn
		if promConf.GatewayAddr != "" {
			promParamFn = func() (string, int) {
				return promConf.GatewayAddr, promConf.IntervalSec
			}
		}
		prometheusMiddleware := middleware.NewPrometheus("go_prometheus-service", "file-explorer",
			stats.HTTPGather, promParamFn)
		engine.Use(prometheusMiddleware.HandlerFunc())
	}

	apiGroup := engine.Group("/")
	api := v1.NewAPI(configFunc, fileService, userService)

	authorization := middleware.NewAuthorization(userService, configFunc)
	register.APIGINRegister(api, apiGroup, authorization)

	addr := net.JoinHostPort(httpConf.Host, httpConf.Port)
	var err error
	serv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	if httpConf.SSL {
		if httpConf.AutoTLS.Enable {
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      autocert.DirCache(httpConf.AutoTLS.CacheDir),
				HostPolicy: autocert.HostWhitelist(httpConf.AutoTLS.Host...),
			}
			serv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		}
		err = serv.ListenAndServeTLS(httpConf.CertPath, httpConf.KeyPath)
	} else {
		err = serv.ListenAndServe()
	}

	if err != nil {
		log.Error("msg", "occur an error when run http service", "error", err.Error())
	}
}

func loadConfig() {
	conf := config.Load([]string{".", "./config"}, "config", "yaml")
	if conf == nil {
		os.Exit(1)
	}
	globalConfig = conf
}
