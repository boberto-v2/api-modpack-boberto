package event_service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

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
	eventId := ctx.Query("id")
	if err != nil {
		wsSession.Close()
		log.Println(err)
	}
	if eventId == "" {
		wsSession.Close()
		return
	}
	uid := uuid.New()
	log.Printf("A client connect to %s", eventId)
	if _, ok := sessionGroupMap[eventId]; ok {
		sessionGroupMap[eventId][uid] = wsSession
	} else {
		sessionGroupMap[eventId] = make(map[uuid.UUID]*websocket.Conn)
		sessionGroupMap[eventId][uid] = wsSession
	}
	defer wsSession.Close()
	// _, found := GetById(eventId)
	// if !found {
	// 	wsSession.WriteMessage(1, []byte("Event not found or expired"))
	// 	wsSession.Close()
	// }
	echo(wsSession, eventId, uid)
}

func echo(wsSession *websocket.Conn, eventId string, uid uuid.UUID) {
	for { //An endlessloop
		messageType, messageContent, err := wsSession.ReadMessage()
		if messageType == 1 {
			log.Printf("Recv:%s from %s", messageContent, eventId)
			emit(eventId, "FREEEEEMAN!!!")
		}
		if err != nil {
			wsSession.Close()
			delete(sessionGroupMap[eventId], uid)
			//I don't think it's recommended to deal with connection closing like this, but it's the easiest way.
			//Or you have to maintain a hashmap to indicate if a session is open or closed? No idea.
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Client disconnected in %s", eventId)
			} else {
				log.Printf("Reading Error in %s. %s", eventId, err)
			}
			break //To escape from the endless loop
		}
	}
}

func (event Event) Emit(messageContent string) {
	emit(event.Id, messageContent)
}

func emit(eventId, messageContent string) {
	log.Printf("emit event to %s %s", eventId, messageContent)
	for _, wsSession := range sessionGroupMap[eventId] {
		err := wsSession.WriteJSON(map[string]any{
			"event": eventId,
			"data":  messageContent,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func Complete() {

}
