package server

type Server struct {
    // only modified at startup
    config          Config
    sensorConfig    map[string]*SensorConfig

    // state private to run()
    devices         map[string]*Device
    stats           map[string]Stat

    deviceChan          chan Device     // add/remove Device
    statChan            chan Stat       // read in from Devices
    influxChan          chan Stat       // write out to influx

    // HTTP API requests, handled by run()
    apiStatusChan       chan chan APIStatus
    apiStatChan         chan chan APIStat
}

func New() (*Server, error) {
    server := &Server{
        sensorConfig:   make(map[string]*SensorConfig),

        devices:        make(map[string]*Device),
        stats:          make(map[string]Stat),

        deviceChan:     make(chan Device),
        statChan:       make(chan Stat),

        apiStatusChan:  make(chan chan APIStatus),
        apiStatChan:    make(chan chan APIStat),
    }

    go server.run()

    return server, nil
}

func (s *Server) run() {
    var stat Stat

    for {
        select {
        case device := <-s.deviceChan:
            if device.avrtemp != nil {
                s.addDevice(&device)

            } else if runningDevice := s.devices[device.String()]; runningDevice != nil {
                s.removeDevice(runningDevice)
            }

        case stat = <-s.statChan:
            s.stat(stat)

        case statusChan := <-s.apiStatusChan:
            s.apiStatus(statusChan)
        case statChan := <-s.apiStatChan:
            s.apiStat(statChan)
        }
    }
}


