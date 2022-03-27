package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Terms struct {
	List []Term `yaml:"terms"`
}

type Term struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Human       string `yaml:"human"`
}

func (t *Terms) Get() *Terms {

	yamlFile, err := ioutil.ReadFile("terms.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, t)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return t
}

func (t *Terms) Search(word string) (string, string) {
	for _, term := range t.List {
		if strings.ToLower(term.Name) == strings.ToLower(word) {
			return term.Description, term.Human
		}
	}
	return "", ""
}

func Find(c *cli.Context) error {
	word := c.Args().First()
	if word == "" {
		cli.ShowAppHelpAndExit(c, 1)
	}
	var terms Terms
	terms.Get()
	description, human := terms.Search(word)
	if description != "" && human != "" {
		fmt.Println("search result for:", c.Args().First())
		fmt.Println("description:", description)
		fmt.Println("human:", human)
	}

	return nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "find word",
				Action:  Find,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
