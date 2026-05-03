package config

		import (
			"fmt"
			"log"
			"os"

			"github.com/joho/godotenv"
		)

		type Config struct {
			Port       string
			DBDialect  string
			DBHost     string
			DBPort     string
			DBUser     string
			DBPassword string
			DBName     string
			DBStorage  string
			JWTSecret  string
		}

		var Defaults = map[string]string{
			"PORT":        "8080",
			"DB_DIALECT":  "sqlite",
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASSWORD": "password",
			"DB_NAME":     "database",
			"DB_STORAGE":  "./database.sqlite",
			"JWT_SECRET":  "changeme",
		}

		func (c *Config) Initialize() {
			env := os.Getenv("APP_ENV")
			if env == "" {
				env = "local"
			}

			var envFile string
			switch env {
			case "local":
				envFile = ".env.local"
			case "development":
				envFile = ".env.development"
			case "test":
				envFile = ".env.test"
			default:
				envFile = ".env"
			}

			err := godotenv.Load(envFile)
			if err != nil {
				log.Printf("Warning: Could not load %s file: %v", envFile, err)
				if envFile != ".env" {
					err = godotenv.Load(".env")
					if err != nil {
						log.Printf("Warning: Could not load .env file: %v", err)
					}
				}
			}

			c.Port = c.Load("PORT")
			c.DBDialect = c.Load("DB_DIALECT")
			c.DBHost = c.Load("DB_HOST")
			c.DBPort = c.Load("DB_PORT")
			c.DBUser = c.Load("DB_USER")
			c.DBPassword = c.Load("DB_PASSWORD")
			c.DBName = c.Load("DB_NAME")
			c.DBStorage = c.Load("DB_STORAGE")

			if c.DBDialect == "sqlite" {
				if c.DBName == "database" && c.DBStorage != "" {
					c.DBName = c.DBStorage
				}
				if c.DBName == "" || c.DBName == "database" {
					c.DBName = "./database.sqlite"
				}
			}
			c.JWTSecret = c.Load("JWT_SECRET")

			fmt.Printf("Environment: %s\n", env)
		}

		func (c *Config) Load(e string) string {
			var data string
			data = os.Getenv(e)
			if data == "" {
				data = Defaults[e]
			}
			return data
		}
		