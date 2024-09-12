package log

import (
	"fmt"
	"log"
)

func Red(s string) {
	log.Println("\033[31m" + s + "\033[0m")
}

func Redf(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	log.Println("\033[31m" + msg + "\033[0m")
}

func Green(s string) {
	log.Println("\033[32m" + s + "\033[0m")
}

func Greenf(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	log.Println("\033[32m" + msg + "\033[0m")
}

func Yellowf(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	log.Println("\033[33m" + msg + "\033[0m")
}

func Yellow(s string) {
	log.Println("\033[33m" + s + "\033[0m")
}

func Bluef(s string, args ...interface{}) {
	msg := fmt.Sprintf(s, args...)
	log.Println("\033[34m" + msg + "\033[0m")
}

func Blue(s string) {
	log.Println("\033[34m" + s + "\033[0m")
}
