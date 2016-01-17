package udev

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
)

type DeviceNode string // path

func (self DeviceNode) Name() string {
    return path.Base(string(self))
}

func (self DeviceNode) Path() string {
    return string(self)
}

func (self DeviceNode) UeventInfo() (UeventInfo, error) {
    return readUevent(path.Join(string(self), "uevent"))
}

func (self DeviceNode) DevFile() (string, error) {
    if ueventInfo, err := self.UeventInfo(); err != nil {
        return "", err
    } else if devFile := ueventInfo.DevFile(); devFile == "" {
        return "", fmt.Errorf("No uevent DEVNAME=")
    } else {
        return devFile, nil
    }
}

func ListClass(class string) (nodes []DeviceNode, err error) {
    dirPath := path.Join("/sys/class", class)

    readdir, err := ioutil.ReadDir(dirPath)
    if err != nil {
        return nil, err
    }

    for _, fileInfo := range readdir {
        /* if fileInfo.Mode() & os.ModeDir == 0 || fileInfo.Mode() & os.ModeSymlink == 0 {
            continue
        } */

        filePath := path.Join(dirPath, fileInfo.Name())

        linkPath, err := os.Readlink(filePath)
        if err != nil {
            continue
        }

        devicePath := path.Clean(path.Join(dirPath, linkPath))

        nodes = append(nodes, DeviceNode(devicePath))
    }

    return nodes, nil
}
