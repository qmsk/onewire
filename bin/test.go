package main

import (
    "github.com/qmsk/onewire/avrtemp"
    "flag"
    "encoding/hex"
    "log"
    "os"
    "github.com/qmsk/onewire/hidraw"
)

var (
    deviceConfig hidraw.DeviceConfig
)

func init() {
    flag.UintVar(&deviceConfig.Vendor, "device-vendor", avrtemp.HIDRAW_CONFIG.Vendor,
        "Select device vendor")
    flag.UintVar(&deviceConfig.Product, "device-product", avrtemp.HIDRAW_CONFIG.Product,
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

    if hidrawList, err := hidraw.List(); err != nil {
        log.Fatalf("hidraw.List: %v\n", err)
    } else {
        log.Printf("hidraw.List:\n")
        for _, deviceInfo := range hidrawList {
            log.Printf("\t%#v\n", deviceInfo)
        }
    }

    hidrawDevice, err := hidraw.Find(deviceConfig)
    if err != nil {
        log.Fatalf("hidraw.Select %v: %v\n", deviceConfig, err)
    }

    if devInfo, err := hidrawDevice.DevInfo(); err != nil {
        log.Fatalf("hidraw.Device %v: DevInfo: %v\n", hidrawDevice, err)
    } else {
        log.Printf("hidraw.Device %v: DevInfo: %#v\n", hidrawDevice, devInfo)
    }

    if reportDescriptor, err := hidrawDevice.ReportDescriptor(); err != nil {
        log.Fatalf("hidraw.Device %v: ReportDescriptor: %v\n", hidrawDevice, err)
    } else {
        log.Printf("hidraw.Device %v: ReportDescriptor:\n%v\n", hidrawDevice, hex.Dump(reportDescriptor))
    }

    avrtempDevice, err := avrtemp.Open(hidrawDevice)
    if err != nil {
        log.Fatalf("avrtemp.Open %v: %v\n", hidrawDevice, err)
    } else {
        log.Printf("avrtemp.Open %v: %v\n", hidrawDevice, avrtempDevice)
    }

    log.Printf("avrtempDevice %v: Read...\n", avrtempDevice)
    for {
        if report, err := avrtempDevice.Read(); err != nil {
            log.Fatalf("avrtemp.Device %v: Read: %v\n", avrtempDevice, err)
        } else {
            log.Printf("%v\n", report)
        }
    }
}
