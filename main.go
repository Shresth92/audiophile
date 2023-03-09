package main

import (
	"github.com/Shresth92/audiophile/api/controller"
	"github.com/Shresth92/audiophile/api/middlewares"
	"github.com/Shresth92/audiophile/api/routes"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/services"
	"github.com/Shresth92/audiophile/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func startServer(
	db *internal.Database,
	router *internal.RequestHandler,
	route *routes.Routes,
	lifecycle fx.Lifecycle) {
	route.Setup()

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router.Gin,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			db.MigrateUpDb()
			models.FirebaseClient, _ = utils.GetFirebaseClient()
			go func(srv *http.Server) {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.Error(err)
				}
			}(srv)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if dbCloseErr := db.CloseDb(); dbCloseErr != nil {
				return dbCloseErr
			}
			if serverCloseErr := srv.Shutdown(ctx); serverCloseErr != nil {
				return serverCloseErr
			}
			return nil
		},
	})
}

func main() {
	var CommonModules = fx.Options(
		controller.Module,
		routes.Module,
		services.Module,
		internal.Module,
		middlewares.Module,
	)
	app := fx.New(CommonModules, fx.Invoke(startServer))
	if app.Err() != nil {
		logrus.Error("failed to start the app")
	}
	app.Run()
}
