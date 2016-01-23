package main

import (
    "github.com/qmsk/onewire/avrtemp"
    "flag"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "log"
    "github.com/qmsk/onewire/server"
)

var (
    deviceConfig hidraw.DeviceConfig
    httpListen  string
)

func init() {
    flag.UintVar(&deviceConfig.VendorID, "device-vendor", avrtemp.HIDRAW_CONFIG.VendorID,
        "Select device vendor")
    flag.UintVar(&deviceConfig.ProductID, "device-product", avrtemp.HIDRAW_CONFIG.ProductID,
        "Select device product")

    flag.StringVar(&httpListen, "http-listen", ":8283",
        "HTTP Listen: HOST:PORT")
}

func main() {
    flag.Parse()

    s, err := server.New()
    if err != nil {
        log.Fatalf("server.New: %v\n", err)
    }

    http.Handle("/api/", http.StripPrefix("/api", s))


    if hidrawList, err := hidraw.List(deviceConfig); err != nil {
        log.Fatalf("hidraw.List %v: %v\n", deviceConfig, err)
    } else {
        log.Printf("hidraw.List...\n")
        for _, deviceInfo := range hidrawList {
            go s.AddHidrawDevice(deviceInfo)
        }
    }

    if monitorChan, err := hidraw.Monitor(deviceConfig); err != nil {
        log.Fatalf("hidraw.Monitor %v: %v\n", deviceConfig, err)
    } else {
        log.Printf("hidraw.Monitor...\n")

        go s.MonitorHidraw(monitorChan)
    }

    // run
    log.Printf("http.ListenAndServe %v...\n", httpListen)

    if err := http.ListenAndServe(httpListen, nil); err != nil {
        log.Fatalf("http.ListenAndServe %v: %v\n", httpListen, err)
    }
}
