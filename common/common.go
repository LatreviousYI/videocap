package common

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var CronTab *cron.Cron

var entryid cron.EntryID

var enidlock sync.RWMutex

var mjpgPid int = -1

var mjpgPidlock sync.RWMutex


func init(){
	CronTab = cron.New(cron.WithSeconds())
	CronTab.Start()
}

func GetEntryId() cron.EntryID{
	enidlock.RLock()
	defer enidlock.RUnlock()
	return entryid
}

func EditEntryId(id cron.EntryID){
	enidlock.Lock()
	defer enidlock.Unlock()
	entryid = id
}

func GetMjpgPid() int{
	mjpgPidlock.RLock()
	defer mjpgPidlock.RUnlock()
	return mjpgPid
}

func EditMjpgPid(id int){
	mjpgPidlock.Lock()
	defer mjpgPidlock.Unlock()
	mjpgPid = id
}