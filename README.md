# Samsung TV Remote

A command-line tool to control your Samsung TV via WebSocket.

## Features
- Discover Samsung TVs on the network using UPnP.
- Connect and authenticate with Samsung TVs via WebSocket.
- Send basic commands to control the TV, such as:
  - Mute/Unmute.
  - Adjust volume.
  - Power off (if supported).
  - Launch apps (optional in future versions).
- Cross-platform support: Works on Linux, macOS, and Windows.

## Installation
To install the tool, run:
```bash
go install github.com/mrqiwi/samsung-tv-remote/cmd/samsung-tv-remote@latest
```

## Usage
Once installed, you can use the tool via the command line. Here are some examples:

### Discover TVs on the network
```bash
samsung-tv-remote --discovery-timeout 10
```

### Specify a custom UPnP search target
```bash
samsung-tv-remote --search-target "urn:schemas-upnp-org:device:MediaRenderer:1"
```

## Troubleshooting
### TV not discovered
- Ensure your TV and computer are connected to the same network.
- Check if UPnP is enabled on your TV.
- Verify the `--search-target` matches your TV's UPnP type (e.g., `MediaRenderer`).

### Connection errors
- Check if the TV's IP address is accessible.
- Ensure the correct port is used (default: `8002`).
- If the connection is refused, the TV might need to enable remote control in its settings.

## Contributing
Contributions are welcome! To contribute:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request.

Ensure your code follows best practices and includes tests where appropriate.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.