package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func NewFieldInline(name string, value string) (res *Field) {
	return &Field{
		Name:   name,
		Value:  value,
		Inline: true,
	}
}

func NewField(name string, value string) (res *Field) {
	return &Field{
		Name:   name,
		Value:  value,
		Inline: false,
	}
}

///////////////////////////////////////////////////////////////////////////

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type Embed struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Color       uint     `json:"color"`
	Fields      []*Field `json:"fields"`
	Author      *Author  `json:"author"`
}

type WebhookMessage struct {
	Content   string   `json:"content"`
	Username  string   `json:"username"`
	AvatarUrl string   `json:"avatar_url"`
	Embeds    []*Embed `json:"embeds"`
}

func (m *WebhookMessage) Marshall() (res string, err error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *WebhookMessage) Send(u string) (r *http.Response, err error) {
	if len(u) <= 0 {
		return nil, errors.New("webhook-url was empty")
	}

	js, err := m.Marshall()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader([]byte(js))
	req, err := http.NewRequest("POST", u, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
