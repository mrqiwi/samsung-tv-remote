package tv

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"golang.org/x/exp/maps"
)

var (
	ErrTokenMissing      = errors.New("Token is missing")
	ErrUnauthorized      = errors.New("Connection unauthorized")
	ErrUnexpectedEvent   = errors.New("Unexpected event type received")
)

type WSClient interface {
	SendMessage(msg []byte) error
	ReadMessage() ([]byte, error)
}

type TVClient struct {
	Client WSClient
	Token  string
}

func NewTVClient(client WSClient) *TVClient {
	return &TVClient{
		Client: client,
	}
}

func (c *TVClient) Authenticate() error {
	if err := c.GetToken(); err != nil {
		return fmt.Errorf("get token failed: %w", err)
	}

	if err := c.RegisterToken(); err != nil {
		return fmt.Errorf("register token failed: %w", err)
	}

	return nil
}

type EventResponse struct {
	Event string `json:"event"`
	Data  struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (c *TVClient) GetToken() error {
	message, err := c.Client.ReadMessage()
	if err != nil {
		return err
	}

	var response EventResponse

	if err := json.Unmarshal(message, &response); err != nil {
		return err
	}

	switch response.Event {
	case "ms.channel.connect":
		if response.Data.Token == "" {
			return ErrTokenMissing
		}
		c.Token = response.Data.Token
		return nil
	case "ms.channel.unauthorized":
		return ErrUnauthorized
	default:
		return fmt.Errorf("%w: %s", ErrUnexpectedEvent, response.Event)
	}
}

type RemoteControlMessage struct {
	Method string              `json:"method"`
	Params RemoteControlParams `json:"params"`
}

type RemoteControlParams struct {
	Cmd         string `json:"Cmd"`
	DataOfCmd   string `json:"DataOfCmd"`
	Option      string `json:"Option,omitempty"`
	TypeOfRemote string `json:"TypeOfRemote,omitempty"`
}

func (c *TVClient) RegisterToken() error {
	message := RemoteControlMessage{
		Method: "ms.remote.control",
		Params: RemoteControlParams{
			Cmd:       "Register",
			DataOfCmd: fmt.Sprintf(`{"auth_Type":"token","ClientToken":"%s"}`, c.Token),
		},
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.Client.SendMessage(bytes)
}

func (c *TVClient) SendClickCommand(command string) error {
	message := RemoteControlMessage{
		Method: "ms.remote.control",
		Params: RemoteControlParams{
			Cmd:         "Click",
			DataOfCmd:   command,
			Option:      "false",
			TypeOfRemote: "SendRemoteKey",
		},
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.Client.SendMessage(bytes)
}

var commandMap = map[string]string{
	"Mute":        "KEY_MUTE",        // Mute the sound
	"VolumeUp":    "KEY_VOLUP",       // Increase the volume
	"VolumeDown":  "KEY_VOLDOWN",     // Decrease the volume
	"ChannelUp":   "KEY_CHUP",        // Switch to the next channel
	"ChannelDown": "KEY_CHDOWN",      // Switch to the previous channel
	"PowerOff":    "KEY_POWEROFF",    // Turn off the TV
	"Source":      "KEY_SOURCE",      // Switch the input source (e.g., HDMI, AV)
	"Home":        "KEY_HOME",        // Return to the home screen
	"Menu":        "KEY_MENU",        // Open the menu
	"Enter":       "KEY_ENTER",       // Press the "OK" button
	"Back":        "KEY_RETURN",      // Navigate back
	"ArrowUp":     "KEY_UP",          // Move the cursor up
	"ArrowDown":   "KEY_DOWN",        // Move the cursor down
	"ArrowLeft":   "KEY_LEFT",        // Move the cursor left
	"ArrowRight":  "KEY_RIGHT",       // Move the cursor right
	"Play":        "KEY_PLAY",        // Start media playback
	"Pause":       "KEY_PAUSE",       // Pause media playback
	"Stop":        "KEY_STOP",        // Stop media playback
	"Rewind":      "KEY_REWIND",      // Rewind media
	"FastForward": "KEY_FF",          // Fast-forward media
	"Info":        "KEY_INFO",        // Display information about the current playback
	"Exit":        "KEY_EXIT",        // Exit the current mode
}

func (c *TVClient) AvailableCommands() []string {
	commands := maps.Keys(commandMap)

	sort.Strings(commands)

	return commands
}

func (c *TVClient) ExecuteCommand(command string) error {
	keyCode, exists := commandMap[command]
	if !exists {
		return fmt.Errorf("Unknown command: %s", command)
	}
	return c.SendClickCommand(keyCode)
}

