package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rustamyusupov/wishes/internal/database"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		SendError(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	link := r.FormValue("link")
	priceStr := r.FormValue("price")
	currencyCode := r.FormValue("currency")
	categoryName := r.FormValue("category")

	if name == "" || link == "" || priceStr == "" || currencyCode == "" || categoryName == "" {
		SendError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		SendError(w, "Invalid price value", http.StatusBadRequest)
		return
	}

	category, err := database.GetCategoryByName(categoryName)
	if err != nil {
		log.Printf("Error getting category: %v", err)
		SendError(w, "Invalid category", http.StatusBadRequest)
		return
	}

	currency, err := database.GetCurrencyByCode(currencyCode)
	if err != nil {
		log.Printf("Error getting currency: %v", err)
		SendError(w, "Invalid currency", http.StatusBadRequest)
		return
	}

	_, err = database.CreateWish(link, name, category.ID, price, currency.ID)
	if err != nil {
		log.Printf("Error adding wish: %v", err)
		SendError(w, "Failed to add wish", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Patch(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		SendError(w, "Missing wish ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		SendError(w, "Invalid wish ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		SendError(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	link := r.FormValue("link")
	priceStr := r.FormValue("price")
	currencyCode := r.FormValue("currency")
	categoryName := r.FormValue("category")

	if name == "" || link == "" || priceStr == "" || currencyCode == "" || categoryName == "" {
		SendError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	category, err := database.GetCategoryByName(categoryName)
	if err != nil {
		log.Printf("Error getting category: %v", err)
		SendError(w, "Invalid category", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		SendError(w, "Invalid price value", http.StatusBadRequest)
		return
	}

	currency, err := database.GetCurrencyByCode(currencyCode)
	if err != nil {
		log.Printf("Error getting currency: %v", err)
		SendError(w, "Invalid currency", http.StatusBadRequest)
		return
	}

	if err := database.UpdateWish(id, link, name, category.ID); err != nil {
		log.Printf("Error updating wish: %v", err)
		SendError(w, "Failed to update wish", http.StatusInternalServerError)
		return
	}

	latestPrice, err := database.GetLatestPrice(id)
	if err != nil || latestPrice.Amount != price || latestPrice.CurrencyID != currency.ID {
		_, err = database.CreatePrice(id, price, currency.ID)
		if err != nil {
			log.Printf("Error adding price: %v", err)
			SendError(w, "Failed to update price", http.StatusInternalServerError)
			return
		}
	}

	SendSuccess(w, map[string]interface{}{
		"message": fmt.Sprintf("Wish with ID %d updated successfully", id),
		"id":      id,
	}, http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		SendError(w, "Missing wish ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		SendError(w, "Invalid wish ID", http.StatusBadRequest)
		return
	}

	if err := database.DeletePricesByWishID(id); err != nil {
		log.Printf("Error deleting prices for wish: %v", err)
	}

	if err := database.DeleteWish(id); err != nil {
		log.Printf("Error deleting wish: %v", err)
		SendError(w, "Failed to delete wish", http.StatusInternalServerError)
		return
	}

	SendSuccess(w, map[string]interface{}{
		"message": fmt.Sprintf("Wish with ID %d deleted successfully", id),
		"id":      id,
	}, http.StatusOK)
}
