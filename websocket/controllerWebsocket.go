package websocket

import (
	"encoding/json"
	"github.com/Joan-Goma/websocketTemplate/controller"
	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader websocket.Upgrader

var upgrader_default = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SetUpgrader(upgrader_modified websocket.Upgrader) {
	upgrader = upgrader_modified
}

func StartWebSocketServer(context *gin.Context) {

	//upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		engine.Warning.Println(err)
		return
	}
	err = ws.SetReadDeadline(time.Now().Add(45 * time.Minute))
	if err != nil {
		engine.Warning.Println("Could not set dead time out for the new client", err)
		return
	}
	defer func(ws *websocket.Conn) {

	}(ws)

	engine.Debug.Println("New client connected!")

	c := controller.GenerateClient(ws, ws.RemoteAddr().String())

	if err != nil {
		err := ws.Close()
		if err != nil {
			engine.Warning.Println("Can not close this connection")
			return
		}
		engine.Warning.Println("error generating new client", err)
	}

	for {

		err = ReadMessage(ws, &c.IncomingMessage)

		if err != nil {
			engine.Warning.Println(err)
			c.LastMessage.Command = "invalid_message"
			c.LastMessage.Data["error"] = "server could not read correctly the message received, please try again"
			c.SendMessage()
			return
		}

		c.LastMessage.RequestID = c.IncomingMessage.RequestID
		c.LastMessage.Command = c.IncomingMessage.Command

		//Execute the command readed
		c.ExecuteCommand(c.IncomingMessage.Command)
		c.MessageReader <- c.IncomingMessage

		//TODO fix this
		//if c.User.Banned {
		//	err := c.WS.Close()
		//	if err != nil {
		//		engine.Debug.Println("could not close connection by user banned")
		//		return
		//	}
		//	return
		//}
		//Reset the data of the incomming message
		data := make(map[string]interface{})
		cM := controller.Message{
			Data: data,
		}
		dataTwo := make(map[string]interface{})
		cS := controller.Message{
			Data: dataTwo,
		}
		c.IncomingMessage = cM
		c.LastMessage = cS
		engine.Debug.Println("command processed")
	}
}

func ReadMessage(ws *websocket.Conn, dest interface{}) error {
	//	engine.Debug.Println("New Incoming message")
	//engine.Debug.Printf("Client, %d sent another request %d", message.RequestID, message.RequestID)
	//Read Message from client
	_, m, err := ws.ReadMessage()
	if err != nil {
		engine.Warning.Println(err)
		return err
	}

	err = json.Unmarshal(m, &dest)
	if err != nil {
		engine.Warning.Println(err)
		return err
	}
	return nil
}
