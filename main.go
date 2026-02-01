package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"kasir2-api/database"
	"kasir2-api/handlers"
	"kasir2-api/repositories"
	"kasir2-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// ===== Load ENV =====
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DBConn == "" {
		log.Fatal("DB_CONN tidak ditemukan")
	}

	// ===== Connect Database =====
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("gagal koneksi database:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// ===== Dependency Injection =====
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// ===== Routes =====
	// categories === products table
	http.HandleFunc("/api/categories", productHandler.HandleProducts)
	http.HandleFunc("/api/categories/", productHandler.HandleProductByID)

	// optional health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// ===== Start Server =====
	log.Println("Server running on port", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
}
