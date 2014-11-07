package main

import (
	"os"

	"github.com/dynport/urknall"
)

func main() {
	defer urknall.OpenLogger(os.Stdout).Close()
	t, err := urknall.NewSshTarget("root@128.199.54.55")
	if err != nil {
		panic("error parsing ssh target: " + err.Error())
	}
	err = urknall.Run(t, &Template{})
	if err != nil {
		logger.Fatal("ERROR: ", err)
	}
}
