package controllers

import (
	"sort"

	"main/models"
)

func groupByCategory(wishes []models.Wish) []Category {
	categories := make(map[string][]models.Wish)
	for _, wish := range wishes {
		categories[wish.Category] = append(categories[wish.Category], wish)
	}

	var result []Category
	for name, wishes := range categories {
		result = append(result, Category{Name: name, Wishes: wishes})
	}
	return result
}

func sortCategories(categories []Category) []Category {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	return categories
}

func getOptions(values []string) []Option {
	var options []Option
	for _, value := range values {
		options = append(options, Option{Label: value, Value: value})
	}
	return options
}
