package log

import (
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
