package utils

import (
	"os"
)

var Tmp string = ""

func InitTmp() error {
	var err error
	Tmp, err = os.MkdirTemp("", "wwise-teller-")
	if err != nil {
		return err
	}
	return nil
}

func CleanTmp() error {
	return os.RemoveAll(Tmp)
}


