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
	"github.com/go-chi/cors"
	planetscale "github.com/harshav17/planet_scale"
	docs "github.com/harshav17/planet_scale/docs"
	utilities "github.com/harshav17/planet_scale/utilites"
	slogchi "github.com/samber/slog-chi"
)

var (
	//go:embed templates/*
	templateFS embed.FS

	//go:embed css/*
	css embed.FS

	//parsed templates
	templates = template.Must(template.ParseFS(templateFS, "templates/*"))
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln         net.Listener
	server     *http.Server
	router     chi.Router
	middleware *Middleware
}

func NewServer(controllers *planetscale.ControllerProvider, middleware *Middleware) *Server {
	s := &Server{
		server:     &http.Server{},
		router:     chi.NewRouter(),
		middleware: middleware,
	}

	logger := utilities.GetLogger()
	s.router.Use(slogchi.New(logger))

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Assuming your CSS file is in a directory named 'css'
	s.router.Handle("/css/output.css", http.FileServer(http.FS(css)))

	// Swagger
	fileServer := http.FileServer(http.FS(docs.Docs))
	s.router.Handle("/swagger/*", http.StripPrefix("/swagger/", fileServer))

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

	s.router.Route("/wh", func(r chi.Router) {
		r.Post("/put_user", controllers.User.HandlePutUser)
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.middleware.withClerkSession())

		r.Route("/groups", func(r chi.Router) {
			r.Get("/", controllers.ExpenseGroup.HandleGetExpenseGroups)
			r.Post("/", controllers.ExpenseGroup.HandlePostExpenseGroup)
			r.Route("/{groupID}", func(r chi.Router) {
				r.Patch("/", controllers.ExpenseGroup.HandlePatchExpenseGroup)
				r.Delete("/", controllers.ExpenseGroup.HandleDeleteExpenseGroup)
				r.Get("/", controllers.ExpenseGroup.HandleGetExpenseGroup)
				r.Route("/members", func(r chi.Router) {
					r.Get("/", controllers.GroupMember.HandleGetGroupMembers)
					r.Post("/", controllers.GroupMember.HandlePostGroupMember)
					r.Delete("/{userID}", controllers.GroupMember.HandleDeleteGroupMember)
				})
				r.Get("/expenses", controllers.Expense.HandleGetGroupExpenses)
				r.Get("/settlements", controllers.Settlement.HandleGetGroupSettlements)
				r.Get("/balances", controllers.ExpenseGroup.HandleGetGroupBalances)
			})
		})

		r.Route("/expenses", func(r chi.Router) {
			r.Post("/", controllers.Expense.HandlePostExpense)
			r.Route("/{expenseID}", func(r chi.Router) {
				r.Patch("/", controllers.Expense.HandlePatchExpense)
				r.Delete("/", controllers.Expense.HandleDeleteExpense)
				r.Get("/", controllers.Expense.HandleGetExpense)
			})
		})

		r.Route("/settlements", func(r chi.Router) {
			r.Post("/", controllers.Settlement.HandlePostSettlement)
			r.Route("/{settlementID}", func(r chi.Router) {
				r.Patch("/", controllers.Settlement.HandlePatchSettlement)
				r.Delete("/", controllers.Settlement.HandleDeleteSettlement)
				r.Get("/", controllers.Settlement.HandleGetSettlement)
			})
		})

		r.Route("/split_types", func(r chi.Router) {
			r.Get("/", controllers.SplitType.HandleGetAllSplitTypes)
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
