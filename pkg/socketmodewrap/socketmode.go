package socketmodewrap

import (
	"fmt"

	"github.com/slack-go/slack/socketmode"
)

// the pkg folder is for packages intended to be used in the another projects
// do we have such plans for this one?

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
