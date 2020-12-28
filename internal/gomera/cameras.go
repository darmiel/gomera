package gomera

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

func (camera *Camera) isCamera(cameraId string, secret string) (res bool) {
	return strings.EqualFold(camera.Id, cameraId) && camera.Secret == secret
}

func (camera *Camera) createJsonPayload() (payload string) {
	var timeStr = time.Now().Format("02.01.2006 15:04:05")

	return `
{
"content": "@everyone",
"username": "Notify for ` + camera.Name + `",
"avatar_url": "` + camera.Avatar + `",
"embeds": [
    {
      "title": "👉 View Stream",
      "description": "**Detected motion**",
      "url": "` + camera.StreamUrl + `",
      "color": 16725044,
      "fields": [
        {
          "name": "📸",
          "value": "` + camera.Name + `",
          "inline": true
        },
        {
          "name": "📷",
          "value": "` + camera.Id + `",
          "inline": true
        },
        {
          "name": "⏰",
          "value": "` + timeStr + `",
          "inline": false
        }
      ],
      "author": {
        "name": "` + camera.Id + `/` + camera.Name + `",
        "icon_url": "` + camera.Avatar + `"
      }
    }
  ]
}`
}

func (camera *Camera) sendWebhook() (r *http.Response, err error) {
	// send http request to discord webhook
	if len(camera.Webhook) == 0 {
		return nil, errors.New("webhook was empty")
	}

	log.Println("Creating payload")
	reader := bytes.NewReader([]byte(camera.createJsonPayload()))
	req, err := http.NewRequest("POST", camera.Webhook, reader)
	if err != nil {
		return nil, err
	}
	log.Println("Created request:", req)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)

	log.Println("Client did ... things")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp, nil
}
