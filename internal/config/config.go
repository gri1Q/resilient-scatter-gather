package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config представляет корневую структуру конфигурации приложения,
// соответствующую ожидаемому YAML-файлу (например, config/local.yaml).
// Поля структуры соответствуют секциям в YAML и могут быть заполнены
// через cleanenv.ReadConfig.
type Config struct {
	// Env указывает окружение запуска приложения ("local", "dev", "prod").
	Env string `env:"ENV" env-default:"local" env-required:"true"`

	// StoragePath задаёт путь к локальному хранилищу (каталог для файлов).
	StoragePath string `env:"STORAGE_PATH" env-required:"true"`

	// HTTPServer содержит настройки HTTP-сервера.
	HTTPServer

	// DatabaseServer содержит настройки подключения к базе данных.
	DatabaseServer
}

// HTTPServer описывает конфигурацию HTTP-сервера: адрес и таймауты.
type HTTPServer struct {

	// Addr задаёт адрес и порт сервера (например "localhost:8080" или ":8080").
	Addr string `env:"HTTP_ADDRESS" env-default:"localhost:8080"`

	// Timeout указывает максимальное время ожидания обработки запроса.
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"4s"`

	// IdleTimeout задаёт таймаут простоя соединения.
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

// DatabaseServer описывает параметры подключения к базе данных.
// Используется для конфигурации различных СУБД (PostgreSQL, MySQL, SQLite и др.).
type DatabaseServer struct {
	// Driver определяет тип базы данных, которую использует приложение.
	// Возможные значения: "postgres", "mysql", "sqlite" и т.д.
	// Если приложение поддерживает только одну СУБД (например PostgreSQL),
	// можно задать значение по умолчанию и не менять его.
	Driver string `env:"DB_DRIVER" env-required:"true"`

	// Host задаёт адрес сервера базы данных.
	// В Docker Compose это может быть имя сервиса (например "db"),
	// а при локальном запуске — "localhost" или IP-адрес.
	Host string `env:"DB_HOST" env-required:"true"`

	// Port — порт, на котором слушает сервер базы данных (обычно 5432 для PostgreSQL).
	Port int `env:"DB_PORT" env-required:"true"`

	// User — имя пользователя, под которым приложение подключается к базе данных.
	User string `env:"DB_USER" env-required:"true"`

	// Password — пароль пользователя базы данных.
	Password string `env:"DB_PASSWORD" env-required:"true"`

	// DBName — имя базы данных, к которой нужно подключиться.
	DBName string `env:"DB_NAME" env-required:"true"`

	// SSLMode определяет режим SSL-соединения.
	// Возможные значения: "disable", "require", "verify-ca", "verify-full".
	// Для локальной разработки обычно используется "disable".
	SSLMode string `env:"DB_SSLMODE" env-default:"disable" env-required:"true"`
}

func MustLoad() *Config {
	configPath := ".env"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist:", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config:", err)
	}

	return &cfg
}
