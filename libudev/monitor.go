package libudev

// #cgo LDFLAGS: -ludev
// #include <libudev.h>
// #include <stdlib.h>
import "C"
import "unsafe"

import (
    "syscall"
)

type MonitorDevice struct {
    Action      string
    Device
}

type Monitor struct {
    udev            *C.struct_udev
    udev_monitor    *C.struct_udev_monitor
}

func MonitorUdev(subsystem string) (*Monitor, error) {
    cName := C.CString("udev")
    defer C.free(unsafe.Pointer(cName))

    udev := C.udev_new()
    udev_monitor := C.udev_monitor_new_from_netlink(udev, cName)

    if udev_monitor == nil {
        C.udev_unref(udev)
        return nil, Error{"udev_monitor_new_from_netlink"}
    }

    if subsystem != "" {
        cSubsystem := C.CString(subsystem)
        defer C.free(unsafe.Pointer(cSubsystem))

        if C.udev_monitor_filter_add_match_subsystem_devtype(udev_monitor, cSubsystem, nil) < 0 {
            C.udev_monitor_unref(udev_monitor)
            C.udev_unref(udev)
            return nil, Error{"udev_monitor_filter_add_match_subsystem_devtype"}
        }
    }

    if C.udev_monitor_enable_receiving(udev_monitor) < 0 {
        C.udev_monitor_unref(udev_monitor)
        C.udev_unref(udev)
        return nil, Error{"udev_monitor_enable_receiving"}
    }

    // set blocking
    if fd := C.udev_monitor_get_fd(udev_monitor); fd < 0 {
        C.udev_monitor_unref(udev_monitor)
        C.udev_unref(udev)
        return nil, Error{"udev_monitor_get_fd"}
    } else if err := syscall.SetNonblock(int(fd), false); err != nil {
        C.udev_monitor_unref(udev_monitor)
        C.udev_unref(udev)
        return nil, err
    }

    return &Monitor{udev, udev_monitor}, nil
}

func (self *Monitor) Recv() (MonitorDevice, error) {
    for {
        udev_device := C.udev_monitor_receive_device(self.udev_monitor)

        if udev_device == nil {
            continue
        }

        // device
        var device MonitorDevice

        if err := device.fromUdev(udev_device); err != nil {
            return device, err
        }

        cAction := C.udev_device_get_action(udev_device)
        if cAction != nil {
            device.Action = C.GoString(cAction)
        }

        return device, nil
    }
}

func (self *Monitor) Close() {
    C.udev_monitor_unref(self.udev_monitor)
    C.udev_unref(self.udev)

    self.udev_monitor = nil
    self.udev = nil
}
