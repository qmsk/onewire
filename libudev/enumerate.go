package libudev

// #cgo LDFLAGS: -ludev
// #include <libudev.h>
// #include <stdlib.h>
import "C"
import "unsafe"

/*
 * Filter by match fields:
 *  * Subsystem
 */
func Enumerate(match Device) ([]Device, error) {
    udev := C.udev_new()
    defer C.udev_unref(udev)

    udev_enumerate := C.udev_enumerate_new(udev)
    defer C.udev_enumerate_unref(udev_enumerate)

    if match.Subsystem != "" {
        // libudev strdup's the string
        cSubsystem := C.CString(match.Subsystem)
        defer C.free(unsafe.Pointer(cSubsystem))

        if C.udev_enumerate_add_match_subsystem(udev_enumerate, cSubsystem) < 0 {
            return nil, Error{"udev_enumerate_add_match_subsystem"}
        }
    }

    if C.udev_enumerate_scan_devices(udev_enumerate) < 0 {
        return nil, Error{"udev_enumerate_scan_devices"}
    }

    // list
    var devices []Device

    udev_list_entry := C.udev_enumerate_get_list_entry(udev_enumerate)
    for ; udev_list_entry != nil; udev_list_entry = C.udev_list_entry_get_next(udev_list_entry) {
        udev_device := C.udev_device_new_from_syspath(udev, C.udev_list_entry_get_name(udev_list_entry))
        if udev_device == nil {
            return nil, Error{"udev_device_new_from_syspath"}
        }
        defer C.udev_device_unref(udev_device)

        var device Device

        if err := device.fromUdev(udev_device); err != nil {
            return nil, err
        } else {
            devices = append(devices, device)
        }
    }

    return devices, nil
}

