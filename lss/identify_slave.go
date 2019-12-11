package lss

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/brutella/can"
	"github.com/brutella/canopen"
)

type identifySlaveRequest struct {
	IDNumber   uint32
	BitChecked uint8
	LSSSub     uint8
	LSSNext    uint8
}

func newIdentifySlaveRequest(idNumber uint32, bitChecked uint8, lssSub uint8, lssNext uint8) *identifySlaveRequest {
	return &identifySlaveRequest{
		IDNumber:   idNumber,
		BitChecked: bitChecked,
		LSSSub:     lssSub,
		LSSNext:    lssNext,
	}
}

func (fs *identifySlaveRequest) MarshalBytes() []byte {
	data := make([]byte, 8)
	data[0] = 0x51
	binary.LittleEndian.PutUint32(data[1:5], fs.IDNumber)
	data[5] = fs.BitChecked
	data[6] = fs.LSSSub
	data[7] = fs.LSSNext
	return data
}

func (msg *identifySlaveRequest) Do(bus *can.Bus, timeout time.Duration) error {
	c := canopen.Client{
		Bus:     bus,
		Timeout: timeout,
	}
	b := msg.MarshalBytes()
	frm := canopen.NewFrame(MessageTypeLSSSlave, b)
	req := canopen.NewRequest(frm, MessageTypeLSSMaster)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	cs := resp.Frame.Data[0]
	if cs != 0x50 {
		return fmt.Errorf("LSS: unexpected command specifier 0x%x", cs)
	}
	return nil
}
