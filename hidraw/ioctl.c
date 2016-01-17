#include "ioctl.h"
#include <sys/ioctl.h>

int hidraw_ioctl_getrdescsize(int fd, int *value) {
    return ioctl(fd, HIDIOCGRDESCSIZE, value);
}

int hidraw_ioctl_getrawinfo(int fd, struct hidraw_devinfo *value) {
    return ioctl(fd, HIDIOCGRAWINFO, value);
}
