package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "POP Mie", Description: "Makanan"},
	{ID: 2, Name: "Teh Gelas", Description: "Minuman"},
	{ID: 3, Name: "Susu Indomilk", Description: "Minuman"},
}

// GET /api/categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// PUT /api/categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updated Category
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			updated.ID = id
			categories[i] = updated

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// DELETE /api/categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category berhasil dihapus",
			})
			return
		}
	}
	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

func main() {

	// GET /api/categories/{id}
	// PUT /api/categories/{id}
	// DELETE /api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// GET /api/categories
	// POST /api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)

		} else if r.Method == "POST" {
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	fmt.Println("Server running di localhost:8080")
	http.ListenAndServe(":8080", nil)
}
