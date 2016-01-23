package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "encoding/json"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "log"
)

type APIStatus struct {
    Name            string  `json:"name"`
    HidrawDevice    string  `json:"hidraw_device"`
    AvrtempDevice   string  `json:"avrtemp_device"`
}

type Server struct {
    hidrawDevices   map[string]*hidraw.Device
    avrtempDevices  map[string]*avrtemp.Device
}

func New() (*Server, error) {
    server := &Server{
        hidrawDevices:  make(map[string]*hidraw.Device),
        avrtempDevices: make(map[string]*avrtemp.Device),
    }

    return server, nil
}

func (s *Server) AddHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
        log.Printf("AddHidrawDevice %#v: hidraw.Open: %v\n", deviceInfo, err)
    } else {
        s.hidrawDevices[deviceInfo.String()] = hidrawDevice

        if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
            log.Printf("AddHidrawDevice %#v: avrtemp.Open: %v\n", deviceInfo, err)
        } else {
            log.Printf("AddHidrawDevice %#v: %#v\n", deviceInfo, avrtempDevice)

            s.avrtempDevices[deviceInfo.String()] = avrtempDevice
        }
    }
}

func (s *Server) RemoveHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    log.Printf("RemoveHidrawDevice %v...\n", deviceInfo)

    if avrtempDevice := s.avrtempDevices[deviceInfo.String()]; avrtempDevice != nil {
        avrtempDevice.Close()
    } else if hidrawDevice := s.hidrawDevices[deviceInfo.String()]; hidrawDevice != nil {
        hidrawDevice.Close()
    }

    delete(s.hidrawDevices, deviceInfo.String())
    delete(s.avrtempDevices, deviceInfo.String())
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

func (s *Server) GetStatus() (interface{}, error) {
    var statusList []APIStatus

    for name, hidrawDevice := range s.hidrawDevices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   hidrawDevice.String(),
        }

        if avrtempDevice := s.avrtempDevices[name]; avrtempDevice != nil {
            status.AvrtempDevice = avrtempDevice.String()
        }

        statusList = append(statusList, status)
    }

    return statusList, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if jsonData, err := s.GetStatus(); err != nil {
        w.WriteHeader(500)

        fmt.Fprintf(w, "%v\n", err)
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)

        if err := json.NewEncoder(w).Encode(jsonData); err != nil {
            panic(err)
        }
    }
}
