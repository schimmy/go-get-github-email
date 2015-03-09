package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	user := os.Args[1]
	fmt.Printf("getting email for username: %s\n", user)
	userEndpoint := fmt.Sprintf("https://api.github.com/users/%s", user)
	resp, err := http.Get(userEndpoint)
	if err != nil {
		log.Fatal("error connecting to github!")
	}
	defer resp.Body.Close()
	var userData map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	for {
		if err := dec.Decode(&userData); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	name := userData["name"]

	fmt.Printf("Full name is: %s\n", name)
	fmt.Printf("Profile email is: %s\n", userData["email"])

	eventEndpoint := fmt.Sprintf("https://api.github.com/users/%s/events", user)
	resp, err = http.Get(eventEndpoint)
	if err != nil {
		log.Fatal("error connecting to github!")
	}
	defer resp.Body.Close()
	var dat []interface{}
	for {
		dec = json.NewDecoder(resp.Body)
		if err := dec.Decode(&dat); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	//fmt.Println(dat)
	for _, e := range dat {
		event := e.(map[string]interface{})
		// get only pushes
		if t := event["type"]; t == "PushEvent" {
			payload := event["payload"].(map[string]interface{})
			commits := payload["commits"].([]interface{})
			for _, c := range commits {
				a := c.(map[string]interface{})["author"]
				author := a.(map[string]interface{})
				// find where author matches the name of our user -
				// there could be other people's commits in a merged PR
				if author["name"] == name {
					fmt.Printf("Commit email is: %s\n", author["email"])
					return
				}
			}
		}
	}
	log.Fatal("couldn't find a commit with a name that matches our user!")
}
