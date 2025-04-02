package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"challenge4/cmd/web"
	"challenge4/internal/notes"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	USER_NAME = "benji"
	PASSWORD  = "pass1"

	API_USAGE = `Usage

All routes are relative to '/api/v1'

METHOD		ROUTE					RESPONSE CODE		RESPONSE BODY
GET			/						200					this message
GET			/notes					200					list of notes
GET			/notes/ids				200					list of note ids
GET			/notes/ids/[id]			200|404				note with that id | note with id not found
GET			/notes/titles			200					list of note titles
GET			/notes/titles/[title]	200|404				note with that title | note with title not found

POST		/						405					this message
POST		/notes					201|4xx|5xx			created note with that id | error message
POST		/notes/ids				405					this message
POST		/notes/ids/[id]			405					this message
POST		/notes/titles			405					this message
POST		/notes/titles/[title]	405					this message

PUT			/						405					this message
PUT			/notes					405					this message
PUT			/notes/ids				405					this message
PUT			/notes/ids/[id]			200|201|4xx|5xx		updated|created note with that id | error message
PUT			/notes/titles			405					this message
PUT			/notes/titles/[title]	200|201|4xx|5xx		updated|created note with that title | error message

PATCH		/*						405					this message

DELETE		/						405					this message
DELETE		/notes					405					this message
DELETE		/notes/ids				405					this message
DELETE		/notes/ids/[id]			204|4xx|5xx			no content | error message
DELETE		/notes/titles			405					this message
DELETE		/notes/titles/[title]	204|4xx|5xx			no content | error message
`
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == USER_NAME && password == PASSWORD, nil
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.FS(web.Files))

	apiGroup := e.Group("/api/v1")
	apiGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == USER_NAME && password == PASSWORD, nil
	}))

	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.GET("/health", s.healthHandler)
	e.GET("/", echo.WrapHandler(templ.Handler(web.Home())))

	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	// api
	apiGroup.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Quicknotes!\n\n"+API_USAGE)
	})
	apiGroup.GET("/notes", s.getNotesHandler)
	apiGroup.GET("/notes/:route", s.getNotesRouteHandler)

	apiGroup.POST("/*", s.methodNotAllowedHandler)
	apiGroup.POST("/notes", s.postNotesHandler)

	apiGroup.PUT("/*", s.methodNotAllowedHandler)
	apiGroup.PUT("/notes/ids/:id", s.putNotesIDHandler)
	apiGroup.PUT("/notes/titles/:title", s.putNotesTitleHandler)

	apiGroup.PATCH("/*", s.methodNotAllowedHandler)

	apiGroup.DELETE("/*", s.methodNotAllowedHandler)
	apiGroup.DELETE("/notes/ids/:id", s.deleteNotesIDHandler)
	apiGroup.DELETE("/notes/titles/:title", s.deleteNotesTitleHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) methodNotAllowedHandler(c echo.Context) error {
	ret := fmt.Sprintf("[%s] %s is not a valid request\n\n%s", c.Request().Method, c.Path(), API_USAGE)
	return c.String(http.StatusMethodNotAllowed, ret)
}

func (s *Server) getNotesHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	var sb strings.Builder
	sb.WriteString("NOTES\n-----\n\n")
	for _, note := range notes.NOTES {
		sb.WriteString(fmt.Sprintf("# %s\n\n", note))
	}
	return c.String(http.StatusOK, sb.String())
}

func (s *Server) getNotesRouteHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	var sb strings.Builder
	parts := strings.Split(c.Param("route"), "/")
	route, id := parts[0], ""
	if len(parts) == 2 {
		id = parts[1]
	}

	if route == "titles" {
		if id == "" {
			sb.WriteString("NOTES TITLES\n------------\n\n")
			for _, v := range notes.NOTES {
				sb.WriteString(fmt.Sprintf("%s\n", v.Title))
			}
			return c.String(http.StatusOK, sb.String())
		} else {
			for _, note := range notes.NOTES {
				if note.Title == id {
					return c.String(http.StatusOK, fmt.Sprintf("# %s\n\n", note))
				}
			}
			return c.String(http.StatusNotFound, fmt.Sprintf("note with title [%s] doesn't exist", id))
		}
	} else if route == "ids" {
		if id == "" {
			sb.WriteString("NOTES IDS\n---------\n\n")
			for k := range notes.NOTES {
				sb.WriteString(fmt.Sprintf("%d\n", k))
			}
			return c.String(http.StatusOK, sb.String())
		} else {
			n, err := strconv.Atoi(id)
			if err != nil {
				return c.String(http.StatusBadRequest, fmt.Sprintf("[%s] is an invalid id", id))
			}

			if note, ok := notes.NOTES[n]; ok {
				ret := fmt.Sprintf("# %s", note)
				return c.String(http.StatusOK, ret)
			} else {
				return c.String(http.StatusNotFound, fmt.Sprintf("note with id [%d] doesn't exist", n))
			}
		}
	} else {
		return s.methodNotAllowedHandler(c)
	}
}

func (s *Server) postNotesHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	note := new(notes.Note)
	note.Id = notes.CUR_ID
	if err := c.Bind(note); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	notes.NOTES[note.Id] = note
	notes.CUR_ID++
	return c.String(http.StatusCreated, fmt.Sprintf("# %s", note))
}

func (s *Server) putNotesIDHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("[%s] is an invalid id", strId))
	}

	if note, found := notes.NOTES[id]; found {
		note.Contents = c.FormValue("contents")
		return c.String(http.StatusOK, fmt.Sprintf("# %s", note))
	}

	note := &notes.Note{Id: notes.CUR_ID, Title: c.FormValue("title"), Contents: c.FormValue("contents")}
	notes.NOTES[notes.CUR_ID] = note
	notes.CUR_ID++

	return c.String(http.StatusCreated, fmt.Sprintf("# %s", note))
}

func (s *Server) putNotesTitleHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	title := c.Param("title")
	for _, note := range notes.NOTES {
		if note.Title == title {
			note.Contents = c.FormValue("contents")
			return c.String(http.StatusOK, fmt.Sprintf("# %s", note))
		}
	}

	note := &notes.Note{Id: notes.CUR_ID, Title: c.FormValue("title"), Contents: c.FormValue("contents")}
	notes.NOTES[note.Id] = note
	notes.CUR_ID++

	return c.String(http.StatusCreated, fmt.Sprintf("# %s", note))
}

func (s *Server) deleteNotesIDHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("[%s] is an invalid id", strId))
	}

	if _, found := notes.NOTES[id]; found {
		delete(notes.NOTES, id)
		return c.NoContent(http.StatusNoContent)
	}

	return c.String(http.StatusNotFound, fmt.Sprintf("note with id [%d] could not be found", id))
}

func (s *Server) deleteNotesTitleHandler(c echo.Context) error {
	notes.LOCK.Lock()
	defer notes.LOCK.Unlock()

	title := c.Param("title")
	for _, note := range notes.NOTES {
		if note.Title == title {
			delete(notes.NOTES, note.Id)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.String(http.StatusNotFound, fmt.Sprintf("note with title [%s] could not be found", title))
}
