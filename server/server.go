package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "github.com/qmsk/onewire/hidraw"
    "log"
)

type Server struct {
    hidrawDevices   map[string]hidraw.DeviceInfo
    avrtempDevices  map[string]*avrtemp.Device
}

func New() (*Server, error) {
    server := &Server{
        hidrawDevices:  make(map[string]hidraw.DeviceInfo),
        avrtempDevices: make(map[string]*avrtemp.Device),
    }

    return server, nil
}

func (s *Server) AddHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    s.hidrawDevices[deviceInfo.String()] = deviceInfo

    if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
        log.Printf("AddHidrawDevice %#v: hidraw.Open: %v\n", deviceInfo, err)
    } else if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
        log.Printf("AddHidrawDevice %#v: avrtemp.Open: %v\n", deviceInfo, err)
    } else {
        log.Printf("AddHidrawDevice %#v: %#v\n", deviceInfo, avrtempDevice)

        s.avrtempDevices[deviceInfo.String()] = avrtempDevice
    }
}

func (s *Server) RemoveHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    log.Printf("RemoveHidrawDevice %v...\n", deviceInfo)

    if avrtempDevice := s.avrtempDevices[deviceInfo.String()]; avrtempDevice != nil {
        avrtempDevice.Close()
    }

    delete(s.avrtempDevices, deviceInfo.String())
    delete(s.hidrawDevices, deviceInfo.String())
}

func (s *Server) MonitorHidraw(monitorChan chan hidraw.MonitorEvent) {
    for monitorEvent := range monitorChan {
        switch monitorEvent.Action {
        case "add":
            s.AddHidrawDevice(monitorEvent.DeviceInfo)
        case "remove":
            s.RemoveHidrawDevice(monitorEvent.DeviceInfo)
        default:
            log.Printf("MonitorHidraw: %v?! %v\n", monitorEvent.Action, monitorEvent.DeviceInfo)
        }
    }

    log.Printf("server.Server %v: MonitorHidraw: exit\n", s)
}
