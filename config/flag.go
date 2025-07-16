package config

import "flag"

func parseMode() string {
	mode := flag.String("mode", "dev", "Execution mode: dev | prod")
	flag.Parse()
	return *mode
}
