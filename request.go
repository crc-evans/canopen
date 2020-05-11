package canopen

import (
	"github.com/crc-evans/can"
)

// A Request represents a CANopen request published on a CAN bus and received by another CANopen node.
type Request struct {
	// The Frame of the request
	Frame Frame

	// The filter used for matching responses to requests
	FilterFunc func(can.Frame) bool
}

// NewRequest returns a request containing the frame to be sent
// and the expected response frame id.
func NewRequest(frm Frame, respID uint32) *Request {
	return &Request{
		Frame:      frm,
		FilterFunc: sdoFilter(frm, respID),
	}
}

func sdoFilter(req Frame, respID uint32) func(can.Frame) bool {
	return func(frm can.Frame) bool { return respID == frm.ID }
}
