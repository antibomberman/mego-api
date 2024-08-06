package api

import (
	"github.com/antibomberman/mego-api/internal/clients"
	"github.com/antibomberman/mego-api/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type Server struct {
	router        chi.Router
	userClient    *clients.UserClient
	postClient    *clients.PostClient
	storageClient *clients.StorageClient
	cfg           *config.Config
}

func NewServer(cfg *config.Config,
	userClient *clients.UserClient,
	postClient *clients.PostClient,
	storageClient *clients.StorageClient) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)

	s := &Server{
		router:        r,
		userClient:    userClient,
		postClient:    postClient,
		storageClient: storageClient,
		cfg:           cfg,
	}

	r.Route("/user", func(r chi.Router) {
		r.Post("/login/send_code", s.UserLoginSendCode)
		r.Post("/login/", s.UserLogin)
		r.Get("/show/{id}", s.UserShow)
		r.Put("/{id}", s.UserUpdate)
	})
	r.Route("/post", func(r chi.Router) {
		r.Post("/", s.PostCreate)
		r.Get("/show/{id}", s.PostShow)
		r.Get("/", s.PostList)
		r.Put("/{id}", s.PostUpdate)
		r.Delete("/{id}", s.PostDelete)
	})

	return s
}

func (s *Server) Start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
