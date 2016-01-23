package main

import (
    "github.com/qmsk/onewire/avrtemp"
    "flag"
    "encoding/hex"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "log"
    "os"
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

    flag.StringVar(&httpListen, "http-listen", "",
        "HTTP Listen: HOST:PORT")
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
    if httpListen != "" {
        log.Printf("http.ListenAndServe %v...\n", httpListen)

        if err := http.ListenAndServe(httpListen, nil); err != nil {
            log.Fatalf("http.ListenAndServe %v: %v\n", httpListen, err)
        }
    }

    /*
    log.Printf("avrtempDevice %v: Read...\n", avrtempDevice)
    for {
        if report, err := avrtempDevice.Read(); err != nil {
            log.Fatalf("avrtemp.Device %v: Read: %v\n", avrtempDevice, err)
        } else {
            log.Printf("%v\n", report)
        }
    }*/
}
