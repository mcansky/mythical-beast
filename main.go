package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"os"
	"time"
	"encoding/json"
)

type mergeData struct {
	Id string `json:"id"`
	Action string `json:"action"`
	PullRequest struct {
		Merged bool `json:"merged"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Repo struct {
			FullName string `json:"full_name"`
			Sha string `json:"sha"`
		} `json:"repo"`
	} `json:"pull_request"`
}

type responseData struct {
	Id string `json:"id"`
	Success bool `json:"success"`
	Message string `json:"message"`
}

var ChannelID = os.Getenv("DISCORD_CHANNEL_ID")
var BadRequest = "Somebody tried something on %s, but I couldn't deal with it"

// utils

func console_logger(level string, topic string, message string) {
	fmt.Println("[%s] %20s %20s %20s > %s", level, ChannelID, time.Now().Format(time.Stamp), topic, message)
}


// Discord messenger

func DiscordMessage(channel_name string, message string) {
	var Token = os.Getenv("DISCORD_SECRET")
	var ChannelID = os.Getenv("DISCORD_CHANNEL_ID")

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}

	// Get the account information.
	u, err := discord.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

    _, err = discord.ChannelMessageSend(ChannelID, message)
	if err != nil {
		fmt.Println("could not send message to discord,", err)
	}
	// Register messageCreate as a callback for the messageCreate events.
	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	return
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


// controllers

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
	fmt.Println("%20s %20s %20s > %s", ChannelID, time.Now().Format(time.Stamp), "Received query", "/hello")
	DiscordMessage("", "Hello")
}

func github_deploy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Github deploy request received !")
	decoder := json.NewDecoder(r.Body)
	var github_data mergeData
	err := decoder.Decode(&github_data)
	if err != nil {
		fmt.Println("Could not decode properly github_deploy request")
	}
	var response responseData
	msg, err := shape_message(github_data)
	if err == nil {
		fmt.Println(msg)
		DiscordMessage("", ":satellite_orbital: " + msg)
		response.Message = "ok"
		response.Success = true
	} else {
		fmt.Println("Could not read github data: %s", err)
		msg = fmt.Sprintf(BadRequest, "/github_deploy")
		DiscordMessage("", "<:megaphone:295327332858593280> " + msg)
		response.Message = "ko"
		response.Success = false
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}


func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/github_deploy", github_deploy)

	http.ListenAndServe(":8080", nil)
}
