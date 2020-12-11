package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/websocket/v2"
	"github.com/lighthouse-p2p/hub/internal/models"
	"github.com/lighthouse-p2p/hub/internal/utils"
	"golang.org/x/crypto/nacl/sign"
)

const (
	initState          = 0
	codeSentState      = 1
	authenticatedState = 2
)

// Signaling handles all the webrtc signaling and real-time comms
func (h *Handlers) Signaling(c *websocket.Conn) {
	// To start signaling, the peer first needs to authenticate
	// On connection, the server sends a random string
	// It expects to recv a signature signed with the peer's private key
	// The signature is validated, and authentication succeeds
	// A redis pubsub connection is made
	// The client gets all the signaling data via the "sub"
	// It can publish the signaling data via "pub"
	// Channels used for pubsub are "signaling_{base64_of_the_public_key}"

	socketState := initState

	// pubKey is the base64 of the public key, same as stored in the DB
	pubKey := c.Query("pub_key")
	if pubKey == "" {
		c.Close()

		return
	}

	var user models.DHT
	tx := h.Cfg.Database.Model(&models.DHT{}).Where("pub_key = ?", pubKey).First(&user)
	if tx.Error != nil {
		c.Close()

		return
	}

	challange, err := utils.GenerateRandomStringURLSafe(32)
	if err != nil {
		log.Printf("%s\n", err)
		c.Close()

		return
	}

	pubKeyBytesSlice, err := base64.StdEncoding.DecodeString(user.PubKey)
	if err != nil {
		log.Printf("%s\n", err)
		c.Close()

		return
	}

	var pubKeyBytes [32]byte
	copy(pubKeyBytes[:], pubKeyBytesSlice)

	c.WriteMessage(1, []byte(challange))
	socketState = codeSentState

	var (
		// mt is the message type
		// noFrame = -1, TextMessage = 1, BinaryMessage = 2
		msg []byte
	)

	var pubsub *redis.PubSub

	for {
		if _, msg, err = c.ReadMessage(); err != nil {
			break
		}

		switch socketState {
		case codeSentState:
			_, valid := sign.Open(nil, msg, &pubKeyBytes)
			if !valid {
				c.Close()

				return
			}

			pubsub = h.Cfg.Redis.Subscribe(context.Background(), getRedisPSKey(pubKey))
			go func() {
				for msg := range pubsub.Channel() {
					c.WriteMessage(1, []byte(msg.Payload))
				}
			}()

			defer pubsub.Unsubscribe(context.Background(), getRedisPSKey(pubKey))

			c.WriteMessage(1, []byte("OK"))

			if err != nil {
				log.Fatalln(err)
			}

			socketState = authenticatedState

			break

		case authenticatedState:
			if _, msg, err = c.ReadMessage(); err != nil {
				break
			}

			var signal models.Signal
			err = json.Unmarshal(msg, &signal)
			if err != nil || signal.To == "" {
				continue
			}

			channel := getRedisPSKey(signal.To)
			h.Cfg.Redis.Publish(context.Background(), channel, string(msg))

			break

		default:
			break
		}
	}
}

func getRedisPSKey(pubKey string) string {
	return fmt.Sprintf("sps_%s", pubKey)
}
