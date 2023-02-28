package models

import (
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"context"
)

var FirebaseClient *App

type App struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}
