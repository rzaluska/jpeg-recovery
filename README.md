# JPEG File Recovery tool

This program is able to restore JPEG files after accidental remove or partition format.

## Install

```sh
go get -u https://github.com/rzaluska/jpeg-recovery
```

## Usage

### Recover from image file
```sh
jpeg-recovery -f discImage.img
```

JPEG files will be saved to the current directory.


### Recover directly from device

```
jpeg-recovery -f /dev/sda
```

### Changing block size
```sh
jpeg-recovery -f discImage.img -b 4096
```

If you find recovery process slow you can increase block size to look
for JPEG header only every x bytes. By decreasing block size you are allowing
for much more accurate search but time will increase. It is not recommended to
go below 512 bytes (default value for -b) of block size
because filesystems tend not to go below that limit of allocation size.
