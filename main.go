package main

/*
#cgo LDFLAGS: -L./lib -lhnsw
*/
import "C"
import (
	"context"
	"crypto/rand"
	"eigen_db/api"
	"eigen_db/cfg"
	"eigen_db/constants"
	"eigen_db/vector_io"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func displayAsciiArt() { // because it looks cool
	fmt.Println(`
███████╗██╗ ██████╗ ███████╗███╗   ██╗██████╗ ██████╗ 
██╔════╝██║██╔════╝ ██╔════╝████╗  ██║██╔══██╗██╔══██╗
█████╗  ██║██║  ███╗█████╗  ██╔██╗ ██║██║  ██║██████╔╝
██╔══╝  ██║██║   ██║██╔══╝  ██║╚██╗██║██║  ██║██╔══██╗
███████╗██║╚██████╔╝███████╗██║ ╚████║██████╔╝██████╔╝
╚══════╝╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚═════╝ 			
	`)
}

func setupRedisClient(ctx context.Context, host string, port int, pass string, apiKey string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       0,
	})

	if _, err := client.Get(ctx, constants.REDIS_API_KEY_NAME).Result(); err != nil { // key does not exist
		keyBytes := make([]byte, 32)
		_, err := rand.Read(keyBytes)
		if err != nil {
			panic(err)
		}
		apiKey = hex.EncodeToString(keyBytes)
	}

	if apiKey != "" {
		status := client.Set(ctx, "apiKey", apiKey, 0)
		if status.Err() != nil {
			panic(status.Err())
		}
	}

	val, err := client.Get(ctx, constants.REDIS_API_KEY_NAME).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("EIGENDB API KEY: %s\n", val)

	return client
}

func main() {
	displayAsciiArt()

	var apiKey string
	var redisHost string
	var redisPort int
	var redisPass string

	flag.StringVar(&apiKey, "api-key", "", "EigenDB API key")
	flag.StringVar(&redisHost, "redis-host", "127.0.0.1", "Redis server host IP")
	flag.IntVar(&redisPort, "redis-port", 6379, "Redis server host port")
	flag.StringVar(&redisPass, "redis-pass", "", "Redis server password (default \"\")")
	flag.Parse()

	cfg.NewConfig()                              // creates a empty Config struct in memory
	config := (&cfg.ConfigFactory{}).GetConfig() // get pointer to Config in memory
	config.LoadConfig(constants.CONFIG_PATH)     // load config from config.yml into the Config struct in memory

	vector_io.SetupDB(config)

	ctx := context.Background()
	redisClient := setupRedisClient(ctx, redisHost, redisPort, redisPass, apiKey)

	if err := vector_io.StartPersistenceLoop(config); err != nil {
		panic(err)
	}

	if err := api.StartAPI(ctx, fmt.Sprintf("%s:%d", config.GetAPIAddress(), config.GetAPIPort()), redisClient); err != nil {
		panic(err)
	}
}
