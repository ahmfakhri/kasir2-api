package main

import (
	"kasir2-api/database"
	"kasir2-api/handlers"
	"kasir2-api/models"
	"kasir2-api/repositories"
	"kasir2-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var categories = []models.Category{
	{ID: 1, Name: "POP Mie", Description: "Makanan"},
	{ID: 2, Name: "Teh Gelas", Description: "Minuman"},
	{ID: 3, Name: "Susu Indomilk", Description: "Minuman"},
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func maskDBConn(db string) string {
	if db == "" {
		return ""
	}
	// sembunyikan password
	// postgresql://user:PASS@host -> postgresql://user:****@host
	if i := strings.Index(db, "://"); i != -1 {
		rest := db[i+3:]
		if at := strings.Index(rest, "@"); at != -1 {
			userPass := rest[:at]
			if colon := strings.Index(userPass, ":"); colon != -1 {
				return db[:i+3] + userPass[:colon+1] + "****@" + rest[at+1:]
			}
		}
	}
	return db
}

func main() {
	log.Println("Starting Kasir2 API")
	// ===== Load ENV =====
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Println(".env found but failed to read:", err)
		} else {
			log.Println(".env file loaded")
		}
	} else {
		log.Println(".env file not found, using environment variables only")
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// ===== DEBUG ENV =====
	log.Println("🔍 ENV CHECK")
	log.Println("PORT   :", config.Port)
	log.Println("DB_CONN:", maskDBConn(config.DBConn))
	log.Println("======================")

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
	row := db.QueryRow("SELECT current_database()")
	var dbName string
	row.Scan(&dbName)
	log.Println(" Database connected")
	log.Println(" Database name:", dbName)
	defer db.Close()

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
