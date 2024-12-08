package model

import "time"

type Notification struct {
	TelegramChat      int64
	NotificationDelay int
	SubjectName       string
	Room              string
	Teacher           string
	StartTime         time.Time
	EndTime           time.Time
}
