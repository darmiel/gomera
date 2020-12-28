package gomera

import (
	"encoding/json"
	"errors"
	"gomera/internal/discord"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	Cameras = make(map[string]*Camera)
)

type Camera struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Webhook   string `json:"webhook"`
	Avatar    string `json:"avatar"`
	StreamUrl string `json:"stream_url"`
	Secret    string `json:"secret"`
}

func parseCameras() (res map[string]*Camera, err error) {
	data, err := ioutil.ReadFile(opt.CameraFile)
	if err != nil {
		return nil, err
	}

	var cameras []*Camera
	err = json.Unmarshal(data, &cameras)
	if err != nil {
		return nil, err
	}

	// clear old cameras
	Cameras = make(map[string]*Camera)
	for _, camera := range cameras {
		Cameras[camera.Id] = camera
	}

	return Cameras, nil
}

func FindCamera(cameraId string, secret string) (cam *Camera, err error) {
	for _, camera := range Cameras {
		if camera.IsCamera(cameraId, secret) {
			cam = camera
			return cam, nil
		}
	}
	return nil, errors.New("camera not found or access denied")
}

func (camera *Camera) IsCamera(cameraId string, secret string) (res bool) {
	return strings.EqualFold(camera.Id, cameraId) && camera.Secret == secret
}

func (camera *Camera) CreatePayload() (message *discord.WebhookMessage) {
	datestr := time.Now().Format("02.01.2006 15:04:05")

	var color uint
	var content string

	if !opt.DevEnvironment {
		color = 16725044      // light red
		content = "@everyone" // ping everyone
	} else {
		color = 3971831               // blue
		content = "**Dev/Test-Mode**" // ping nobody
	}

	message = &discord.WebhookMessage{
		Content:   content,
		Username:  "Notify for " + camera.Name,
		AvatarUrl: camera.Avatar,
		Embeds: []*discord.Embed{
			{
				Title:       "👉 View Stream",
				Description: "**Detected motion**",
				Url:         camera.StreamUrl,
				Color:       color,
				Fields: []*discord.Field{
					discord.NewFieldInline("📸", camera.Name),
					discord.NewFieldInline("📷", camera.Id),
					discord.NewField("⏰", datestr),
				},
				Author: &discord.Author{
					Name:    camera.Id + `/` + camera.Name,
					IconUrl: camera.Avatar,
				},
			},
		},
	}

	return message
}

func (camera *Camera) Send() (r *http.Response, err error) {
	message := camera.CreatePayload()
	r, err = message.Send(camera.Webhook)
	return r, err
}
