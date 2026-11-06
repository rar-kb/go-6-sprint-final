package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const indexFile = "index.html"

// Index отдаёт HTML-форму из файла index.html.
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, indexFile)
}

// Upload принимает файл из формы, конвертирует его содержимое
func Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := safeName(time.Now().UTC().String()) + filepath.Ext(header.Filename)

	out, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := out.WriteString(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(result)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// safeName заменяет символы, недопустимые в имени файла
func safeName(name string) string {
	replacer := strings.NewReplacer(
		":", "-",
		" ", "_",
		".", "_",
		"/", "-",
		"\\", "-",
	)
	return replacer.Replace(name)
}
