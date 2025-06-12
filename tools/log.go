package tools

import (
	"log"

	"github.com/fatih/color"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func LogInfo(s string) {
	log.Printf("["+green("INFO")+"] %s\n", s)
}

func LogWarn(s string) {
	log.Printf("["+yellow("WARN")+"] %s\n", s)
}

func LogFatal(v ...any) {
	log.Fatal(v...)
}

func LogOk(s string) {
	log.Printf("["+green("OK")+"] %s\n", s)
}
