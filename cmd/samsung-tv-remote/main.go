package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	disc "github.com/mrqiwi/samsung-tv-remote/internal/discover"
	"github.com/mrqiwi/samsung-tv-remote/internal/tui"
	"github.com/mrqiwi/samsung-tv-remote/internal/tv"
	"github.com/mrqiwi/samsung-tv-remote/internal/ws"
	"github.com/spf13/cobra"
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

	commands := tvClient.AvailableCommands()
	model := tui.NewCommandListModel(commands)

	for {
		err = chooseTVCommand(&model)
		if err != nil {
			fmt.Println(err)
			return
		}

		err := tvClient.ExecuteCommand(model.SelectedCommand())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func chooseTV(discover *disc.DeviceDiscover) (disc.DeviceInfo, error) {
	devices, err := discover.DiscoverSamsungTVs()
	if err != nil || len(devices) == 0 {
		return disc.DeviceInfo{}, fmt.Errorf("No TVs found on the network")
	}

	model := tui.NewDeviceListModel(devices)
	p := tea.NewProgram(&model)

	if _, err := p.Run(); err != nil {
		return disc.DeviceInfo{}, fmt.Errorf("Error running program: %v", err)
	}

	if model.Esc() {
		return disc.DeviceInfo{}, fmt.Errorf("Exit running program")
	}

	return model.SelectedDevice(), nil
}

func chooseTVCommand(model *tui.CommandListModel) error {
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error running program: %v", err)
	}

	if model.Esc() {
		return fmt.Errorf("Exit running program")
	}

	return nil
}
