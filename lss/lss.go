package lss

import (
	"time"

	"github.com/brutella/can"
)

const (
	MessageTypeLSSMaster = 0x7e4
	MessageTypeLSSSlave  = 0x7e5
)

func ResetAll(bus *can.Bus, timeout time.Duration) error {
	req := newIdentifySlaveRequest(0, uint8(0x80), 0, 0)
	if err := req.Do(bus, timeout); err != nil {
		return err
	}
	return nil
}

func Fastscan(bus *can.Bus, timeout time.Duration,
	vendorID uint32,
	productCode uint32,
	revisionNumber uint32,
	serialNumber uint32) (*Slave, error) {
	s := &Slave{
		VendorID:       vendorID,
		ProductCode:    productCode,
		RevisionNumber: revisionNumber,
		SerialNumber:   serialNumber,
	}
	if s.VendorID == 0 {
		id, err := lssFindID(bus, 0, timeout)
		if err != nil {
			return nil, err
		}
		s.VendorID = id
	}
	if s.ProductCode == 0 {
		id, err := lssFindID(bus, 1, timeout)
		if err != nil {
			return nil, err
		}
		s.ProductCode = id
	}
	if s.RevisionNumber == 0 {
		id, err := lssFindID(bus, 2, timeout)
		if err != nil {
			return nil, err
		}
		s.RevisionNumber = id
	}
	if s.SerialNumber == 0 {
		id, err := lssFindID(bus, 3, timeout)
		if err != nil {
			return nil, err
		}
		s.SerialNumber = id
	}
	return s, nil
}

func lssFindID(bus *can.Bus, lssSub uint8, timeout time.Duration) (uint32, error) {
	lssNumber := uint32(0)
	lssNext := uint8(lssSub)
	for bitChecked := 31; bitChecked >= 0; bitChecked-- {
		time.Sleep(timeout) // wait for any in-flight messages to pass...
		req := newIdentifySlaveRequest(lssNumber, uint8(bitChecked), lssSub, lssNext)
		if err := req.Do(bus, timeout); err != nil {
			if _, ok := err.(*can.ErrTimeout); !ok {
				// if any error occurred other than a timeout, return
				// fmt.Println("1")
				return 0, err
			}
			// if timeout, bit = 1; confirm
			lssNumber |= (1 << uint8(bitChecked))
			req = newIdentifySlaveRequest(lssNumber, uint8(bitChecked), lssSub, lssNext)
			if err := req.Do(bus, timeout); err != nil {
				// if timeout still occurred, no more nodes are answering; return
				// fmt.Println("2")
				return 0, err
			}
			continue
		}
		// if response, bit = 0
	}
	// send one last message to confirm the lssNumber
	if lssSub < 3 {
		lssNext = lssSub + 1
	}
	if lssSub == 3 {
		lssNext = 0
	}
	req := newIdentifySlaveRequest(lssNumber, 0, lssSub, lssNext)
	err := req.Do(bus, timeout)
	if err != nil {
		// we always expect a reply to the confirmation here, so
		// if any error occurred (INCLUDING a timeout), return
		// fmt.Println("3")
		return 0, err
	}
	return lssNumber, nil
}
