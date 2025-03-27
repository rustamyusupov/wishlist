package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var templates map[string]*template.Template

func Initialize() error {
	templates = make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"formatPrice": FormatPrice,
	}

	pages := []string{"home", "new", "edit"}
	for _, page := range pages {
		tmpl, err := template.New("layout.tmpl").Funcs(funcMap).ParseFiles(
			filepath.Join("web", "templates", "layout.tmpl"),
			filepath.Join("web", "templates", fmt.Sprintf("%s.tmpl", page)),
		)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", page, err)
		}
		templates[page] = tmpl
	}

	return nil
}

func FormatPrice(price float64) string {
	s := fmt.Sprintf("%.2f", price)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	decimalPart := parts[1]

	var result string
	for i := len(intPart); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}

		if len(result) > 0 {
			result = intPart[start:i] + " " + result
		} else {
			result = intPart[start:i]
		}
	}

	return result + "," + decimalPart
}

func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, ok := templates[tmplName]
	if !ok {
		http.Error(w, fmt.Sprintf("Template %s not found", tmplName), http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmplName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func SendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func SendError(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := Response{
		Success: false,
		Error:   errorMsg,
	}
	SendJSON(w, response, statusCode)
}

func SendSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	response := Response{
		Success: true,
		Data:    data,
	}
	SendJSON(w, response, statusCode)
}

func HandleError(w http.ResponseWriter, err error, status int, logMessage string) {
	if logMessage != "" {
		log.Printf("%s: %v", logMessage, err)
	}
	http.Error(w, err.Error(), status)
}
