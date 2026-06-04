package handler

import (
	"diplom/internal/auth"
	"diplom/internal/service"

	authmw "diplom/pkg/middleware/auth"

	"github.com/go-chi/chi/v5"
)

// RegisterAPI mounts public auth and JWT-protected routes.
func RegisterAPI(r chi.Router, svc *service.Services, jwt *auth.Manager) {
	authH := NewAuthHandler(svc.Auth)
	me := NewMeHandler(svc)
	users := NewUserHandler(svc.Users)
	subjects := NewSubjectHandler(svc.Subjects)
	topics := NewTopicHandler(svc.Topics)
	tasks := NewTaskHandler(svc.Tasks)
	sessions := NewSessionHandler(svc.Sessions)
	repetitions := NewRepetitionHandler(svc.Repetitions)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authmw.Bearer(jwt))

			r.Route("/me", func(r chi.Router) {
				r.Get("/", me.Get)
				r.Put("/", me.Update)

				r.Get("/subjects", me.ListSubjects)
				r.Post("/subjects", me.CreateSubject)

				r.Get("/sessions", me.ListSessions)
				r.Get("/repetitions", me.ListRepetitions)
				r.Get("/calendar", me.Calendar)

				r.Route("/stats", func(r chi.Router) {
					r.Get("/", me.StatsOverview)
					r.Get("/topics", me.StatsTopics)
					r.Get("/sessions", me.StatsSessions)
					r.Get("/upcoming", me.StatsUpcoming)
				})
			})

			r.Get("/users", users.List)

			r.Route("/subjects", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", subjects.Get)
					r.Put("/", subjects.Update)
					r.Delete("/", subjects.Delete)

					r.Route("/topics", func(r chi.Router) {
						r.Post("/", topics.Create)
						r.Get("/", topics.ListBySubject)
					})
				})
			})

			r.Route("/topics", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", topics.Get)
					r.Put("/", topics.Update)
					r.Delete("/", topics.Delete)

					r.Route("/tasks", func(r chi.Router) {
						r.Post("/", tasks.Create)
						r.Get("/", tasks.ListByTopic)
					})
				})
			})

			r.Route("/tasks", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", tasks.Get)
					r.Put("/", tasks.Update)
					r.Delete("/", tasks.Delete)
				})
			})

			r.Route("/sessions", func(r chi.Router) {
				r.Post("/", sessions.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", sessions.Get)
					r.Post("/finish", sessions.Finish)
					r.Post("/attempts", sessions.SubmitAttempt)
					r.Get("/attempts", sessions.ListAttempts)
					r.Get("/stats", sessions.GetStats)
				})
			})

			r.Route("/repetitions", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Patch("/", repetitions.Update)
				})
			})
		})
	})
}
