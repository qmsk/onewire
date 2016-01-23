package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "encoding/json"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "strings"
    "time"
)

type APIHandler func(*http.Request, ...string) (interface{}, error)

type APIStatus struct {
    Name            string              `json:"name"`
    HidrawDevice    hidraw.DeviceInfo   `json:"hidraw_device"`
    AvrtempDevice   avrtemp.Status      `json:"avrtemp_device"`
    Stats           []string            `json:"stats"`
}

func (s *Server) GetStatus(_ *http.Request, path ...string) (interface{}, error) {
    var statusList []APIStatus

    for name, hidrawDevice := range s.hidrawDevices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   hidrawDevice,
        }

        if avrtempDevice := s.avrtempDevices[name]; avrtempDevice != nil {
            status.AvrtempDevice = avrtempDevice.Status()

            for statName, stat := range s.stats {
                if stat.Device == avrtempDevice {
                    status.Stats = append(status.Stats, statName)
                }
            }
        }

        statusList = append(statusList, status)
    }

    return statusList, nil
}

type APIStat struct {
    Stat        string      `json:"stat"`
    Time        time.Time   `json:"time"`
    Temperature float64     `json:"temperature"`
}

func (s *Server) GetStats(_ *http.Request, path ...string) (interface{}, error) {
    var statsList []APIStat

    for name, stat := range s.stats {
        apiStat := APIStat{
            Stat:           name,
            Time:           stat.Time,
            Temperature:    stat.Temperature.Float64(),
        }

        statsList = append(statsList, apiStat)
    }

    return statsList, nil
}

func (s *Server) lookupAPI(path []string) (APIHandler, []string) {
    switch path[0] {
    case "":
        return s.GetStatus, nil
    case "stats":
        return s.GetStats, path[1:]
    default:
        return nil, path
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := strings.Split(r.URL.Path, "/")
    if path[0] == "" {
        path = path[1:]
    }

    apiHandler, apiPath := s.lookupAPI(path)

    if apiHandler == nil {
        w.WriteHeader(404)

    } else if jsonData, err := apiHandler(r, apiPath...); err != nil {
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
