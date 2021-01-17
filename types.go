package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Payload struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type Config struct {
	Cookie struct {
		Host    string `yaml:"host"`
		Method  string `yaml:"method"`
		URL     string `yaml:"url"`
		Headers []struct {
			Key string `yaml:"key"`
			Val string `yaml:"val"`
		} `yaml:"headers"`
	} `yaml:"cookie"`
	Login struct {
		Method string `yaml:"method"`
		URL    string `yaml:"url"`
		Auth   struct {
			Un string `yaml:"un"`
			Pw string `yaml:"pw"`
		} `yaml:"auth"`
		Headers []struct {
			Key string `yaml:"key"`
			Val string `yaml:"val"`
		} `yaml:"headers"`
	} `yaml:"login"`
	Historysearch struct {
		URL     string `yaml:"url"`
		Method  string `yaml:"method"`
		Headers []struct {
			Key string `yaml:"key"`
			Val string `yaml:"val"`
		} `yaml:"headers"`
	} `yaml:"historysearch"`
	Historytext struct {
		URL     string `yaml:"url"`
		Method  string `yaml:"method"`
		Headers []struct {
			Key string `yaml:"key"`
			Val string `yaml:"val"`
		} `yaml:"headers"`
	} `yaml:"historytext"`
	Currentfaults struct {
		URL     string `yaml:"url"`
		Method  string `yaml:"method"`
		Headers []struct {
			Key string `yaml:"key"`
			Val string `yaml:"val"`
		} `yaml:"headers"`
	} `yaml:"currentfaults"`
	MachineNames []string `yaml:"machinenames"`
}

type lhJSON struct {
	Count int `json:"count"`
	Data  []struct {
		ID        int    `json:"id"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
		Text      string `json:"text"`
		Textmap   struct {
			DE string `json:"DE"`
			EN string `json:"EN"`
		} `json:"textmap"`
		NodeID     string `json:"nodeId"`
		NodePath   string `json:"nodePath"`
		TKey       string `json:"tKey"`
		TIdx       int    `json:"tIdx"`
		Definition struct {
			Key      string `json:"key"`
			Type     string `json:"type"`
			Severity string `json:"severity"`
			DocNr    int    `json:"docNr"`
		} `json:"definition"`
	} `json:"data"`
}

func GetConfig() {
	if _, err := os.Stat("./config/config.yml"); err == nil { // check if config file exists
		configYaml, err := ioutil.ReadFile("./config/config.yml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(configYaml, &CONFIG)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		// config file not included, use embedded config
		// this took more time to fix than I'd like to admit
		// configYaml, err := Asset("./config/config.yml")
		configYaml, err := Asset("config/config.yml")
		if err != nil {
			fmt.Println("Asset was not found")
			panic(err)
		}
		err = yaml.Unmarshal(configYaml, &CONFIG)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Schrodinger: file may or may not exist. See err for details.")
	}
}
