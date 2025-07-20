package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DatabaseConnection string
	Port               string
}

func getConnectionStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"), os.Getenv("DATABASE"))
}

func LoadConfig() (Config, error) {
	connection := getConnectionStr()
	if connection == "host= port= user= password= database= sslmode=disable" {
		log.Fatal("evn variables not passed!")
		return Config{}, errors.New("evn variables not passed!")
	}
	return Config{DatabaseConnection: connection, Port: os.Getenv("PORT")}, nil
}
