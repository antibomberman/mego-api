package api

import (
	localMiddleware "github.com/antibomberman/mego-api/internal/adapters/api/middleware"
	"github.com/antibomberman/mego-api/internal/clients"
	"github.com/antibomberman/mego-api/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type Server struct {
	router        chi.Router
	authClient    *clients.AuthClient
	userClient    *clients.UserClient
	postClient    *clients.PostClient
	storageClient *clients.StorageClient
	cfg           *config.Config
}

func NewServer(cfg *config.Config, authClient *clients.AuthClient, userClient *clients.UserClient, postClient *clients.PostClient, storageClient *clients.StorageClient) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)
	s := &Server{
		router:        r,
		authClient:    authClient,
		userClient:    userClient,
		postClient:    postClient,
		storageClient: storageClient,
		cfg:           cfg,
	}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login/send_code", s.AuthLoginSendCode)
		r.Post("/login", s.AuthLogin)

	})
	r.Route("/user", func(r chi.Router) {
		r.Get("/show/{id}", s.UserShow)

		r.Group(func(r chi.Router) {
			r.Use(localMiddleware.JwtMiddleware)
			r.Get("/me", s.UserMe)
			r.Post("/update/profile", s.UserUpdateProfile)
			r.Post("/update/theme", s.UserUpdateTheme)
			r.Post("/update/lang", s.UserUpdateLang)
			r.Post("/update/email/send_code", s.UserUpdateEmailSendCode)
			r.Post("/update/email", s.UserUpdateEmail)

		})
	})
	r.Route("/post", func(r chi.Router) {
		r.Get("/show/{id}", s.PostShow)
		r.Get("/", s.PostList)
		r.Group(func(r chi.Router) {
			r.Use(localMiddleware.JwtMiddleware)
			r.Post("/", s.PostCreate)
			r.Get("/my", s.PostMyList)
			r.Get("/show/{id}", s.PostShow)
			r.Put("/{id}", s.PostUpdate)
			r.Delete("/{id}", s.PostDelete)
		})

	})

	return s
}

func (s *Server) Start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
