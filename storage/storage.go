package storage

import (
	"hash/crc32"
	"os"
	"time"

	"github.com/mguzelevich/go.log"
)

var (
	dataStorage *Storage
	crc32q      *crc32.Table
)

type Storage struct {
	engine StorageEngine

	engineShutdownChan chan bool
	doneChan           chan bool
}

func (s Storage) DoneChan() chan bool {
	return s.doneChan
}

func (s *Storage) shutdownLoop(shutdownChan chan bool) {
	<-shutdownChan
	close(s.engineShutdownChan)

	done := s.engine.DoneChan()
	f := func() bool {
		for {
			timeout := time.After(10 * time.Second)
			select {
			case <-done:
				done = nil
			case <-timeout:
				return false
			default:
				if done == nil {
					return true
				}
			}
		}
	}

	if ok := f(); ok {
		log.Info.Printf("storage gracefully stopped")
	} else {
		log.Warning.Printf("storage shutdowned by timer")
	}

	close(s.doneChan)
}

func Init(engineFabric func(params map[string]string, shutdownChan chan bool) (StorageEngine, error), params map[string]string, shutdownChan chan bool) *Storage {
	dataStorage = &Storage{
		engineShutdownChan: make(chan bool),
		doneChan:           make(chan bool),
	}

	if engine, err := engineFabric(params, dataStorage.engineShutdownChan); err != nil {
		log.Error.Printf("storage init error %s", err)
		os.Exit(1)
	} else {
		dataStorage.engine = engine
		log.Info.Printf("storage started")
	}

	go dataStorage.shutdownLoop(shutdownChan)
	return dataStorage
}

func DumpAll() ([]byte, error) {
	return dataStorage.engine.DumpAll()
}

func Create(e Entity) (Entity, error) {
	return dataStorage.engine.Create(e)
}

func Read(e Entity) (Entity, error) {
	return dataStorage.engine.Read(e)
}

func Update(e Entity) (Entity, error) {
	return dataStorage.engine.Update(e)
}

func Delete(e Entity) (Entity, error) {
	return dataStorage.engine.Delete(e)
}

func Accounts() ([]*Account, error) {
	accounts, err := dataStorage.engine.Accounts()
	return accounts, err
}

func init() {
	crc32q = crc32.MakeTable(0xD5828281)
}
