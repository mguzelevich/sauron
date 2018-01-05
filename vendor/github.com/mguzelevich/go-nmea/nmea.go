package nmea

import (
	//"fmt"

	"github.com/mguzelevich/go-nmea/messages"
)

func Marshal(msg *messages.NmeaMessage) ([]byte, error) {
	return nil, nil
}

func Unmarshal(data []byte) (messages.NmeaMessage, error) {
	return messages.ParseMessage(data)
}
