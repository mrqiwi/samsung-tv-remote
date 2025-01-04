package discover

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/koron/go-ssdp"
)

type DeviceDiscover struct {
	SearchType string
	WaitSec int
}

func NewDeviceDiscover(searchType string, waitSec int) *DeviceDiscover {
	return &DeviceDiscover{
		SearchType: searchType,
		WaitSec:    waitSec,
	}
}

type DeviceInfo struct {
	Name      string
	IPAddress string
}

func (d *DeviceDiscover) DiscoverSamsungTVs() ([]DeviceInfo, error) {
	var discoveredDevices []DeviceInfo

	services, err := ssdp.Search(d.SearchType, d.WaitSec, "")
	if err != nil {
		return discoveredDevices, err
	}

	for _, service := range services {
		device, err := discoverDevice(service)
		if err != nil {
			continue
		}

		discoveredDevices = append(discoveredDevices, device)
	}

	return discoveredDevices, nil
}

type UPnPDevice struct {
	Device struct {
		FriendlyName string `xml:"friendlyName"`
	} `xml:"device"`
}

func discoverDevice(service ssdp.Service) (DeviceInfo, error) {
	resp, err := http.Get(service.Location) // http://192.168.0.107:9197/dmr
	if err != nil {
		return DeviceInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DeviceInfo{}, errors.New("status code is not ok")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DeviceInfo{}, err
	}

	var upnpDevice UPnPDevice
	err = xml.Unmarshal(body, &upnpDevice)
	if err != nil {
		return DeviceInfo{}, err
	}

	parsedURL, err := url.Parse(service.Location)
	if err != nil {
		return DeviceInfo{}, err
	}

	return DeviceInfo{
		Name:      upnpDevice.Device.FriendlyName,
		IPAddress: strings.Split(parsedURL.Host, ":")[0],
	}, nil
}
