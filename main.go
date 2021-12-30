package main

/**

 */

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const VERSION = "v1.0.1"

var (
	configFile  string
	showVersion bool
	config      *Config
	cache       *freecache.Cache
	creator     *Creator
	counter     *Counter
)

type Config struct {
	ServerConfig
	CacheConfig
}

type CacheConfig struct {
	Size      int
	GCPercent int
	Expire    int
}

type ServerConfig struct {
	Listen  string
	Token   string
	Creator int
	Pool    int
}

func init() {

	flag.StringVar(&configFile, "c", "default.conf", "configuration file.")
	flag.BoolVar(&showVersion, "v", false, "show current version.")
	flag.Parse()
	//Show Version
	if showVersion {
		fmt.Printf("Current version: %s", VERSION)
		os.Exit(0)
	}
	//Default Value
	config = &Config{
		ServerConfig{
			Listen:  "0.0.0.0:20724",
			Token:   "brant",
			Creator: 10,
			Pool:    200,
		},
		CacheConfig{
			Size:      1024,
			GCPercent: 20,
			Expire:    300,
		},
	}
	//Default Value
	err := ini.MapTo(&config, configFile)
	if err != nil {
		fmt.Printf("Fail to read configuration file: %v", err)
		os.Exit(1)
	}
	//Cache Initialization
	cache = freecache.NewCache(config.Size * 1024 * 1024)
	creator = NewCreator(config.Creator, config.Pool)
	counter = NewCounter()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

}

func main() {

	r := gin.Default()
	// 计数器
	r.Use(func(c *gin.Context) {
		defer c.Next()
		if strings.HasPrefix(c.Request.URL.Path, "/lookup") {
			counter.OnLookup()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/create") {
			counter.OnCreate()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/notify") {
			counter.OnSubmit()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/submit") {
			counter.OnSubmit()
			return
		}
		// 跨域
	}, func(c *gin.Context) {
		defer c.Next()
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	})

	// 存活探测
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// 度量
	r.GET("/actuator/:metric", func(c *gin.Context) {
		var value int64 = 0
		metric := c.Param("metric")
		if "cache" == metric {
			c.JSON(http.StatusOK, gin.H{
				"r":              0,
				"EntryCount":     cache.EntryCount(),
				"LookupCount":    cache.LookupCount(),
				"HitCount":       cache.HitCount(),
				"HitRate":        cache.HitRate(),
				"MissCount":      cache.MissCount(),
				"OverwriteCount": cache.OverwriteCount(),
				"EvacuateCount":  cache.EvacuateCount(),
				"ExpiredCount":   cache.ExpiredCount(),
			})
			return
		}
		switch metric {
		case "create":
			value = counter.GetCreate()
		case "lookup":
			value = counter.GetLookup()
		case "submit":
			value = counter.GetSubmit()
		case "notify":
			value = counter.GetNotify()
		default:
			value = -1
		}
		if value == -1 {
			c.JSON(http.StatusOK, gin.H{
				"r": -1,
				"m": "allowed metric : cache, create, lookup, submit, notify",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"r": 0,
				"v": value,
			})
		}
	})
	// 创建二维码
	r.GET("/create", func(c *gin.Context) {
		q := creator.Get()
		err := cache.Set([]byte(q.uuid), nil, config.Expire)
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -1,
				"m": "error",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"r": 0,
				"u": q.uuid,
				"i": q.image,
			})
		}
	})
	//轮询二维码
	r.GET("/lookup/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		//uuid不合法
		if len(uuid) != 36 {
			c.JSON(http.StatusOK, gin.H{
				"r": -1,
				"m": "uuid illegal",
			})
			return
		}
		val, err := cache.Get([]byte(uuid))
		//缓存获取错误
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -2,
				"m": "uuid not found in cache",
			})
			return
		}
		//还未屏幕扫码
		if nil == val {
			c.JSON(http.StatusOK, gin.H{
				"r": 1,
				"m": "uuid found, but value is empty",
			})
			return
		}
		//已经屏幕扫码
		if "brant" == string(val) {
			c.JSON(http.StatusOK, gin.H{
				"r": 2,
				"m": "notified",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"r": 0,
				"v": val,
			})
		}
	})

	//提交扫码结果
	r.POST("/submit/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		// uuid不合法
		if len(uuid) != 36 {
			c.JSON(http.StatusOK, gin.H{
				"r": -1,
				"m": "uuid illegal",
			})
			return
		}
		// uuid过期
		_, err := cache.Get([]byte(uuid))
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -2,
				"m": "uuid not found",
			})
			return
		}
		formData := c.PostForm("data")
		// 扫码为空
		if len(formData) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"r": -3,
				"m": "data is empty",
			})
			return
		}
		// 扫码结果解b64
		data, err := base64.StdEncoding.DecodeString(formData)
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -4,
				"m": "b64 decode failed",
			})
			return
		}
		// 扫码结果写缓存
		err = cache.Set([]byte(uuid), data, config.Expire)
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -5,
				"m": "data is error",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"r": 0,
				"m": "submit success",
			})
		}
	})
	//扫码触发服务端
	r.POST("/notify/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		//uuid不合法
		if len(uuid) != 36 {
			c.JSON(http.StatusOK, gin.H{
				"r": -1,
				"m": "uuid illegal",
			})
			return
		}
		_, err := cache.Get([]byte(uuid))
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -2,
				"m": "uuid not found",
			})
			return
		}
		err = cache.Set([]byte(uuid), []byte("brant"), config.Expire)
		if nil != err {
			c.JSON(http.StatusOK, gin.H{
				"r": -3,
				"m": "data is error",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"r": 0,
				"m": "notify success",
			})
		}
	})
	log.Println("Server is listen on", config.Listen)
	err := r.Run(config.Listen)
	if err != nil {
		fmt.Printf("Start Server Failed, : %v", err)
	}

}
