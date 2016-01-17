#include <linux/hidraw.h>

int hidraw_ioctl_getrdescsize(int fd, int *value);
int hidraw_ioctl_getrawinfo(int fd, struct hidraw_devinfo *value);
