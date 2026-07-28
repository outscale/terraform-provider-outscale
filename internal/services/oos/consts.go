package oos

import "time"

const (
	CreateDefaultTimeout time.Duration = 10 * time.Minute
	ReadDefaultTimeout   time.Duration = 2 * time.Minute
	UpdateDefaultTimeout time.Duration = 10 * time.Minute
	DeleteDefaultTimeout time.Duration = 5 * time.Minute
)
