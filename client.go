package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/koo04/kdeck-server/api"
	"github.com/koo04/kdeck-server/proto/data"
	"github.com/wailsapp/wails"
	"google.golang.org/grpc"
)

type ServerSettings struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type Client struct {
	log     *wails.CustomLogger
	r       *wails.Runtime
	conn    *grpc.ClientConn
	client  data.DataServiceClient
	buttons *wails.Store
	Server  ServerSettings `json:"server"`
}

func (c *Client) WailsInit(runtime *wails.Runtime) error {
	c.r = runtime
	c.log = runtime.Log.New("client")
	c.buttons = runtime.Store.New("Buttons", []api.Button{})

	c.LoadSettings()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", c.Server.Url, strconv.Itoa(c.Server.Port)), grpc.WithInsecure())
	// conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		c.log.Errorf("Could not connect: %s", err)
	}
	c.conn = conn
	client := data.NewDataServiceClient(c.conn)
	c.client = client

	var connection = make(chan bool)
	go func() {
		for {
			if c.conn.GetState().String() == "READY" {
				c.UpdateButtons()
				c.r.Events.Emit("ready")
				connection <- true
			} else {
				c.log.Error("We don't have a connection to the server")
				c.r.Events.Emit("error")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			if c.conn.GetState().String() != "READY" {
				c.log.Debug("Server is not ready")
				c.r.Events.Emit("ready", false)
				time.Sleep(1 * time.Second)
			} else {
				c.log.Debug("Server ready!")
				c.r.Events.Emit("ready", true)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	return nil
}

func (c *Client) WailsShutdown() {
	defer c.conn.Close()
}

func (c *Client) LoadSettings() {
	c.log.Debug("Trying to get the working directory")
	workDir, err := os.Getwd()
	if err != nil {
		c.log.Fatalf("Could not get work directory: %s", err)
	}

	c.log.Debug("Checking to see if the settings file exists")
	settingsFile, err := os.Open(fmt.Sprintf("%s/settings.json", workDir))
	if err != nil {
		c.log.Debugf("Could not load settings: %s", err)
		c.log.Debug("Creating default settings")
		c.Server = ServerSettings{}
		c.Server.Url = "localhost"
		c.Server.Port = 9000

		c.SaveSettings()
	} else {
		c.log.Debug("Found an existing settings file. Loading...")
		jsonParser := json.NewDecoder(settingsFile)
		if err = jsonParser.Decode(&c); err != nil {
			c.log.Fatalf("Could not parse settings: %s", err)
		}

		j, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			c.log.Fatalf("Could not marshal default settings: %s", err)
		}
		c.log.Debugf("... Loaded: %s", j)
	}
}

func (c *Client) SaveSettings() {
	workDir, err := os.Getwd()
	if err != nil {
		c.log.Fatalf("Could not get work directory: %s", err)
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		c.log.Fatalf("Could not marshal default settings: %s", err)
	}

	defer os.WriteFile(fmt.Sprintf("%s/settings.json", workDir), j, 0777)
	if err != nil {
		c.log.Fatalf("Could not save settings file")
	}
}

func (c *Client) UpdateButtons() {
	if c.conn.GetState().String() == "READY" {
		response, err := c.client.GetButtons(context.Background(), &data.Empty{})
		if err != nil {
			c.log.Fatalf("Error updating buttons: %s", err)
		}
		var buttons []api.Button
		json.Unmarshal([]byte(response.Body), &buttons)
		c.log.Debugf("Update Buttons: %s", buttons)

		c.buttons.Set(buttons)
	}
}

func (c *Client) GetButtons() string {
	if c.conn.GetState().String() == "READY" {
		buttons, _ := json.Marshal(c.buttons.Get())
		c.log.Debugf("Get Buttons: %s", string(buttons))
		return string(buttons)
	} else {
		return "[]"
	}
}

func (c *Client) PressButton(plugin string, action string) {
	c.log.Debugf("Plugin: %s | Action: %s", plugin, action)
	_, err := c.client.PressButton(context.Background(), &data.PressButtonRequest{Plugin: plugin, Action: action})
	if err != nil {
		c.log.Fatalf("Error: %s", err)
	}
}
