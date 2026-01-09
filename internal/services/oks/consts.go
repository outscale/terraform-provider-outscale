package oks

import "time"

const (
	CreateDefaultTimeout  time.Duration = 15 * time.Minute
	ReadDefaultTimeout    time.Duration = 2 * time.Minute
	UpdateDefaultTimeout  time.Duration = 10 * time.Minute
	DeleteDefaultTimeout  time.Duration = 10 * time.Minute
	UpgradeDefaultTimeout time.Duration = 30 * time.Minute
)
