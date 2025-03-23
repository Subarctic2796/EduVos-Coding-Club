package web

import (
	"challenge3/internal/todoinfo"
	"challenge3/internal/urlinfo"
	"fmt"
	"log"
	"net/http"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	urlinfo.LOCK.Lock()
	defer urlinfo.LOCK.Unlock()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	fullurl := r.FormValue("url")
	shorturl := Shorten(fullurl)
	if _, found := urlinfo.URL_CACHE[urlinfo.ShortUrl(shorturl)]; !found {
		urlinfo.URL_CACHE[urlinfo.ShortUrl(shorturl)] = urlinfo.FullUrl(fullurl)
	}
	component := ShortenPost(fmt.Sprintf("/short/%s", shorturl))
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in ShortenHandler: %e", err)
	}
}

func Shorten(s string) string {
	var hash uint = 5381
	for _, i := range s {
		hash = ((hash << 5) + hash) + uint(i)
	}
	return fmt.Sprintf("%d", hash)[:15]
}

func BackSeeAllHandler(w http.ResponseWriter, r *http.Request) {
	urlinfo.LOCK.Lock()
	defer urlinfo.LOCK.Unlock()

	strmap := make(map[string]string)
	for k, v := range urlinfo.URL_CACHE {
		strmap[fmt.Sprintf("/short/%s", k)] = string(v)
	}

	component := BackSeeAll(strmap)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in BackSeeAllHandler: %e", err)
	}
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoinfo.LOCK.Lock()
	defer todoinfo.LOCK.Unlock()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	component := Item(name, fmt.Sprintf("%d", todoinfo.CUR_ID+1))
	todoinfo.CUR_ID++
	todoinfo.TODOS[todoinfo.CUR_ID] = todoinfo.NewTodo(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in AddTodoHandler: %e", err)
	}
}
