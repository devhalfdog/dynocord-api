package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/devhalfdog/dynocord-api/utils"
)

var (
	WebhookID    = utils.Environment("DISCORD_ID")
	WebhookToken = utils.Environment("DISCORD_TOKEN")
)

func UploadImage(filepath string, streamer string) (string, error) {
	out, err := os.Open(filepath)
	if err != nil {
		log.Println("upload image error :", err)
		return "", err
	}
	defer out.Close()

	param := &discordgo.WebhookParams{
		Files: []*discordgo.File{
			{
				Name:   fmt.Sprintf("%s.png", streamer),
				Reader: out,
			},
		},
	}

	client, err := discordgo.New("")
	if err != nil {
		log.Println("discordgo instance error :", err)
		return "", err
	}

	p, err := client.WebhookExecute(WebhookID, WebhookToken, true, param)
	if err != nil {
		log.Println("execute webhook error :", err)
		return "", err
	}

	return p.Attachments[0].URL, nil

}
