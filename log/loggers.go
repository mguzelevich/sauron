package log

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Logger struct {
	Trace   io.Writer
	Debug   io.Writer
	Info    io.Writer
	Warning io.Writer
	Error   io.Writer
}

func InitLoggers(logger *Logger) {
	if logger.Trace != nil {
		Trace = log.New(logger.Trace, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logger.Debug != nil {
		Debug = log.New(logger.Debug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logger.Info != nil {
		Info = log.New(logger.Info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logger.Warning != nil {
		Warning = log.New(logger.Warning, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if logger.Error != nil {
		Error = log.New(logger.Error, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func init() {
	Trace = log.New(ioutil.Discard, "", log.Ldate)
	Debug = log.New(ioutil.Discard, "", log.Ldate)
	Info = log.New(ioutil.Discard, "", log.Ldate)
	Warning = log.New(ioutil.Discard, "", log.Ldate)
	Error = log.New(ioutil.Discard, "", log.Ldate)
}

// UUID generates a random UUID according to RFC 4122
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

/*

init:

```
func main() {
	log.InitLoggers(&log.Logger{
		ioutil.Discard,
		ioutil.Discard,
		os.Stdout,
		os.Stdout,
		os.Stderr,}
	)

	log.InitLoggers(&log.Logger{ Error: os.Stderr })

}
```

usage

```
import ".../log"

log.Debug.Printf("some debug message")
...

log.Error.Printf("some error message %s", err)

```
*/
