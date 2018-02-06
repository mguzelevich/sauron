package loggers

import (
	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/storage"
)

// type LoggerServer interface {
// 	New(addr string)
// 	ListenAndServe(shutdownChan chan bool)
// 	DoneChan() chan bool
// }

type Message interface {
	Device() *storage.Device
	Telemetry() *storage.Telemetry
}

type msgHandler func(raw string) (Message, error)

type Logger struct {
	handler func(raw string) (Message, error)

	processMsgChan chan string
}

func (l *Logger) processLoop() {
	for raw := range l.processMsgChan {
		msg, err := l.handler(raw)
		if err != nil {
			log.Error.Printf("parse packet [%q] error [%s]", raw, err)
			continue
		}
		l.save(msg)
	}
}

func (l *Logger) Log(raw string) {
	go func() {
		l.processMsgChan <- raw
	}()
}

func (l *Logger) save(msg Message) {
	/*
		получить хеш девайса из телеметрии
		true:
			получить инстанс девайса из стораджа
			получить инстанс аккаунта
		false:
			создать инстанс девайса
			приааттачить к анонимусу (создать новый аккаунт)

		получить инстанс девайса из аккаунта
		записать телеметрию в девайс
	*/

	/*
		получить хеш девайса из телеметрии
		получить инстанс девайса из аккаунта
		записать телеметрию в девайс
	*/

	account := &storage.Account{}

	device, err := storage.GetDevice(msg.Device())
	if err != nil {
		switch err {
		case storage.ErrEntityNotFound:
			account, err = storage.ReadAccount(account)
			if err != nil {
				switch err {
				case storage.ErrEntityNotFound:
					account, _ = storage.CreateAccount(account)
				default:
					panic(err)
				}
			}
			device, _ = account.CreateDevice(device)
		default:
			panic(err)
		}
	} else {
		//account, _ := storage.ReadAccount(&storage.Account{Id: device.UserId})
	}
	device.AddTelemetry(msg.Telemetry())
}

func New(handler msgHandler) *Logger {
	logger := &Logger{
		handler:        handler,
		processMsgChan: make(chan string),
	}
	go logger.processLoop()
	return logger
}
