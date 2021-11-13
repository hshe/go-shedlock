package shedlock

import (
	"github.com/robfig/cron"
	"gorm.io/gorm"
	"time"
)

//https://blog.csdn.net/waltonhuang/article/details/106555195
//https://www.baeldung.com/shedlock-spring
//func _(db *gorm.DB) *DbLocker {
//
//	return nil
//}

//LockerDb
type LockerDb struct {
	db       *gorm.DB
	LockTime int `json:"lock_time" 5min`
}

func NewLockerDbFor(db *gorm.DB, lockTime int) *LockerDb {
	return &LockerDb{db: db, LockTime: lockTime}
}

func NewLockerDb(db *gorm.DB) *LockerDb {
	return &LockerDb{db: db, LockTime: 5}
}

func (l LockerDb) Add(name string, spec string, cmd func()) error {
	c := cron.New()
	err := c.AddFunc(spec, func() {
		if l.DoLock(name) {
			cmd()
		}
	})
	c.Start()
	return err
}

func (l LockerDb) Adds(schedules []*Schedules) error {
	c := cron.New()
	for i := range schedules {
		c.AddFunc(schedules[i].Spec, func() {
			if l.DoLock(schedules[i].Name) {
				schedules[i].cmd()
			}
		})
	}
	c.Start()
	return nil
}

func (l LockerDb) DoLock(name string) bool {
	return l.Insert(name) || l.Update(name)
}

func (l LockerDb) Insert(name string) bool {
	s := &ShedLock{
		Name:      name,
		LockUntil: time.Now().Add(time.Duration(l.LockTime) * time.Hour),
		LockedAt:  time.Now(),
		LockedBy:  LocalHostName(),
	}
	create := l.db.Table("shedlock").Create(&s)
	if create.Error == nil && create.RowsAffected > 0 {
		return true
	}
	return false
}

func (l LockerDb) Find(name string) *Schedules {
	res := &Schedules{}
	l.db.Table("shedlock").Where("locked_by =?", LocalHostName()).Where("name=?", name).Find(&res)
	return res
}

func (l LockerDb) Update(name string) bool {
	now := time.Now()
	s := &ShedLock{
		Name:      name,
		LockedAt:  now,
		LockedBy:  LocalHostName(),
		LockUntil: now.Add(time.Duration(l.LockTime) * time.Minute),
	}
	update := l.db.Table("shedlock").Where("name=?", name).Where("lock_until<=?", now).Updates(&s)
	if update.Error == nil && update.RowsAffected > 0 {
		return true
	}
	return false
}

func (l LockerDb) Unlock(name string) bool {
	panic("implement me")
}
