package oapi

import "time"

const (
	CreateDefaultTimeout time.Duration = 10 * time.Minute
	ReadDefaultTimeout   time.Duration = 5 * time.Minute
	UpdateDefaultTimeout time.Duration = 10 * time.Minute
	DeleteDefaultTimeout time.Duration = 5 * time.Minute

	MinPort     int   = 1
	MaxPort     int   = 65535
	MinIops     int32 = 100
	MaxIops     int32 = 13000
	DefaultIops int32 = 150
	MaxSize     int32 = 14901

	AwaitActiveStateDefaultValue          bool = true
	RemoveDefaultOutboundRuleDefaultValue bool = false
)
