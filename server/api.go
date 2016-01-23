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

func (s *Server) GetConfig(_ *http.Request, path ...string) (interface{}, error) {
    return s.config, nil
}

type APIStatus struct {
    Name            string              `json:"name"`
    HidrawDevice    hidraw.DeviceInfo   `json:"hidraw"`
    AvrtempDevice   avrtemp.Status      `json:"avrtemp"`
    Stats           map[string]string   `json:"stats"`
}

func (s *Server) GetStatus(_ *http.Request, path ...string) (interface{}, error) {
    var statusList []APIStatus

    for name, device  := range s.devices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   device.hidraw,
            Stats:          make(map[string]string),
        }

        if device.avrtemp != nil {
            status.AvrtempDevice = device.avrtemp.Status()
        }

        for statID, stat := range s.stats {
            if stat.Device != device {
                continue
            } else {
                // nil-safe
                status.Stats[statID] = s.sensorConfig[statID].String()
            }
        }

        statusList = append(statusList, status)
    }

    return statusList, nil
}

type APIStat struct {
    ID          string      `json:"id"`
    Family      string      `json:"family"`
    SensorName  string      `json:"sensor_name"`
    Time        time.Time   `json:"time"`
    Temperature float64     `json:"temperature"`
}

func (s *Server) GetStats(_ *http.Request, path ...string) (interface{}, error) {
    var statsList []APIStat

    for id, stat := range s.stats {
        apiStat := APIStat{
            ID:             id,
            Family:         stat.ID.Family(),
            Time:           stat.Time,
            Temperature:    stat.Temperature.Float64(),
        }

        if stat.SensorConfig != nil {
            apiStat.SensorName = stat.SensorConfig.String()
        }

        statsList = append(statsList, apiStat)
    }

    return statsList, nil
}

func (s *Server) lookupAPI(path []string) (APIHandler, []string) {
    switch path[0] {
    case "":
        return s.GetStatus, nil
    case "config":
        return s.GetConfig, nil
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
