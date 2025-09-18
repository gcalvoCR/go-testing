package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gcalvocr/go-testing/logger"
)

// IndexHandler serves the API documentation page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Serving API documentation page", nil)

	// Get the template file path
	templatePath := filepath.Join("templates", "index.html")

	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Error("Failed to parse template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Error("Failed to execute template", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
