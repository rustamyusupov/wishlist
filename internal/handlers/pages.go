package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/rustamyusupov/wishes/internal/database"
	"github.com/rustamyusupov/wishes/internal/models"
)

type DisplayCategory struct {
	Name   string
	Wishes []models.Wish
}

func Home(w http.ResponseWriter, r *http.Request) {
	wishes, err := database.GetWishes()
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get wishes")
		return
	}

	categories := organizeWishesByCategory(wishes)
	RenderTemplate(w, "home", struct{ Categories []DisplayCategory }{Categories: categories})
}

func New(w http.ResponseWriter, r *http.Request) {
	currencyOptions, err := getCurrencyOptions()
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get currencies")
		return
	}

	categoryOptions, err := getCategoryOptions()
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get categories")
		return
	}

	data := struct {
		Currencies []models.Option
		Categories []models.Option
	}{
		Currencies: currencyOptions,
		Categories: categoryOptions,
	}

	RenderTemplate(w, "new", data)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		HandleError(w, fmt.Errorf("id is required"), http.StatusBadRequest, "Missing wish ID")
		return
	}

	wish, err := database.GetWishByID(id)
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get wish")
		return
	}

	currencyOptions, err := getCurrencyOptions()
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get currencies")
		return
	}

	categoryOptions, err := getCategoryOptions()
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError, "Failed to get categories")
		return
	}

	data := struct {
		Wish       models.Wish
		Currencies []models.Option
		Categories []models.Option
	}{
		Wish:       wish,
		Currencies: currencyOptions,
		Categories: categoryOptions,
	}

	RenderTemplate(w, "edit", data)
}

func organizeWishesByCategory(wishes []models.Wish) []DisplayCategory {
	categoryMap := make(map[string][]models.Wish)
	for _, wish := range wishes {
		categoryMap[wish.Category] = append(categoryMap[wish.Category], wish)
	}

	var displayCategories []DisplayCategory
	for name, categoryWishes := range categoryMap {
		displayCategories = append(displayCategories, DisplayCategory{
			Name:   name,
			Wishes: categoryWishes,
		})
	}

	sort.Slice(displayCategories, func(i, j int) bool {
		return displayCategories[i].Name < displayCategories[j].Name
	})

	return displayCategories
}

func getCategoryOptions() ([]models.Option, error) {
	categories, err := database.GetCategories()
	if err != nil {
		return nil, err
	}

	options := make([]models.Option, len(categories))
	for i, category := range categories {
		options[i] = models.Option{
			Label: category.Name,
			Value: category.Name,
		}
	}

	return options, nil
}

func getCurrencyOptions() ([]models.Option, error) {
	currencies, err := database.GetCurrencies()
	if err != nil {
		return nil, err
	}

	options := make([]models.Option, len(currencies))
	for i, currency := range currencies {
		options[i] = models.Option{
			Label: currency.Code,
			Value: currency.Code,
		}
	}

	return options, nil
}
