package handlers

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/gomodule/redigo/redis"
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

	var pubSubConn redis.PubSubConn

	var (
		// mt is the message type
		// noFrame = -1, TextMessage = 1, BinaryMessage = 2
		mt  int
		msg []byte
	)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			break
		}

		if mt != 2 {
			continue
		}

		switch socketState {
		case codeSentState:
			_, valid := sign.Open(nil, msg, &pubKeyBytes)
			if !valid {
				c.Close()

				return
			}

			redisPoolConn := h.Cfg.Redis.Pool.Get()
			pubSubConn = redis.PubSubConn{Conn: redisPoolConn}

			if err := pubSubConn.PSubscribe(getRedisPSKey(pubKey)); err != nil {
				log.Printf("Redis error: %s", err)

				c.Close()
				return
			}

			go func() {
				for {
					switch v := pubSubConn.Receive().(type) {
					case redis.Message:
						c.WriteMessage(1, v.Data)
					}
				}
			}()

			c.SetCloseHandler(func(code int, msg string) error {
				go log.Println("Bye")
				go pubSubConn.Close()

				return nil
			})

			c.WriteMessage(1, []byte("OK"))

			if err != nil {
				log.Fatalln(err)
			}

			socketState = authenticatedState

			break

		case authenticatedState:
			break

		default:
			break
		}
	}
}

func getRedisPSKey(pubKey string) string {
	return fmt.Sprintf("sps_%s", pubKey)
}
