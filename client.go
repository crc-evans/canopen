package canopen

import (
	"time"

	"github.com/brutella/can"
)

// A Client handles message communication by sending a request
// and waiting for the response.
type Client struct {
	Bus     *can.Bus
	Timeout time.Duration
}

// Do sends a request and waits for a response. It uses a custom CAN frame filter.
// If the response frame doesn't arrive on time, an error is returned.
func (c *Client) Do(req *Request) (*Response, error) {
	ch := can.WaitFunc(c.Bus, req.FilterFunc, c.Timeout)
	if err := c.Bus.Publish(req.Frame.CANFrame()); err != nil {
		return nil, err
	}
	resp := <-ch
	if resp.Err != nil {
		// return nil, fmt.Errorf("request Cob ID %#v failed: %v", req.Frame.CobID, resp.Err)
		return nil, resp.Err
	}
	return &Response{CANopenFrame(resp.Frame), req}, nil
}
