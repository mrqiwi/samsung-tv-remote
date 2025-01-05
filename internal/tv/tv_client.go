package tv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"golang.org/x/exp/maps"
)

var (
	ErrUnauthorized    = errors.New("connection unauthorized")
	ErrUnexpectedEvent = errors.New("unexpected event type received")
	ErrDisconnected    = errors.New("TV disconnected")
)

type WSClient interface {
	SendMessage(msg []byte) error
	ReadMessage() ([]byte, error)
}

type TVClient struct {
	Client WSClient
	Cancel context.CancelFunc
}

func NewTVClient(client WSClient) (*TVClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &TVClient{
		Client: client,
		Cancel: cancel,
	}

	if err := c.handleEventMessage(); err != nil {
		return nil, err
	}

	go c.startListeningForMessages(ctx)

	return c, nil
}

type EventResponse struct {
	Event string `json:"event"`
	Data  struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (c *TVClient) handleEventMessage() error {
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
		return nil
	case "ms.channel.unauthorized":
		return ErrUnauthorized
	case "ms.channel.clientDisconnect":
		return ErrDisconnected
	default:
		return fmt.Errorf("%w: %s", ErrUnexpectedEvent, response.Event)
	}
}

func (c *TVClient) startListeningForMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.handleEventMessage()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (c *TVClient) Close() {
	c.Cancel()
}

type RemoteControlMessage struct {
	Method string              `json:"method"`
	Params RemoteControlParams `json:"params"`
}

type RemoteControlParams struct {
	Cmd          string `json:"Cmd"`
	DataOfCmd    string `json:"DataOfCmd"`
	Option       string `json:"Option,omitempty"`
	TypeOfRemote string `json:"TypeOfRemote,omitempty"`
}

func (c *TVClient) SendClickCommand(command string) error {
	message := RemoteControlMessage{
		Method: "ms.remote.control",
		Params: RemoteControlParams{
			Cmd:          "Click",
			DataOfCmd:    command,
			Option:       "false",
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
	"Mute":        "KEY_MUTE",     // Mute the sound
	"VolumeUp":    "KEY_VOLUP",    // Increase the volume
	"VolumeDown":  "KEY_VOLDOWN",  // Decrease the volume
	"ChannelUp":   "KEY_CHUP",     // Switch to the next channel
	"ChannelDown": "KEY_CHDOWN",   // Switch to the previous channel
	"PowerOff":    "KEY_POWEROFF", // Turn off the TV
	"Source":      "KEY_SOURCE",   // Switch the input source (e.g., HDMI, AV)
	"Home":        "KEY_HOME",     // Return to the home screen
	"Menu":        "KEY_MENU",     // Open the menu
	"Enter":       "KEY_ENTER",    // Press the "OK" button
	"Back":        "KEY_RETURN",   // Navigate back
	"ArrowUp":     "KEY_UP",       // Move the cursor up
	"ArrowDown":   "KEY_DOWN",     // Move the cursor down
	"ArrowLeft":   "KEY_LEFT",     // Move the cursor left
	"ArrowRight":  "KEY_RIGHT",    // Move the cursor right
	"Play":        "KEY_PLAY",     // Start media playback
	"Pause":       "KEY_PAUSE",    // Pause media playback
	"Stop":        "KEY_STOP",     // Stop media playback
	"Rewind":      "KEY_REWIND",   // Rewind media
	"FastForward": "KEY_FF",       // Fast-forward media
	"Info":        "KEY_INFO",     // Display information about the current playback
	"Exit":        "KEY_EXIT",     // Exit the current mode
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
