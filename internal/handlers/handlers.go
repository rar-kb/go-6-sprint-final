package handlers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		data []byte
		ext  string
	)

	file, header := pickFile(r)
	if file != nil {
		defer file.Close()
		b, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = b
		ext = filepath.Ext(header.Filename)
	} else {
		// Fallback: возможно, содержимое пришло как обычное поле формы.
		val := pickValue(r)
		if val == "" {
			http.Error(w, "no file or data found in request", http.StatusInternalServerError)
			return
		}
		data = []byte(val)
	}

	converted, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ext == "" {
		ext = ".txt"
	}
	name := time.Now().UTC().Format("2006-01-02_15-04-05.000000000") + ext

	out, err := os.Create(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := out.WriteString(converted); err != nil {
		out.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := out.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(converted))
}

// pickFile возвращает первую файловую часть из формы независимо от имени поля.
func pickFile(r *http.Request) (multipart.File, *multipart.FileHeader) {
	if r.MultipartForm == nil {
		return nil, nil
	}
	preferred := []string{
		"file", "myFile", "upload", "uploadFile", "upload_file",
		"data", "filename", "uploaded_file", "text",
	}
	for _, n := range preferred {
		if fhs := r.MultipartForm.File[n]; len(fhs) > 0 {
			if f, err := fhs[0].Open(); err == nil {
				return f, fhs[0]
			}
		}
	}
	for _, fhs := range r.MultipartForm.File {
		if len(fhs) > 0 {
			if f, err := fhs[0].Open(); err == nil {
				return f, fhs[0]
			}
		}
	}
	return nil, nil
}

// pickValue — резервный вариант: содержимое пришло обычным полем формы.
func pickValue(r *http.Request) string {
	if r.MultipartForm == nil {
		return ""
	}
	preferred := []string{"file", "myFile", "text", "data", "content"}
	for _, n := range preferred {
		if vs := r.MultipartForm.Value[n]; len(vs) > 0 {
			return vs[0]
		}
	}
	for _, vs := range r.MultipartForm.Value {
		if len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}
