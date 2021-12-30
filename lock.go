package shedlock

import (
	"github.com/robfig/cron"
	"os"
	"time"
)

type QuartzJob interface {
	cron.Job
}

type ShedLock struct {
	Name      string    `gorm:"name"`
	LockUntil time.Time `gorm:"lock_until"`
	LockedAt  time.Time `gorm:"locked_at"`
	LockedBy  string    `gorm:"locked_by" hostname+name`
}

type Schedule struct {
	Name     string
	Spec     string
	Cmd      func()
	Job      cron.Job
	LockTime int `json:"lock_time" 20 sec`
}

type Locker interface {
	AddFun(name string, spec string, cmd func()) error
	AddSchedules(schedules []*Schedule) error
	DoLock(name string) bool
	Insert(name string) bool
	Find(name string) *Schedule
	Update(name string) bool
	Unlock(name string) bool

	Start()
	Stop()
}

func LocalHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
