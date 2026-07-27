package main

import (
	"fmt"
	"log"

	"gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error al leer la configuración inicial: %v", err)
	}

	err = cfg.SetUser("lane")
	if err != nil {
		log.Fatalf("error al actualizar el usuario: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error al releer la configuración: %v", err)
	}

	fmt.Printf("%+v\n", cfg)
}
