package socketio

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"strconv"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/SteakBarbare/RPGBot/utils"
)

var (
	router *mux.Router
	Server *gosocketio.Server
)

type Message struct {
	Text string `json:"message"`
}

func init() {
	Server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	fmt.Println("Socket Inititalize...")
}

func LoadSocket() {
	// Socket connection
	Server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		fmt.Println("Connected", c.Id())

		c.Emit("/welcome", "")
	})

	Server.On("/joinRoom", func(c *gosocketio.Channel, message Message) string {
		dungeonId := message.Text
		c.Join("Room" + dungeonId)

		return "Join Room successfully"
	})

	// Socket disconnection
	Server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		fmt.Println("Disconnected", c.Id())

		c.Emit("/disconnected", "")
	})

	Server.On("/leaveRoom", func(c *gosocketio.Channel, message Message) string {
		dungeonId := message.Text
		c.Leave("Room" + dungeonId)

		return "Leave Room successfully"
	})

	// Dungeon socket
	Server.On("/getDungeonSquares", func(c *gosocketio.Channel, message Message) string {
		dungeonId, err := strconv.Atoi(message.Text)
		if err != nil {
			log.Println(err)
		}

		dungeonTiles, err := utils.GetFullDungeonTiles(dungeonId)
		if err != nil {
			log.Println(err)
		}

		c.BroadcastTo("Room" + message.Text, "/dungeonSquares", dungeonTiles)

		return "message sent successfully."
	})
}

func UpdateDungeonTiles(dungeonId int) {
	s := strconv.Itoa(dungeonId)
	Server.BroadcastTo("Room" + s, "/updateDungeonSquares", "")
}

func CreateRouter() {
	router = mux.NewRouter()
}

func InititalizeRoutes() {
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
	router.Handle("/socket.io/", Server)
}

func StartServer() {
	apiPort := os.Getenv("API_PORT")
	fmt.Println("Socketio server started at http://localhost:" + apiPort)
	log.Fatal(http.ListenAndServe(`:` + apiPort, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

func Connect() {
	LoadSocket()
	CreateRouter()
	InititalizeRoutes()
	StartServer()
}
