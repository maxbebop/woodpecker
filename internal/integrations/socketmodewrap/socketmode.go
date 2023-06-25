package socketmodewrap

import (
	"fmt"

	"github.com/slack-go/slack/socketmode"
)

type SocketmodeClient struct {
	c *socketmode.Client
}

func New(c *socketmode.Client) *SocketmodeClient {
	return &SocketmodeClient{c: c}
}

func (c *SocketmodeClient) EventsIn() <-chan socketmode.Event { return c.c.Events }
func (c *SocketmodeClient) Run() error {
	return fmt.Errorf("%w", c.c.Run())
}

func (c *SocketmodeClient) Ack(req socketmode.Request, payload ...any) { c.c.Ack(req, payload...) }
