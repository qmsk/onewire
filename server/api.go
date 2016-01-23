package server

import (
    "fmt"
    "encoding/json"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "strings"
)

type APIHandler func(*http.Request, ...string) (interface{}, error)

type APIStatus struct {
    Name            string              `json:"name"`
    HidrawDevice    hidraw.DeviceInfo   `json:"hidraw_device"`
    AvrtempDevice   string              `json:"avrtemp_device"`
}


func (s *Server) GetStatus(_ *http.Request, path ...string) (interface{}, error) {
    var statusList []APIStatus

    for name, hidrawDevice := range s.hidrawDevices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   hidrawDevice,
        }

        if avrtempDevice := s.avrtempDevices[name]; avrtempDevice != nil {
            status.AvrtempDevice = avrtempDevice.String()
        }

        statusList = append(statusList, status)
    }

    return statusList, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := strings.Split(r.URL.Path, "/")
    if path[0] == "" {
        path = path[1:]
    }
    apiName := path[0]
    apiPath := path[1:]

    var apiHandler APIHandler

    switch apiName {
    case "":
        apiHandler = s.GetStatus
    default:
        w.WriteHeader(404)
        return
    }

    if jsonData, err := apiHandler(r, apiPath...); err != nil {
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
