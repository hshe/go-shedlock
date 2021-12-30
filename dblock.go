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
	LockTime int `json:"lock_time" 20 sec`
	c        *cron.Cron
}

func NewLockerDbFor(db *gorm.DB, lockTime int) *LockerDb {
	return &LockerDb{db: db, LockTime: lockTime, c: cron.New()}
}

func NewLockerDb(db *gorm.DB) *LockerDb {
	return &LockerDb{db: db, LockTime: 20, c: cron.New()}
}
func NewLockerDbWithLockTime(db *gorm.DB, lockFor int) *LockerDb {
	return &LockerDb{db: db, LockTime: lockFor, c: cron.New()}
}

func (l LockerDb) AddFun(name string, spec string, cmd func()) error {
	if l.c == nil {
		l.c = cron.New()
	}
	err := l.c.AddFunc(spec, func() {
		if l.DoLock(name) {
			defer l.Unlock(name)
			cmd()
		}
	})
	return err
}

func (l LockerDb) AddSchedules(schedules []*Schedule) error {
	if l.c == nil {
		l.c = cron.New()
	}
	for i := range schedules {
		l.c.AddFunc(schedules[i].Spec, func() {
			if !l.DoLock(schedules[i].Name) {
				return
			}
			l.Unlock(schedules[i].Name)
			if schedules[i].Job == nil {
				schedules[i].Cmd()
			} else {
				schedules[i].Job.Run()
			}
		})
	}
	return nil
}

func (l LockerDb) DoLock(name string) bool {
	return l.Insert(name) || l.Update(name)
}

func (l LockerDb) Insert(name string) bool {
	s := &ShedLock{
		Name:      name,
		LockUntil: time.Now().Add(time.Duration(l.LockTime) * time.Second),
		LockedAt:  time.Now(),
		LockedBy:  LocalHostName(),
	}
	create := l.db.Table("shedlock").Create(&s)
	if create.Error == nil && create.RowsAffected > 0 {
		return true
	}
	return false
}

func (l LockerDb) Find(name string) *Schedule {
	res := &Schedule{}
	l.db.Table("shedlock").Where("locked_by =?", LocalHostName()).Where("name=?", name).Find(&res)
	return res
}

func (l LockerDb) Update(name string) bool {
	now := time.Now()
	s := &ShedLock{
		Name:      name,
		LockedAt:  now,
		LockedBy:  LocalHostName(),
		LockUntil: now.Add(time.Duration(l.LockTime) * time.Second),
	}
	update := l.db.Table("shedlock").Where("name=?", name).Where("lock_until<=?", now).Updates(&s)
	if update.Error == nil && update.RowsAffected > 0 {
		return true
	}
	return false
}

func (l LockerDb) Unlock(name string) bool {
	now := time.Now()
	s := &ShedLock{
		LockUntil: now,
	}
	l.db.Table("shedlock").Where("name=?", name).Where("lock_until>?", now).Updates(&s)
	return true
}

func (l LockerDb) Start() {
	l.c.Start()
}

func (l LockerDb) Stop() {
	l.c.Stop()
	l.c = nil
}
