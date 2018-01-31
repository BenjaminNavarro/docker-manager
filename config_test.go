package main

import (
	"fmt"
	"log"
	"os/user"
	"testing"
)

func TestConfigParser(t *testing.T) {
	var configurations []ImageConfiguration
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ReadDockerConfiguration(currentUser.HomeDir+"/.docker-manager.yaml", &configurations)
	for idx := range configurations {
		configurations[idx].Print()
		fmt.Println("\n\tStart command: ", configurations[idx].GenerateStartCommand())
	}
}
