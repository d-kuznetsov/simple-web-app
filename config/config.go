package config

import "flag"

var Port string
var MongoURI string

func init() {
	flag.StringVar(&Port, "port", "8080", "port")
	flag.StringVar(&MongoURI, "mongouri", "mongodb://localhost:27017", "mongouri")
	flag.Parse()
}
