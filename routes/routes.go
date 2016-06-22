package routes

import (
	"github.com/jinzhu/gorm"

	"github.com/hyperionii/medivac/config"
	"github.com/hyperionii/medivac/handlers/auth"
)

// Routes contains all routes for the application.
type Routes struct {
	routes []Route
}

// NewRoutes creates a new Router instance and initializes all routes.
func NewRoutes(cfg *config.Config, db *gorm.DB) *Routes {
	authHandler := auth.NewHandler(cfg)

	r := &Routes{
		routes: []Route{
			&route{
				pattern:       "auth/login/",
				method:        "POST",
				handlerFunc:   authHandler.Login,
				requiresAuth:  false,
				requiredRoles: []string{},
			},
		},
	}

	return r
}
