package websocket

import (
	"os"
	"log"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/SteakBarbare/RPGBot/game"
)

var upgrader = websocket.Upgrader{}

type NotifsReceived struct {
	Message string
	Event string
}
type NotifisToSend struct {
	Message []game.DungeonTile
	Event string
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
			var notifsReceived NotifsReceived
			err := ws.ReadJSON(&notifsReceived)
			if !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
				break
			}
			log.Println(notifsReceived)

			// Get dungeon id from client side
			v := notifsReceived
			var message = getField(&v, "Message");
			dungeonId, err := strconv.Atoi(message)
			if err != nil {
				log.Println(err)
			}

			dungeonTiles, err := utils.GetFullDungeonTiles(dungeonId)
			if err != nil {
				log.Println(err)
			}
			
			// Send DungeonSquares from server
			notifsToSent := NotifisToSend{
				Message: dungeonTiles,
				Event: "EventDungeonSquares",
			}
			if err := ws.WriteJSON(notifsToSent); !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
			}
		}
		return nil
	})

	apiPort := os.Getenv("DOCKER_API_PORT")
	e.Logger.Fatal(e.Start(`:` + apiPort))
}

func getField(v *NotifsReceived, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
