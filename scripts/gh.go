package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cli/go-gh"
	"github.com/SamWolfs/rofi-web/cmd"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type Repository struct {
	Name string `mapstructure:"name"`
	NameWithOwner string `mapstructure:"nameWithOwner"`
	Url string `mapstructure:"url"`
}

func main() {
	repoList, _, err := gh.Exec("repo", "list", "--json", "name,nameWithOwner,url")
	if err != nil {
		log.Fatal(err)
	}
	var repositories []Repository
	json.Unmarshal(repoList.Bytes(), &repositories)

	links := funk.Map(repositories, func (repo Repository) cmd.Link {
		return cmd.Link{
			Name: repo.Name,
			Tags: repo.NameWithOwner,
			Url: repo.Url,
		}
	})

	var out []byte
	out, err = yaml.Marshal(&links)
	if err != nil {
			log.Fatalf("error: %v", err)
	}
	fmt.Println(string(out))
}
