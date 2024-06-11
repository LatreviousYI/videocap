package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var Rwlock sync.RWMutex


func Save(fileName string,data any) error{
	Rwlock.Lock()
	defer Rwlock.Unlock()
	basePath := GetExcutePath()
	dataFilepath:= filepath.Join(basePath,"data",fileName)
	file,err := os.Create(dataFilepath)
	if err != nil{
		return err
	}
	encoder := toml.NewEncoder(file)
    if err := encoder.Encode(data); err != nil {
        return err
    }
	return nil
}

func Get(fileName string,data any) error{
	basePath := GetExcutePath()
	Rwlock.RLock()
	defer Rwlock.RUnlock()
	dataFilepath:= filepath.Join(basePath,"data",fileName)
	_,err := toml.DecodeFile(dataFilepath,data)
	if err != nil{
		return err
	}
	return nil
}