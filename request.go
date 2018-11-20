package canopen

import (
	"encoding/binary"

	"github.com/brutella/can"
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
	return func(frm can.Frame) bool {
		if respID != frm.ID {
			return false
		}
		reqData := req.Data
		frmData := frm.Data
		reqIndex, reqSubindex := binary.LittleEndian.Uint16(reqData[1:3]), reqData[3]
		frmIndex, frmSubindex := binary.LittleEndian.Uint16(frmData[1:3]), frmData[3]
		if frmIndex != reqIndex || frmSubindex != reqSubindex {
			return false
		}
		return true
	}
}
