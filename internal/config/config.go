package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" enc-required:"./data"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Слово must говорит о том, что функция не будет возвращать ошибку
// А будет просто паниковать, если что-то пойдет не так
// Этим подходом нельзя злоупотреблять!
// Тут мы загружаем конфиг и заполняем его структуру
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// Получает значение конфига
// Либо из переменной окружения, либо из флага
// Приоритет: флаг > переменной окружения
func fetchConfigPath() string {
	var res string

	// --config="path/to/congig.yaml"
	flag.StringVar(&res, "config", "", "path to congig file")
	flag.Parse()

	// Если не был использован флаг, то
	// Достаем путь из переменной окружения
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
