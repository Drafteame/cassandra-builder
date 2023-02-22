package models

import (
	"time"
)

type Consistency uint16

const (
	Any         Consistency = 0x00
	One         Consistency = 0x01
	Two         Consistency = 0x02
	Three       Consistency = 0x03
	Quorum      Consistency = 0x04
	All         Consistency = 0x05
	LocalQuorum Consistency = 0x06
	EachQuorum  Consistency = 0x07
	LocalOne    Consistency = 0x0A
)

// Config is the main cassandra configuration needed
type Config struct {
	Port                     int
	KeyspaceName             string
	Username                 string
	Password                 string
	ContactPoints            []string
	Debug                    bool
	ProtoVersion             int
	Consistency              Consistency
	CaPath                   string
	DisableInitialHostLookup bool
	Timeout                  time.Duration
	ConnectTimeout           time.Duration
	NumRetries               uint
}
