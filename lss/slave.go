package lss

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/brutella/can"
	"github.com/brutella/canopen"
)

type Slave struct {
	VendorID       uint32
	ProductCode    uint32
	RevisionNumber uint32
	SerialNumber   uint32
}

func (s *Slave) Select(bus *can.Bus) error {
	c := canopen.Client{
		Bus:     bus,
		Timeout: 2 * time.Second,
	}
	data := make([]byte, 8)

	data[0] = 0x40
	binary.LittleEndian.PutUint32(data[1:5], s.VendorID)
	frm := canopen.NewFrame(MessageTypeLSSSlave, data)
	err := bus.Publish(frm.CANFrame())
	if err != nil {
		return err
	}

	data[0] = 0x41
	binary.LittleEndian.PutUint32(data[1:5], s.ProductCode)
	frm = canopen.NewFrame(MessageTypeLSSSlave, data)
	err = bus.Publish(frm.CANFrame())
	if err != nil {
		return err
	}

	data[0] = 0x42
	binary.LittleEndian.PutUint32(data[1:5], s.RevisionNumber)
	frm = canopen.NewFrame(MessageTypeLSSSlave, data)
	err = bus.Publish(frm.CANFrame())
	if err != nil {
		return err
	}

	// last ID; expect response from node
	data[0] = 0x43
	binary.LittleEndian.PutUint32(data[1:5], s.SerialNumber)
	frm = canopen.NewFrame(MessageTypeLSSSlave, data)
	req := canopen.NewRequest(frm, MessageTypeLSSMaster)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	cs := resp.Frame.Data[0]
	// TODO: also check resp DLC=8 and rest of bytes are zero?
	if cs != 0x44 {
		return fmt.Errorf("LSS: unexpected command specifier 0x%x", cs)
	}
	return nil
}

func (s *Slave) SetNodeID(bus *can.Bus, nodeID uint8) error {
	c := canopen.Client{
		Bus:     bus,
		Timeout: 2 * time.Second,
	}
	data := make([]byte, 8)
	data[0] = 0x11
	data[1] = nodeID
	frm := canopen.NewFrame(MessageTypeLSSSlave, data)
	req := canopen.NewRequest(frm, MessageTypeLSSMaster)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	cs := resp.Frame.Data[0]
	// TODO: also check resp DLC=8 and rest of bytes are zero?
	if cs != 0x11 {
		return fmt.Errorf("LSS: unexpected command specifier 0x%x", cs)
	}
	return nil
}
