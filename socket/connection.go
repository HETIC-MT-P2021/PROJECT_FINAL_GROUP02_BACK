package socket

import (
	"os"
	"log"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}

func Connect() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create an endoint for WebSocket
	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// Upgrade incoming http connection
		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer ws.Close()

		log.Println("Connected!")

		// Listen to a connection
		for {
			var message Message
			err := ws.ReadJSON(&message)
			if !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
				break
			}
			log.Println(message)

			// Send message from server
			if err := ws.WriteJSON(message); !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
			}
		}
		return nil
	})

	apiPort := os.Getenv("DOCKER_API_PORT")
	e.Logger.Fatal(e.Start(`:` + apiPort))
}
