package app

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rajenderK7/goly/internal/encode"
)

type ReqBody struct {
	LongURL  string `json:"long_url,omitempty" firestore:"long_url,omitempty"`
	ShortURL string `json:"short_url,omitempty" firestore:"short_url,omitempty"`
}

type ErrMessage struct {
	Message string `json:"message"`
}

func (app *App) getURLHelper(key, val string, ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return app.db.Collection("urls").
		Where(key, "==", val).
		Documents(ctx).
		GetAll()
}

func (app *App) shortURLHelper(longURL string, ctx context.Context) (string, error) {
	var (
		shortURL = ""
		docs     []*firestore.DocumentSnapshot
		err      error
	)

	// Keep generating short URLs until there are no collisions
	for {
		shortURL = encode.GenerateShortURL(longURL)
		docs, err = app.getURLHelper("short_url", shortURL, ctx)
		if err != nil {
			return "", err
		}
		if len(docs) == 0 {
			break
		}
	}

	return shortURL, nil
}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to Goly!")
}

// Desc: Create and/or return shortened URL
// Route: POST api/create
func (app *App) createShortURL(c echo.Context) error {
	var req, res ReqBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrMessage{"bad request"})
	}
	ctx := c.Request().Context()
	docs, err := app.getURLHelper("long_url", req.LongURL, ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrMessage{"internal server error"})
	}

	// Return the short URL if already exists
	if len(docs) > 0 {
		docs[0].DataTo(&res)
		return c.JSON(http.StatusOK, res)
	}

	shortURL, err := app.shortURLHelper(req.LongURL, ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrMessage{"internal server error"})
	}

	res = ReqBody{
		LongURL:  req.LongURL,
		ShortURL: shortURL,
	}

	_, _, err = app.db.Collection("urls").Add(ctx, res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrMessage{"internal server error"})

	}
	return c.JSON(http.StatusOK, res)
}

// Desc: Redirects to the original URL if the short URL is valid
// Route: GET api/short-url
func (app *App) redirectURL(c echo.Context) error {
	var req, res ReqBody
	req.ShortURL = c.Param("short-url")

	ctx := c.Request().Context()
	docs, err := app.getURLHelper("short_url", req.ShortURL, ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrMessage{"internal server error"})
	}
	if len(docs) == 0 {
		return c.JSON(http.StatusBadRequest, ErrMessage{"Invalid URL"})
	}

	docs[0].DataTo(&res)
	return c.Redirect(http.StatusMovedPermanently, res.LongURL)
}
