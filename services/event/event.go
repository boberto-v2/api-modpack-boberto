package event_service

//thanks to https://gist.github.com/crosstyan/47e7d3fa1b9e4716c0d6c76760a4a70c
import (
	"log"
	"net/http"
	"net/url"

	event_cache "github.com/brutalzinn/boberto-modpack-api/services/event/cache"
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

func WebSocketHandler(ctx *gin.Context) {
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
		eventName := wsURLParam["name"][0]
		log.Printf("A client connect to %s", eventName)
		if _, ok := sessionGroupMap[eventName]; ok {
			sessionGroupMap[eventName][uid] = wsSession
		} else {
			sessionGroupMap[eventName] = make(map[uuid.UUID]*websocket.Conn)
			sessionGroupMap[eventName][uid] = wsSession
		}
		defer wsSession.Close()
		_, found := event_cache.GetById(eventName)
		if !found {
			Emit(eventName, []byte("Event not found or expired"))
			wsSession.Close()
		}
		echo(wsSession, eventName, uid)
		return
	}
	wsSession.Close()
}

func echo(wsSession *websocket.Conn, eventName string, uid uuid.UUID) {
	//Message Type:
	//Details in
	//https://godoc.org/github.com/gorilla/websocket#pkg-constants
	//TextMessage=1
	//BinaryMessage=2
	for { //An endlessloop
		messageType, messageContent, err := wsSession.ReadMessage()
		if messageType == 1 {
			log.Printf("Recv:%s from %s", messageContent, eventName)
			Emit(eventName, messageContent)
		}
		if err != nil {
			wsSession.Close()
			delete(sessionGroupMap[eventName], uid)
			//I don't think it's recommended to deal with connection closing like this, but it's the easiest way.
			//Or you have to maintain a hashmap to indicate if a session is open or closed? No idea.
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Client disconnected in %s", eventName)
			} else {
				log.Printf("Reading Error in %s. %s", eventName, err)
			}
			break //To escape from the endless loop
		}
	}
}
func Emit(eventName string, messageContent []byte) {
	for _, wsSession := range sessionGroupMap[eventName] {
		err := wsSession.WriteMessage(1, messageContent)
		if err != nil {
			log.Println(err)
		}
	}
}
