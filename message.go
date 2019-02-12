package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type UrbanDictionary struct {
	List []SomeThingy
}

type SomeThingy struct {
	Definition   string
	Permalink    string
	Thumbs_up    int
	Sound_urls   []string
	Author       string
	Word         string
	Defid        int
	Current_vote string
	Written_on   string
	Example      string
	Thumbs_down  int
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	words := strings.Split(strings.ToLower(m.Content), " ")
	if words[0] != "!def" {
		return
	}
	if len(words) == 1 {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> Could not process the request!")
		er.Println("Something wrong with message sent by: " + m.Author.Username)
		return
	}
	resp, err := http.Get("http://api.urbandictionary.com/v0/define?term=" + string(words[1]))
	if err != nil {
		er.Println("Could not make the request to the api!", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		er.Println("The request to urban dictionary failed! The GET did not respond with statuscode 200.")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		er.Println("Failed to read response body! ", err.Error())
		return
	}
	var result UrbanDictionary
	err = json.Unmarshal(body, &result)
	if err != nil {
		er.Println("Error decoding json: ", err.Error())
		return
	}
	if len(result.List) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Thats not a word <@"+m.Author.ID+">!")
		er.Println("The requested word is not defined in urban dictionary!")
		return
	}
	urban := result.List[0]
	s.ChannelMessageSend(m.ChannelID, "Word: "+urban.Word+"\n"+"Definition: "+urban.Definition+"\n"+"Link: "+urban.Permalink)
	suc.Println("Definition processed and sent!")
}
