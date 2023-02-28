package main

import (
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/server"
	"github.com/Shresth92/audiophile/utils"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

func main() {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	err := utils.LoadEnv()
	if err != nil {
		logrus.Errorf("Environment variables loading failed.; %s", err.Error())
		return
	}

	models.FirebaseClient, err = utils.GetFirebaseClient()
	if err != nil {
		logrus.Errorf("Firebase client createtion failed; %s", err.Error())
		return
	}

	port := utils.GetEnvValue("PORT")
	server := server.ServeStart()
	database.ConnectDb()

	go func() {
		if err := server.Start(port); err != nil {
			logrus.Errorf("Server not shut down gracefully; %s", err.Error())
			return
		}
	}()

	<-done
	database.CloseDb()
	if err := server.Stop(5 * time.Second); err != nil {
		logrus.Errorf("Server not shut down gracefully; %s", err.Error())
		return
	}
}
