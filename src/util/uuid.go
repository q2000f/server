package util

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

func InitUUID(machineID func() (uint16, error)) {
	st := sonyflake.Settings{
		StartTime: time.Unix(1555929296, 0), //2019-04-22 do not modify
		MachineID: machineID,
	}

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GetNewID() uint64 {
	if sf == nil {
		panic("sonyflake not created")
	}

	id, err := sf.NextID()
	if err != nil {
		panic(nil)
	}

	return id
}
