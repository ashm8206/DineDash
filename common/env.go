package common

import (
	"log"
	"syscall"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		log.Printf("Port Initialzed %s:", val)
		return val
	}
	log.Printf("Port Initialzed %s:", fallback)
	return fallback
}
