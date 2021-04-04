package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple/handlers"

	gin_gonic "github.com/gin-gonic/gin"
	"github.com/maxstanley/tango/wrapper/gin"
)

func main() {
	// Sets gin to release mode, this remove all debug statements.
	gin_gonic.SetMode(gin_gonic.ReleaseMode)
	// Create a new instance of the gin router.
	r := gin_gonic.New()
	// With recovery middleware.
	r.Use(gin_gonic.Recovery())

	r.POST(gin.Wrapper("/country/:country/city/:city", handlers.NewPostCanDrink))

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	// Starts the HTTP Server in a go routine so the interrupt signals can be
	// handled.
	go func() {
		fmt.Println("Starting HTTP Server.")
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("HTTP Server Error: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	// Wait for a selected signal to interrupt the program.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit
	fmt.Printf("%s Signal has been caught.\n", signal.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the HTTP Server.
	fmt.Println("HTTP Server Shutting down.")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server failed to shutdown gracefully: %s\n", err.Error())
	}
}
