#include "ioctl.h"
#include <sys/ioctl.h>

int hidraw_ioctl_getrawinfo(int fd, struct hidraw_devinfo *value) {
    return ioctl(fd, HIDIOCGRAWINFO, value);
}
