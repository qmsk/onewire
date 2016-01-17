package udev

import (
    "bufio"
    "fmt"
    "os"
    "path"
    "strings"
)

type UeventInfo map[string]string

func (self UeventInfo) DevName() string {
    return self["DEVNAME"]
}

func (self UeventInfo) DevFile() string {
    devName := self.DevName()

    if devName == "" {
        return ""
    } else if devName[0] == '/' {
        return devName
    } else {
        return path.Join("/dev", devName)
    }
}

func readUevent(ueventPath string) (UeventInfo, error) {
    file, err := os.Open(ueventPath)
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(file)
    info := make(UeventInfo)

    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "=", 2)

        if len(parts) != 2 {
            return nil, fmt.Errorf("Invalid line: %#v", parts)
        }

        info[parts[0]] = parts[1]
    }

    return info, nil
}

