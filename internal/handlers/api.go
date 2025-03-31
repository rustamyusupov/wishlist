package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rustamyusupov/wishes/internal/database"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error parsing form")
		return
	}

	name := r.FormValue("name")
	link := r.FormValue("link")
	priceStr := r.FormValue("price")
	currencyCode := r.FormValue("currency")
	categoryName := r.FormValue("category")
	sortStr := r.FormValue("sort")

	if name == "" || link == "" || priceStr == "" || currencyCode == "" || categoryName == "" {
		HandleError(w, fmt.Errorf("all fields are required"), http.StatusBadRequest, "")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		HandleError(w, fmt.Errorf("invalid price value"), http.StatusBadRequest, "")
		return
	}

	sort := 0
	if sortStr != "" {
		sort, err = strconv.Atoi(sortStr)
		if err != nil {
			HandleError(w, fmt.Errorf("invalid sort value"), http.StatusBadRequest, "")
			return
		}
	}

	category, err := database.GetCategoryByName(categoryName)
	if err != nil {
		HandleError(w, err, http.StatusBadRequest, "Error getting category")
		return
	}

	currency, err := database.GetCurrencyByCode(currencyCode)
	if err != nil {
		HandleError(w, err, http.StatusBadRequest, "Error getting currency")
		return
	}

	_, err = database.CreateWish(link, name, category.ID, price, currency.ID, sort)
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error adding wish")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Patch(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		HandleError(w, fmt.Errorf("missing wish ID"), http.StatusBadRequest, "")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		HandleError(w, fmt.Errorf("invalid wish ID"), http.StatusBadRequest, "")
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error parsing form")
		return
	}

	name := r.FormValue("name")
	link := r.FormValue("link")
	priceStr := r.FormValue("price")
	currencyCode := r.FormValue("currency")
	categoryName := r.FormValue("category")
	sortStr := r.FormValue("sort")

	if name == "" || link == "" || priceStr == "" || currencyCode == "" || categoryName == "" {
		HandleError(w, fmt.Errorf("all fields are required"), http.StatusBadRequest, "")
		return
	}

	category, err := database.GetCategoryByName(categoryName)
	if err != nil {
		HandleError(w, err, http.StatusBadRequest, "Error getting category")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		HandleError(w, fmt.Errorf("invalid price value"), http.StatusBadRequest, "")
		return
	}

	sort := 0
	if sortStr != "" {
		sort, err = strconv.Atoi(sortStr)
		if err != nil {
			HandleError(w, fmt.Errorf("invalid sort value"), http.StatusBadRequest, "")
			return
		}
	}

	currency, err := database.GetCurrencyByCode(currencyCode)
	if err != nil {
		HandleError(w, err, http.StatusBadRequest, "Error getting currency")
		return
	}

	if err := database.UpdateWish(id, link, name, category.ID, sort); err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error updating wish")
		return
	}

	latestPrice, err := database.GetLatestPrice(id)
	if err != nil || latestPrice.Amount != price || latestPrice.CurrencyID != currency.ID {
		_, err = database.CreatePrice(id, price, currency.ID)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError, "Error adding price")
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
		HandleError(w, fmt.Errorf("missing wish ID"), http.StatusBadRequest, "")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		HandleError(w, fmt.Errorf("invalid wish ID"), http.StatusBadRequest, "")
		return
	}

	if err := database.DeletePricesByWishID(id); err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error deleting prices for wish")
		return
	}

	if err := database.DeleteWish(id); err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Error deleting wish")
		return
	}

	SendSuccess(w, map[string]interface{}{
		"message": fmt.Sprintf("Wish with ID %d deleted successfully", id),
		"id":      id,
	}, http.StatusOK)
}
