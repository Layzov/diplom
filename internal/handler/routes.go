package handler

import (
	"diplom/internal/service"

	"github.com/go-chi/chi/v5"
)

// RegisterAPI mounts CRUD routes under /api/v1.
func RegisterAPI(r chi.Router, svc *service.Services) {
	users := NewUserHandler(svc.Users)
	subjects := NewSubjectHandler(svc.Subjects)
	topics := NewTopicHandler(svc.Topics)
	tasks := NewTaskHandler(svc.Tasks)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", users.Create)
			r.Get("/", users.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", users.Get)
				r.Put("/", users.Update)
				r.Delete("/", users.Delete)

				r.Route("/subjects", func(r chi.Router) {
					r.Post("/", subjects.Create)
					r.Get("/", subjects.ListByUser)
				})
			})
		})

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
	})
}
