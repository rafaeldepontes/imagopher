package tool

import (
	"log"
	"os"
)

func ChecksEnv(src *string) {
	if _, err := os.Stat(*src); err != nil {
		log.Println("[WARN] Couldn't load the '.env' file, trying to use '.env.example' instead...")
		*src = ".env.example"
	}
}
