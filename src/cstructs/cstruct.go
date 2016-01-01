package cstructs

import (
	"os"
)

type Container struct {
	Name   string
	Status string
	Time   string
	Log    string
}

type Message struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Color   string  `json:"color"`
	Text    string  `json:"text"`
	Pretext string  `json:"pretext"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
