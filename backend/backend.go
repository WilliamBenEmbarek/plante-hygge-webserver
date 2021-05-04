package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type dataPacket struct {
	Cmd  string
	EUI  string
	Ts   int
	Ack  bool
	Bat  int
	Fcnt int
	Port int
	Data string
}

type inputDataPacket struct {
	Data string
}

var fakeData = [10]string{
	"123456789123449373849570400000000",
	"12345678912349648i490390399900001",
	"123456789123496775845380400200002",
	"123456789123496894939990403000003",
	"123456789123496759599730410000004",
	"123456789123496704784880400000005",
	"123456789123496774328970300000006",
	"123456789123496704324230200000007",
	"123456789123496708023230100000008",
	"123456789123496709100010400900009"}

var fakeIndex = 0

var dataChannel chan inputDataPacket
var upgrader = websocket.Upgrader{}

func main() {
	dataChannel = make(chan inputDataPacket)
	c := NewController()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.File("/", "public/index.html")
	e.File("/style.css", "public/style.css")
	e.File("/app.js", "public/app.js")
	e.GET("/ws", c.StreamController)
	e.POST("/dataPush", fetchData)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

// Controller interface has two methods
type Controller interface {
	// Homecontroller renders initial home page
	HomeController(e echo.Context) error

	// StreamController responds with live cpu status over websocket
	StreamController(e echo.Context) error
}

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

func (c *controller) HomeController(e echo.Context) error {
	return e.File("views/index.html")
}

func (c *controller) StreamController(e echo.Context) error {
	//TODO Make this only be localhost
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(e.Response(), e.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		newVal := <-dataChannel
		jsonMessage, _ := json.Marshal(newVal.Data)
		err := ws.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func done(control chan bool) {
	control <- true
}

/*func simulateData() {
	setInterval(func() {
		data := dataPacket{
			Cmd:  "Pepe",
			EUI:  "Papa",
			Ts:   1,
			Ack:  false,
			Bat:  1,
			Fcnt: 1,
			Port: 1,
			Data: fakeData[fakeIndex%10],
		}
		dataChannel <- data
		fakeIndex++
	}, 2500)
	fmt.Println("Yeetus Deletus")
}*/
/*
func fetchDataCibicom(c echo.Context) error {
	//fmt.Println(c.Request())
	decoder := json.NewDecoder(c.Request().Body)
	var data dataPacket
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data packet")
	}
	byteArray, err := hex.DecodeString(data.Data)
	if err != nil {
		fmt.Println("Error decoding HEX Data")
	}
	data.Data = string(byteArray)

	fmt.Println(data)
	dataChannel <- data
	return nil
}
*/
func fetchData(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var input inputDataPacket
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("Error decoding data packet")
	}
	fmt.Println(input)
	select {
	case dataChannel <- input:
		fmt.Println("sent to frontend", input)
		return c.JSON(http.StatusOK, input)
	default:
		return c.JSON(http.StatusOK, input)
	}
}

func setInterval(function func(), milliseconds int) chan bool {

	// How often to fire the passed in function in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticker and the channel to signal
	// the ending of 	the interval
	ticker := time.NewTicker(interval)
	control := make(chan bool)

	// Put the selection in a go routine so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-control:
				return
			case <-ticker.C:
				// This won't block
				go function()

			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return control
}
