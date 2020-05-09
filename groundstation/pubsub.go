package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/crmathieu/daq/data"
//	"gopkg.in/yaml.v2"
//	"io/ioutil"
	"os"
//	"path"
//	"runtime"
	"time"
	"strings"
	b64 "encoding/base64"
)
/*
type EgoConfig struct {
	Redis ERedis `yaml:"redis"`
}

type ERedis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"pass"`
}

var Config struct {
	Ego EgoConfig `yaml:"ego"`
}

type RedisINFO struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"pass"`
}
*/

// default values...
const (
	DEF_REDIS_HOST = "127.0.0.1"
	DEF_REDIS_PORT = "6379"
	DEF_REDIS_PWD  = ""

	PUBSUB_CHANNEL = "DAQChannel"
)

//var Rclient *redis.Client

// InitConfig------------------------------------------------------------------
func InitConfig() bool {

	var env, redisenv string
	var ok bool

	if env, ok = os.LookupEnv("DAQ_ENV"); !ok {
		fmt.Printf("Fatal: environment variable DAQ_ENV not found...")
		return false
	}
/*
	if dbenv, ok = os.LookupEnv("DAQ_DB_DSN_"+env); ok {
		dec, _ := b64.StdEncoding.DecodeString(dbenv)
		data.Config.DB_dsn = string(dec)
	} else {
		fmt.Printf("Fatal: environment variable DAQ_DB_DSN_"+env+" not found...")
		return false
	}
*/	
	if redisenv, ok = os.LookupEnv("DAQ_REDIS_DSN_"+env); ok {
		dec, _ := b64.StdEncoding.DecodeString(redisenv)
		data.CInfo.REDIS_dsn = string(dec)
		params := strings.Split(data.CInfo.REDIS_dsn, ":")
/*		data.Config.Cache_Redis_host = params[0]
		data.Config.Cache_Redis_port = params[1]
		data.Config.Cache_Redis_pass = params[2]
*/
		data.CInfo.RedisHost = params[0]
		data.CInfo.RedisPort = params[1]
		data.CInfo.RedisPass = params[2]
	} else {
		data.CInfo.RedisHost = DEF_REDIS_HOST
		data.CInfo.RedisPort = DEF_REDIS_PORT
		data.CInfo.RedisPass = DEF_REDIS_PWD
		fmt.Printf("Fatal: environment variable DAQ_REDIS_DSN_"+env+" not found...")
	}
	return RedisInit()

/*	fmt.Println("\nInitializing DB...")
	data.Config.DB = psql.DBinit(data.Config.DB_dsn)
	if data.Config.DB == nil {
		fmt.Println(data.Config)
		return false
	}
*/
}


func RedisInit() bool {
	const REDIS_RETRIES = 5
	var err error
	Rclient = redis.NewClient(&redis.Options{
		Addr:     data.CInfo.RedisHost + ":" + data.CInfo.RedisPort,
		Password: data.CInfo.RedisPass,
		DB:       0,
	})
	for i:=0; i<REDIS_RETRIES; i++ {
		_, err = Rclient.Ping().Result()
		if err != nil {
			fmt.Printf(".")
		} else {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func Publish(Payload string) error {
	fmt.Printf("PUBLISHING: %s\n", Payload)
	return Rclient.Publish(PUBSUB_CHANNEL, Payload).Err()
}
