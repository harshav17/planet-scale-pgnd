package http

import (
	"context"
	"embed"
	"log/slog"
	"net"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	planetscale "github.com/harshav17/planet_scale"
	utilities "github.com/harshav17/planet_scale/utilites"
	slogchi "github.com/samber/slog-chi"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/*
	css embed.FS

	//parsed templates
	templates = template.Must(template.ParseFS(templateFS, "templates/*"))
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router chi.Router
}

func NewServer(controllers *planetscale.ControllerProvider) *Server {
	s := &Server{
		server: &http.Server{},
		router: chi.NewRouter(),
	}

	logger := utilities.GetLogger()
	s.router.Use(slogchi.New(logger))

	// Assuming your CSS file is in a directory named 'css'
	s.router.Handle("/css/output.css", http.FileServer(http.FS(css)))

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			Error(w, r, err)
			return
		}
	})

	s.router.Route("/products", func(r chi.Router) {
		r.Route("/{productID}", func(r chi.Router) {
			r.Get("/", controllers.Product.HandleGetProduct)
		})
		r.Post("/", controllers.Product.HandlePostProduct)
		r.Get("/add", controllers.Product.HandleProductAdd)
	})

	s.router.Route("/groups", func(r chi.Router) {
		r.Get("/", controllers.ExpenseGroup.HandleGetExpenseGroups)
		r.Post("/", controllers.ExpenseGroup.HandlePostExpenseGroup)
		r.Route("/{groupID}", func(r chi.Router) {
			r.Patch("/", controllers.ExpenseGroup.HandlePatchExpenseGroup)
			r.Delete("/", controllers.ExpenseGroup.HandleDeleteExpenseGroup)
			r.Get("/", controllers.ExpenseGroup.HandleGetExpenseGroup)
		})
	})

	s.server.Handler = s.router
	return s
}

func (s *Server) Open() error {
	var err error
	if s.ln, err = net.Listen("tcp", ":8080"); err != nil {
		return err
	}
	slog.Info("Listening on :8080")
	go s.server.Serve(s.ln)

	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
