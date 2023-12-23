package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rajenderK7/goly/internal/app"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	PROD := os.Getenv("PROD")
	var opt option.ClientOption
	if PROD == "true" {
		// This is specific to the service on which the backend (server) is hosted
		opt = option.WithCredentialsFile("serviceAccountKey.json")
	} else {
		opt = option.WithCredentialsFile("../../config/firebase-admin/serviceAccountKey.json")
	}
	fbApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := fbApp.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	app := app.NewApp(e, db)
	app.Init()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3001" // fallback to port :3001 in dev mode
	}
	app.Start(PORT, ctx)
}
