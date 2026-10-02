package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-default:"configs/game.db"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type GameConfig struct {
	Categories []struct {
		Name      string `yaml:"name"`
		Questions []struct {
			Question  string `yaml:"question"`
			Answer    string `yaml:"answer"`
			Points    int    `yaml:"points"`
			MediaType string `yaml:"media_type"`
			MediaURL  string `yaml:"media_url"`
		} `yaml:"questions"`
	} `yaml:"categories"`
}

type GameConfigManager struct {
	Config GameConfig
}

func NewGameConfigManager(configPath string) (*GameConfigManager, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config GameConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &GameConfigManager{
		Config: config,
	}, nil
}

func MustLoad() *Config {
	// تجاهل خطأ عدم وجود ملف .env على السحاب
	_ = godotenv.Load()

	// إذا لم يتم تحديد CONFIG_PATH، استخدم المسار الافتراضي محلياً أو على السحاب
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/local.yaml" // أو مسار ملف الإعدادات الخاص بك
	}

	var cfg Config

	// إذا وجد ملف الإعدادات نقرأه، وإذا لم يوجد نعتمد القيم الافتراضية لتجنب الانهيار
	if _, err := os.Stat(configPath); err == nil {
		_ = cleanenv.ReadConfig(configPath, &cfg)
	}

	// دعم منفذ Render التلقائي
	if port := os.Getenv("PORT"); port != "" {
		cfg.HttpServer.Address = ":" + port
	}

	if cfg.HttpServer.Address == "" {
		cfg.HttpServer.Address = ":8080"
	}

	if cfg.StoragePath == "" {
		cfg.StoragePath = "configs/game.db"
	}

	return &cfg
}