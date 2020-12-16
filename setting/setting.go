package setting

import (
	"github.com/go-ini/ini"
	"log"
)

var Cfg *ini.File

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal("Fail to parse app.ini: ", err)
	}

}
