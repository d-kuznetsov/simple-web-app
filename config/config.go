package config

import "flag"

var Port string
var MongoURI string
var SessionKey string

func init() {
	flag.StringVar(&Port, "port", "8080", "port")
	flag.StringVar(&MongoURI, "mongouri", "mongodb://localhost:27017", "mongouri")
	flag.StringVar(&SessionKey, "sessionkey", "session-key", "sessionkey")
	flag.Parse()
}
