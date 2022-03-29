package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
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

const url = "https://raw.githubusercontent.com/PoteeDev/glossary/main/terms.yml"
const password = "g1thub_h@rdc0d3_p@$$w0rd"

func Download(c *cli.Context) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http.Get -> %v", err)
		return err
	}

	// We read all the bytes of the image
	// Types: data []byte
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll -> %v", err)
		return err
	}
	resp.Body.Close()

	filename := "terms.yml"

	if err = ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println("Error Saving:", filename, err)
	} else {
		log.Println("Saved:", filename)
	}
	return nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "Find word",
				Action:  Find,
			},
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "Download latest glossary",
				Action:  Download,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
