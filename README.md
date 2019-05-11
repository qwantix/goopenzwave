# goopenzwave

Go bindings for the [OpenZWave](https://github.com/OpenZWave/open-zwave) library.


OpenZwave 1.6 compatible

## Warning

This package is still fairly new and so the API is changing pretty rapidly, so be careful if you decide to use it. However, the API will try to mimic the C++ OpenZWave library as much as possible, if it doesn't already, so there shouldn't be many breaking changes.

Most of the C++ OpenZWave library is wrapped now, but should you find anything missing please create a new issue or fork it, implement it yourself and submit a pull request.


## Installing OpenZWave

This package requires a system installation of [OpenZWave](https://github.com/OpenZWave/open-zwave). pkg-config is then used during the build of this package to get the open-zwave library and headers.

Note that package managers may install an old version of the library so a manual build/install from source is preferred.

Example install from source:
1. `git clone https://github.com/OpenZWave/open-zwave.git`
2. `cd open-zwave`
3. `make -j$(nproc)`
4. `sudo make install`
5. You may need to call `sudo ldconfig` now on linux systems
5. See the [open-zwave/INSTALL](https://github.com/OpenZWave/open-zwave/blob/master/INSTALL) file for more information


## Get the Package

```
go get github.com/qwantix/goopenzwave
```


## Example: `gominozw`

This package comes with a basic example, `gominozw`, which is a replica of the original C++ OpenZWave MinOZW utility, now written in Go.

It shows how to set up the Manager with various options and listen for Notifications. Once the initial scan of devices is complete, polling for basic values is set up for the devices.

To install and use:

```
go install github.com/qwantix/goopenzwave/gominozw
gominozw --controller /dev/ttyYourUSBDevice
```


## Notes

### open-zwave build fails with `fatal error: libudev.h: No such file or directory` on Debian/Ubuntu

Try installing libudev with apt and build again.

```sh
apt-get install libudev-dev
cd open-zwave && make
```


### Crashes instantly on macOS 10.12

Do you see something like this when trying to run something with the goopenzwave package?

```
$ ./gominozw -h
zsh: killed     ./gominozw -h
```

You should try building with the `-ldflags=-s` option. E.g.: `go build -ldflags=-s`. More info at [golang/go#19734](https://github.com/golang/go/issues/19734).
