package storage

import (
	"hash/crc32"
	"os"
	"time"

	"github.com/mguzelevich/go-log"
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

	func() {
		for {
			timeout := time.After(10 * time.Second)
			select {
			case <-done:
				done = nil
			case <-timeout:
				log.Warning.Printf("storage shutdowned by timer")
				return
			default:
				if done == nil {
					return
				}
			}
		}
	}()
	log.Info.Printf("storage gracefully stopped")
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

func Dump() {

}

func Accounts() ([]*Account, error) {
	accounts, err := dataStorage.engine.Accounts()
	return accounts, err
}

func CreateAccount(account *Account) (*Account, error) {
	account.Id = account.Pk()
	entity, err := dataStorage.engine.Create(account)
	return entity.(*Account), err
}

func ReadAccount(account *Account) (*Account, error) {
	entity, err := dataStorage.engine.Read(account)
	return entity.(*Account), err
}

func UpdateAccount(account *Account) (*Account, error) {
	entity, err := dataStorage.engine.Update(account)
	return entity.(*Account), err
}

func DeleteAccount(account *Account) error {
	_, err := dataStorage.engine.Delete(account)
	return err
}

func GetDevice(device *Device) (*Device, error) {
	entity, err := dataStorage.engine.Read(device)
	return entity.(*Device), err
}

func init() {
	crc32q = crc32.MakeTable(0xD5828281)
}
