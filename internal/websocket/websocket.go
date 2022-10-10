package websocket

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/carfloresf/financial-chat/internal/constants"
	"github.com/carfloresf/financial-chat/internal/queue"
)

func getName(session *melody.Session) string {
	userSession, ok := session.Keys[sessions.DefaultKey].(sessions.Session)
	if !ok {
		log.Errorf("error getting user session")

		return ""
	}

	return cast.ToString(userSession.Get(constants.Userkey))
}

func Handler(websocketServer *melody.Melody, queueClient queue.Queuer) func(c *gin.Context) {
	websocketServer.HandleDisconnect(func(session *melody.Session) {
		melodySession := Session{session}
		userName := getName(session)

		// broadcast <user> left channel message
		err := websocketServer.BroadcastFilter(
			[]byte(fmt.Sprintf("%s left channel",
				userName)),
			melodySession.SameChannel)
		if err != nil {
			log.Errorf("error broadcasting to others: %s", err)
		}
	})

	websocketServer.HandleConnect(func(session *melody.Session) {
		melodySession := Session{session}
		userName := getName(session)

		// broadcast <user> joined channel message
		err := websocketServer.BroadcastFilter(
			[]byte(fmt.Sprintf("%s joined channel",
				userName)),
			melodySession.SameChannel)
		if err != nil {
			log.Errorf("error broadcasting to others: %s", err)
		}
	})

	websocketServer.HandleMessage(func(session *melody.Session, msg []byte) {
		splitMessage := strings.SplitAfter(string(msg), "> ")
		melodySession := Session{session}

		// broadcast original command
		err := websocketServer.BroadcastFilter(msg, melodySession.SameChannel)
		if err != nil {
			log.Errorf("error broadcasting: %v", err)
		}

		if strings.HasPrefix(splitMessage[1], "/") {
			correlationID := uuid.New().String()
			queueClient.StoreSession(correlationID, session)

			// send command to queue for processing
			err := queueClient.Publish([]byte(splitMessage[1]), correlationID, "request")
			if err != nil {
				log.Errorf("error publishing to queue: %s", err)
			}

			msg = []byte("<system> command forwarded: " + splitMessage[1])

			// broadcast "command forwarded" message
			err = websocketServer.BroadcastFilter(msg, melodySession.SameChannel)
			if err != nil {
				log.Errorf("error broadcasting: %v", err)
			}
		}
	})

	return func(c *gin.Context) {
		// handle request with keys to be able to access user session
		err := websocketServer.HandleRequestWithKeys(c.Writer, c.Request, c.Keys)
		if err != nil {
			log.Errorf("error handling request: %s", err)
		}
	}
}

type Session struct {
	*melody.Session
}

// SameChannel to filter sessions by channel
func (s *Session) SameChannel(other *melody.Session) bool {
	return s.Request.URL.Path == other.Request.URL.Path
}
