package server

import (
	"fmt"
	"net/http"
	"strconv"

	"challenge3/cmd/web"
	"challenge3/internal/todoinfo"
	"challenge3/internal/urlinfo"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", echo.WrapHandler(templ.Handler(web.Home())))

	e.GET("/back", echo.WrapHandler(templ.Handler(web.BackHome())))
	e.GET("/back/all", echo.WrapHandler(http.HandlerFunc(web.BackSeeAllHandler)))
	e.GET("/short/:id", echo.WrapHandler(http.HandlerFunc(s.shortRedirect)))

	e.GET("/front", echo.WrapHandler(templ.Handler(web.FrontHome())))

	e.POST("/shorten", echo.WrapHandler(http.HandlerFunc(web.ShortenHandler)))
	e.POST("/front/add", echo.WrapHandler(http.HandlerFunc(web.AddTodoHandler)))

	e.DELETE("/front/delete/:id", s.DeleteTodo)

	return e
}

func (s *Server) shortRedirect(w http.ResponseWriter, r *http.Request) {
	urlinfo.LOCK.Lock()
	defer urlinfo.LOCK.Unlock()

	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Short url is missing", http.StatusBadRequest)
		return
	}

	if fullurl, found := urlinfo.URL_CACHE[urlinfo.ShortUrl(shortKey)]; found {
		http.Redirect(w, r, string(fullurl), http.StatusMovedPermanently)
	} else {
		http.Error(w, "Short url could not be found", http.StatusNotFound)
		return
	}
}

func (s *Server) DeleteTodo(c echo.Context) error {
	urlinfo.LOCK.Lock()
	defer urlinfo.LOCK.Unlock()

	n, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("error parsing id")
	}

	delete(todoinfo.TODOS, n)

	return c.NoContent(200)
}
