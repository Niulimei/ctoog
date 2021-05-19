// Code generated by go-swagger; DO NOT EDIT.

package main

import (
	"ctgb/utils"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sevlyar/go-daemon"
	"gopkg.in/yaml.v2"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"ctgb/restapi"
	"ctgb/restapi/operations"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	lvl = strings.ToLower(lvl)
	if ok {
		if lvl == "debug" {
			log.SetLevel(log.DebugLevel)
		} else if lvl == "info" {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
	}
}

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	usage := "Usage: ./translator-server start | stop | status | restart"
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "start":
			fmt.Println("translator-server is starting.")
			break
		case "stop", "restart":
			pid, _ := ioutil.ReadFile("./translator-server.pid")
			if string(pid) == "" {
				fmt.Println("translator-server is not running.")
			} else if _, err := utils.Exec("ps " + string(pid)); err != nil {
				fmt.Println("translator-server is not running.")
			} else {
				utils.Exec("kill " + string(pid))
				fmt.Println("translator-server has been stopped.")
			}
			if command == "stop" {
				return
			} else {
				fmt.Println("translator-server is starting.")
				break
			}
		case "status":
			pid, _ := ioutil.ReadFile("./translator-server.pid")
			if string(pid) == "" {
				fmt.Println("translator-server is not running.")
			} else if _, err := utils.Exec("ps " + string(pid)); err != nil {
				fmt.Println("translator-server is not running.")
			} else {
				fmt.Println("translator-server is running.")
			}
			return
		default:
			fmt.Println(usage)
			return
		}
	}
	cntxt := &daemon.Context{
		PidFileName: "translator-server.pid",
		PidFilePerm: 0644,
		LogFileName: "translator-server.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"translator-server"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	type conf struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	tmp := &conf{}
	content, _ := ioutil.ReadFile("./translator-server.yaml")
	yaml.Unmarshal(content, tmp)
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTranslatorAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "translator"
	parser.LongDescription = swaggerSpec.Spec().Info.Description
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()
	server.Host = tmp.Host
	server.Port = tmp.Port
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
