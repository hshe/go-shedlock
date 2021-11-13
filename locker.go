package shedlock

import (
	"os"
	"time"
)

/**

	String getInsertStatement() {
return "INSERT INTO " + this.tableName() + "(" + this.name() + ", " + this.lockUntil() + ", " + this.lockedAt() + ", " + this.lockedBy() + ") VALUES(:name, :lockUntil, :now, :lockedBy)";
}

public String getUpdateStatement() {
return "UPDATE " + this.tableName() + " SET " + this.lockUntil() + " = :lockUntil, " + this.lockedAt() + " = :now, " + this.lockedBy() + " = :lockedBy WHERE " + this.name() + " = :name AND " + this.lockUntil() + " <= :now";
}

public String getExtendStatement() {
//return "UPDATE " + this.tableName() + " SET " + this.lockUntil() + " = :lockUntil WHERE " + this.name() + " = :name AND " + this.lockedBy() + " = :lockedBy AND " + this.lockUntil() + " > :now";
//}
//
//public String getUnlockStatement() {
//return "UPDATE " + this.tableName() + " SET " + this.lockUntil() + " = :unlockTime WHERE " + this.name() + " = :name";
//}
//
//String name() {
//return this.configuration.getColumnNames().getName();
//}
//
//String lockUntil() {
//return this.configuration.getColumnNames().getLockUntil();
//}
//
//String lockedAt() {
//return this.configuration.getColumnNames().getLockedAt();
//}
//
//String lockedBy() {
//return this.configuration.getColumnNames().getLockedBy();
//}
//
//String tableName() {
//return this.configuration.getTableName();
//}
*/

type ShedLock struct {
	Name      string    `gorm:"name"`
	LockUntil time.Time `gorm:"lock_until"`
	LockedAt  time.Time `gorm:"locked_at"`
	LockedBy  string    `gorm:"locked_by" hostname+name`
}

type Schedules struct {
	Name     string
	Spec     string
	cmd      func()
	LockTime int `json:"lock_time" 5min`
}

type Locker interface {
	Add(name string, spec string, cmd func()) error
	Adds(schedules []*Schedules) error
	DoLock(name string) bool
	Insert(name string) bool
	Find(name string) *Schedules
	Update(name string) bool
	Unlock(name string) bool
}

func LocalHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
