package schampionne

import (
	"context"
	"log"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	grpc BrokerClient
}

// New instantiate a grpc connection and returns a Client wrapper
func NewClient(addr string) *Client {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[ERR ] Could not create Schampionne client. Reason: %s\n", err)
	}

	client := NewBrokerClient(conn)

	return &Client{
		conn: conn,
		grpc: client,
	}
}

func (c *Client) Close() {
	c.conn.Close()

}

// Whisper fires a Rumor through Schampionne network.
// t (type) string is used as identifier and source for listeners inside the network
// m (message) string contains data sent through the network along the Type
// returns an acknoledgment (Ack). Best case should returb ACK_OK code
func (c *Client) Whisper(t, m string) (*Ack, error) {
	id, err := uuid.NewV4()

	if err != nil {
		log.Printf("[WARN] Could not generate uuid v4")
	}

	return c.WhisperRumor(&Rumor{
		Type:    t,
		Message: m,
		Id:      id.String(),
	})
}

// WhisperRumor sends a rumor through Schampionne network.
func (c *Client) WhisperRumor(r *Rumor) (*Ack, error) {
	log.Printf("[INFO] Sending Rumor %+v\n", r)
	return c.grpc.Whisper(context.Background(), r)
}
