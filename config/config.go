package config

import "flag"

var Port string

func init() {
	flag.StringVar(&Port, "port", "8080", "port")
	flag.Parse()
}
