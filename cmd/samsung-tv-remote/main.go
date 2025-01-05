package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	disc "github.com/mrqiwi/samsung-tv-remote/internal/discover"
	"github.com/mrqiwi/samsung-tv-remote/internal/tv"
	"github.com/mrqiwi/samsung-tv-remote/internal/ws"
)

func main() {
	var (
		port         int
		discTimeout  int
		searchTarget string
	)

	rootCmd := &cobra.Command{
		Use:   "samsung-tv-remote",
		Short: "Control your Samsung TV via WebSocket",
		Run: func(cmd *cobra.Command, args []string) {
			executeCommand(port, discTimeout, searchTarget)
		},
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 8002, "TV port number")
	rootCmd.Flags().StringVar(&searchTarget, "search-target", "urn:schemas-upnp-org:device:MediaRenderer:1", "UPnP search type")
	rootCmd.Flags().IntVar(&discTimeout, "discovery-timeout", 5, "Number of seconds to wait for device discovery")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error run: %v", err)
	}
}

func executeCommand(port, discTimeout int, searchTarget string) {
	fmt.Println("Searching for devices...")

	discover := disc.NewDeviceDiscover(searchTarget, discTimeout)
	device, err := chooseTV(discover)
	if err != nil {
		fmt.Print(err)
		return
	}

	tvURL := fmt.Sprintf("wss://%s:%d/api/v2/channels/samsung.remote.control", device.IPAddress, port)

	wsClient, err := ws.NewWebSocketClient(tvURL)
	if err != nil {
		fmt.Printf("Error connecting to the TV: %v", err)
		return
	}
	defer wsClient.Close()

	fmt.Println("Attempting to connect to the TV. Please approve the connection request on your TV screen...")

	tvClient, err := tv.NewTVClient(wsClient)
	if err != nil {
		fmt.Printf("Error tv authenticating: %v", err)
		return
	}
	defer tvClient.Close()

	for {
		cmd, err := chooseTVCommand(tvClient)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = tvClient.ExecuteCommand(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Command sent successfully")
	}
}

func chooseTV(discover *disc.DeviceDiscover) (disc.DeviceInfo, error) {
	devices, err := discover.DiscoverSamsungTVs()
	if err != nil || len(devices) == 0 {
		return disc.DeviceInfo{}, fmt.Errorf("No TVs found on the network")
	}

	fmt.Println("Discovered TVs:")
	for i, device := range devices {
		fmt.Printf("%d. %s (%s)\n", i+1, device.Name, device.IPAddress)
	}

	fmt.Printf("Enter the number of the TV you want to connect to: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return disc.DeviceInfo{}, fmt.Errorf("Error reading input: %v", err)
	}

	selection, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || selection < 1 || selection > len(devices) {
		return disc.DeviceInfo{}, fmt.Errorf("Invalid selection")
	}

	return devices[selection-1], nil
}

func chooseTVCommand(tvClient *tv.TVClient) (string, error) {
	commands := tvClient.AvailableCommands()

	fmt.Println("Available commands:")
	for i, cmd := range commands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}

	fmt.Printf("Enter the number of the command to execute: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Error reading input: %v", err)
	}

	selection, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || selection < 1 || selection > len(commands) {
		return "", fmt.Errorf("Invalid selection")
	}

	return commands[selection-1], nil
}
