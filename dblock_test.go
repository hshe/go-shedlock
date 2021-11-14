package shedlock

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(6)

	cron := cron.New()
	cron.Start()
	cron.Stop()
	cron.AddFunc("* * * * * ?", func() {
		t.Log("xxxxxx")
		wg.Done()
	})
	cron.Start()
	select {
	case <-time.After(10 * time.Second):
		// No job ran!
	case <-wait(wg):
		t.Fatal("expected stopped cron does not run any job")
	}
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func stop(cron *cron.Cron) chan bool {
	ch := make(chan bool)
	go func() {
		cron.Stop()
		ch <- true
	}()
	return ch
}

func TestLockerDb_Add(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/icc311?charset=utf8&parseTime=true&loc=Local"), &gorm.Config{
		//不打印日志
		//Logger: logger.New(
		//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//	logger.Config{
		//		SlowThreshold: time.Second,   // 慢 SQL 阈值
		//		LogLevel:      logger.Silent, // Log level
		//		Colorful:      false,         // 禁用彩色打印
		//	},
		//),
	})
	l := NewLockerDb(db)
	l.Add("bbbbb", "* * * * * ?", func() {
		fmt.Println("test...")
	})
	go func() {
		l.Add("bbbbb", "* * * * * ?", func() {
			fmt.Println("test...")
		})
	}()
	//select {}
}
