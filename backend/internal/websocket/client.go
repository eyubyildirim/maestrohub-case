package websocket

import (
	"encoding/json"
	"log"
	"time"
	"tower-defence-engine/internal/game"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Command struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type CreateSquadPayload struct {
	UnitIDs []string `json:"unit_ids"`
}

type DeployIndividualsPayload struct {
	UnitIDs []string `json:"unit_ids"`
}

type DeploySquadPayload struct {
	SquadID string `json:"squad_id"`
}

func (c *Client) readPump(commandCh chan<- game.ICommand) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var cmd Command
		if err := json.Unmarshal(message, &cmd); err != nil {
			log.Printf("error unmarshalling command: %v", err)
			continue
		}

		switch cmd.Action {
		case "create_squad":
			var payload CreateSquadPayload
			if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
				log.Printf("error unmarshalling create_group payload: %v", err)
				continue
			}

			commandCh <- game.CreateSquadCommand{UnitIDs: payload.UnitIDs}
		case "deploy_individuals":
			var payload DeployIndividualsPayload
			if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
				log.Printf("error unmarshalling send_individual payload: %v", err)
				continue
			}

			commandCh <- game.DeployIndividualsCommand{UnitIDs: payload.UnitIDs}
		case "deploy_squad":
			var payload DeploySquadPayload
			if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
				log.Printf("error unmarshalling deploy_squad payload: %v", err)
				continue
			}
			commandCh <- game.DeploySquadCommand{SquadID: payload.SquadID}
		default:
			log.Printf("received unknown command action: %s", cmd.Action)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
