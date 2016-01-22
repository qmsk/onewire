package avrtemp

/*
 * Support for the Diamex GmbH / AVR Temp Sensor.
 *
 * This is a simple DS18x20 sensor reader module using an Atmel AT90USB162 (based on the Teensy devboard).
 *
 * The board enumerates the 1wire bus for DS18x20 sensors at startup, and then sends a report every second, cycling through the
 * enumerated sensors.
 */

import (
    "encoding/binary"
    "bytes"
    "fmt"
    "encoding/hex"
    "github.com/qmsk/onewire/hidraw"
)

// hidraw device configuration for hidraw.Find()
var HIDRAW_CONFIG = hidraw.DeviceConfig{
    Vendor:     0x16C0,
    Product:    0x0480,
}

type Device struct {
    hidrawDevice    *hidraw.Device
}

type ID             [8]byte

func (self ID) String() string {
    return hex.EncodeToString(self[:])
}

type Temperature    uint16

func (self Temperature) Float64() float64 {
    return float64(self) / 10.0
}

type PowerState     uint8

func (self PowerState) String() string {
    switch self {
    case 0x00:
        return "parasite"
    case 0x01:
        return "extern"
    default:
        return fmt.Sprintf("%#04x", self)
    }
}

type Report struct {
    Count       uint8           // 0x00
    Index       uint8           // 0x01
    Power       PowerState      // 0x02
    _           [1]byte
    Temp        Temperature     // 0x04
    _           [2]byte         // 0x06
    ID          ID              // 0x08
}

func (self *Report) unpack(buf []byte) error {
    return binary.Read(bytes.NewReader(buf), binary.LittleEndian, self)
}

func (self Report) String() string {
    return fmt.Sprintf("Sensor #%d of %d: %.1fC (Power: %v ID: %v)",
        self.Index, self.Count,
        self.Temp.Float64(),
        self.Power, self.ID,
    )
}

func Open(hidrawDevice *hidraw.Device) (*Device, error) {
    device := &Device{
        hidrawDevice:   hidrawDevice,
    }

    return device, nil
}

func (self *Device) Read() (report Report, err error) {
    buf := make([]byte, 64)

    if readSize, err := self.hidrawDevice.Read(buf); err != nil {
        return report, err
    } else {
        buf = buf[:readSize]
    }

    if err := report.unpack(buf); err != nil {
        return report, err
    }

    return
}
