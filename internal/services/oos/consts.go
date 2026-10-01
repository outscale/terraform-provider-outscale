package oos

import "time"

const (
	CreateDefaultTimeout time.Duration = 20 * time.Minute
	ReadDefaultTimeout   time.Duration = 20 * time.Minute
	UpdateDefaultTimeout time.Duration = 20 * time.Minute
	DeleteDefaultTimeout time.Duration = 60 * time.Minute
)
