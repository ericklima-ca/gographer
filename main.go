package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericklima-ca/gographer/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()
	oauth := router.Group("oauth")
	{
		oauth.GET("/microsoft", routers.MicrosoftLogin)
	}

	callback := router.Group("callback")
	{
		callback.GET("/microsoft", routers.MicrosoftCallback)
	}

	httpSrv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil &&
			errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error on %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting...")
}
