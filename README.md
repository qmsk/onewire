## Dependencies

* libudev-dev

## Configuration

### `/etc/udev/rules.d/90-hidraw.rules`

    KERNEL=="hidraw*", \
        GROUP="plugdev", MODE=0660

