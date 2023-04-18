package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gopher-bell/site/internal/models"
	"github.com/gopher-bell/site/log"
	"go.uber.org/zap"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		log.ZapLogger.Error("failed to get snippet lists", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/pages/home.html",
	// 	"./ui/html/partials/nav.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	log.ZapLogger.Error("failed to parse template files", zap.Error(err))
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	log.ZapLogger.Error("failed to execute template files", zap.Error(err))
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// w.Write([]byte("Hello from Jong!"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			log.ZapLogger.Error(models.ErrNoRecord.Error(), zap.Error(err))
			http.NotFound(w, r)
		} else {
			log.ZapLogger.Error("failed to get data", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		log.ZapLogger.Error("failed to create snippet", zap.Error(err))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
