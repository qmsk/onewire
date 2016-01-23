package main

import (
    "github.com/qmsk/onewire/avrtemp"
    "flag"
    "encoding/hex"
    "github.com/qmsk/onewire/hidraw"
    "log"
    "os"
)

var (
    deviceConfig hidraw.DeviceConfig
)

func init() {
    flag.UintVar(&deviceConfig.VendorID, "device-vendor", avrtemp.HIDRAW_CONFIG.VendorID,
        "Select device vendor")
    flag.UintVar(&deviceConfig.ProductID, "device-product", avrtemp.HIDRAW_CONFIG.ProductID,
        "Select device product")
}

func hexdump(buf []byte) {
    dumper := hex.Dumper(os.Stdout)

    if _, err := dumper.Write(buf); err != nil {
        panic(err)
    }

    dumper.Close()
}

func main() {
    flag.Parse()

    if hidrawList, err := hidraw.List(deviceConfig); err != nil {
        log.Fatalf("hidraw.List %v: %v\n", deviceConfig, err)
    } else {
        log.Printf("hidraw.List...\n")
        for _, deviceInfo := range hidrawList {
            if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
                log.Printf("hidraw.Open: %v\n", deviceInfo, err)
            } else if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
                log.Printf("avrtemp.Open: %v\n", deviceInfo, err)
            } else {
                log.Printf("%#v\n", deviceInfo, avrtempDevice)

                log.Printf("avrtempDevice %v: Read...\n", avrtempDevice)
                for {
                    if report, err := avrtempDevice.Read(); err != nil {
                        log.Fatalf("avrtemp.Device %v: Read: %v\n", avrtempDevice, err)
                    } else {
                        log.Printf("%v\n", report)
                    }
                }
            }
        }
    }
}
