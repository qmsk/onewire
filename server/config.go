package server

import (
    "fmt"
    "github.com/BurntSushi/toml"
)

type SensorConfig struct {
    name        string
    ID          string
}

func (self SensorConfig) String() string {
    return self.name
}

type Config struct {
    Sensors     map[string]*SensorConfig
}

func (s *Server) LoadConfig(filePath string) error {
    if meta, err := toml.DecodeFile(filePath, &s.config); err != nil {
        return err
    } else if len(meta.Undecoded()) > 0 {
        return fmt.Errorf("Undecoded keys: %v", meta.Undecoded())
    }

    for sensorName, sensorConfig := range s.config.Sensors {
        sensorConfig.name = sensorName

        s.sensorConfig[sensorConfig.ID] = sensorConfig
    }

    return nil
}
