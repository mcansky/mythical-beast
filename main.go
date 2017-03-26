package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"os"
	"time"
)

type mergeData struct {
	Id          string `json:"id"`
	Action      string `json:"action"`
	PullRequest struct {
		Merged bool `json:"merged"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		Repo struct {
			FullName string `json:"full_name"`
			Sha      string `json:"sha"`
		} `json:"repo"`
	} `json:"pull_request"`
}

type responseData struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var ChannelID = os.Getenv("DISCORD_CHANNEL_ID")
var BadRequest = "Somebody tried something on %s, but I couldn't deal with it"
var Token = os.Getenv("DISCORD_SECRET")
var BotID string

// utils

func console_logger(level string, topic string, message string) {
	fmt.Println("%s: %20s %20s %s > %s", level, ChannelID, time.Now().Format(time.Stamp), topic, message)
}

// Concerns

func shape_message(data mergeData) (string, error) {
	if data.Id != "" {
		msg := fmt.Sprintf("Github Deploy for %s:%s (@%s)", data.PullRequest.Repo.FullName, data.PullRequest.Repo.Sha, data.PullRequest.User.Login)
		return msg, nil
	} else {
		return "", fmt.Errorf("missing data in github payload")
	}
}

// callbacks
func basics(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

}

// controllers

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
	fmt.Println("%20s %20s %20s > %s", ChannelID, time.Now().Format(time.Stamp), "Received query", "/hello")
	DiscordMessage("", "Hello !!!")
}

func github_deploy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Github deploy request received !")
	decoder := json.NewDecoder(r.Body)
	var github_data mergeData
	err := decoder.Decode(&github_data)
	if err != nil {
		console_logger("ERROR", "github deploy", "Could not decode properly github_deploy request")
	}
	var response responseData
	msg, err := shape_message(github_data)
	if err == nil {
		console_logger("INFO", "github deploy", msg)
		DiscordMessage("", ":satellite_orbital: "+msg)
		response.Message = "ok"
		response.Success = true
	} else {
		fmt.Println("Could not read github data: %s", err)
		msg = fmt.Sprintf(BadRequest, "/github_deploy")
		DiscordMessage("", "<:megaphone:295327332858593280> "+msg)
		response.Message = "ko"
		response.Success = false
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func main() {
	register_bot()
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/github_deploy", github_deploy)

	http.ListenAndServe(":8080", nil)
}
