package logger

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"log"
	"os"
)

type Loggers struct {
	Info *log.Logger
	Err  *log.Logger
}

func InitLogger(m consts.Mode) (loggers Loggers) {
	switch m {
	case consts.Production:
		// TODO filenames to config
		infoF, err := os.OpenFile("../logs/info.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("%s", err)
		}
		defer infoF.Close()
		// TODO filenames to config
		errF, err := os.OpenFile("../logs/info.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		defer errF.Close()

		loggers.Info = log.New(infoF, "[INFO]\t", log.Ldate|log.Ltime)
		loggers.Err = log.New(errF, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	case consts.Development:
		loggers.Info = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
		loggers.Err = log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	}
	loggers.Info.Print("Executing InitLogger.")
	return
}
