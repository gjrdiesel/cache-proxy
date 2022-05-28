package config_factory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Instance struct {
	ProgramName string
	Config      *Config
	File        string
}

func New(Program string) *Instance {
	i := &Instance{
		ProgramName: Program,
	}

	i.EnsureDirectoryExists()
	i.File = i.Dir() + "config.json"

	return i
}

func homeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return homeDir
}

type Config struct {
	SiphonUrl *string `json:"siphonUrl"`
	Port      *string `json:"port"`
}

func (i *Instance) Dir() string {
	return homeDir() + "/.config/" + i.ProgramName + "/"
}

func (i *Instance) ReadConfig() *Config {
	s := os.Getenv("SIPHON_URL")
	p := os.Getenv("PORT")

	if p != "" && s != "" {
		i.Config = &Config{
			SiphonUrl: &s,
			Port:      &p,
		}
		return i.Config
	}

	body, err := ioutil.ReadFile(i.File)
	if err != nil {
		log.Printf("Unable to read config file: %v\n", err)
	}

	err = json.Unmarshal(body, &i.Config)
	if err != nil {
		log.Println("Unable to unmarshal config file: ", err)
		return nil
	}

	return i.Config
}

func (i *Instance) AskConfig() *Config {
	fmt.Printf("Enter a URL to siphon from: ")
	var u string
	fmt.Scanln(&u)

	fmt.Printf("Port for siphon to listen on (8080): ")
	var p string
	fmt.Scanln(&p)

	if p == "" {
		p = "8080"
	}

	i.Config = &Config{
		SiphonUrl: &u,
		Port:      &p,
	}

	return i.Config
}

func (i *Instance) WriteConfig() {
	s, err := json.Marshal(i.Config)
	if err != nil {
		log.Println("Could not write config, ", err)
		return
	}

	err = ioutil.WriteFile(i.File, s, os.ModePerm)
	if err != nil {
		log.Println("Could not write config file", err)
	}
}

func (i *Instance) Settings() *Config {
	config := i.ReadConfig()
	if config != nil {
		return config
	}

	return i.RedoConfiguration()
}

func (i *Instance) RedoConfiguration() *Config {
	config := i.AskConfig()

	i.WriteConfig()

	return config
}

func (i *Instance) EnsureDirectoryExists() {
	err := os.MkdirAll(i.Dir(), os.ModePerm)
	if err != nil {
		log.Println("Unable to ensure directory exists,", err)
		return
	}
}
