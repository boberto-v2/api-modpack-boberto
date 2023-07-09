package event_service

//thanks to https://gist.github.com/crosstyan/47e7d3fa1b9e4716c0d6c76760a4a70c
import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var sessionGroupMap = make(map[string]map[uuid.UUID]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(ctx *gin.Context) {
	wsSession, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	uid := uuid.New()
	wsURL := ctx.Request.URL
	wsURLParam, err := url.ParseQuery(wsURL.RawQuery)
	if err != nil {
		wsSession.Close()
		log.Println(err)
	}
	if _, ok := wsURLParam["name"]; ok {
		threadName := wsURLParam["name"][0]
		log.Printf("A client connect to %s", threadName)
		if _, ok := sessionGroupMap[threadName]; ok {
			sessionGroupMap[threadName][uid] = wsSession
		} else {
			sessionGroupMap[threadName] = make(map[uuid.UUID]*websocket.Conn)
			sessionGroupMap[threadName][uid] = wsSession
		}
		defer wsSession.Close()
		echo(wsSession, threadName, uid)
		return
	}
	wsSession.Close()
}

func echo(wsSession *websocket.Conn, threadName string, uid uuid.UUID) {
	//Message Type:
	//Details in
	//https://godoc.org/github.com/gorilla/websocket#pkg-constants
	//TextMessage=1
	//BinaryMessage=2
	for { //An endlessloop
		messageType, messageContent, err := wsSession.ReadMessage()
		if messageType == 1 {
			log.Printf("Recv:%s from %s", messageContent, threadName)
			Emit(threadName, messageContent)
		}
		if err != nil {
			wsSession.Close()
			delete(sessionGroupMap[threadName], uid)
			//I don't think it's recommended to deal with connection closing like this, but it's the easiest way.
			//Or you have to maintain a hashmap to indicate if a session is open or closed? No idea.
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Client disconnected in %s", threadName)
			} else {
				log.Printf("Reading Error in %s. %s", threadName, err)
			}
			break //To escape from the endless loop
		}
	}
}
func Emit(threadName string, messageContent []byte) {
	for _, wsSession := range sessionGroupMap[threadName] {
		err := wsSession.WriteMessage(1, messageContent)
		if err != nil {
			log.Println(err)
		}
	}
}
