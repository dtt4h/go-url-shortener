package main

import (
	"fmt"
	"go-url-shortener/internal/config"
)

func main() {
	config := config.MustLoad()

	fmt.Println(config)
	// TODO: init config: cleanenv 	+
	// TODO: init logger: slog
	// TODO: init storage: sqlite
	// TODO: init router: chi, render
	// TODO: run server
}
