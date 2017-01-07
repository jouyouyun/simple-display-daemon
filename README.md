# Simple Display Daemon

Display Daemon provides some dbus interfaces to list connected outputs

## Dependencies


### Build Dependencies

* golang
* [go-lib](https://github.com/linuxdeepin/go-lib)
* [dde-api](https://github.com/linuxdeepin/dde-api)
* [xgb](https://github.com/BurntSushi/xgb)
* [xgbutil](https://github.com/BurntSushi/xgbutil)
* [goconvey](https://github.com/smartystreets/goconvey/convey)

### Runtime Dependencies

* libx11
* libxtst
* gtk+3
* xrandr
* deepin-terminal

## Installation

Install prerequisites

```shell
$ go get github.com/BurntSushi/xgb
$ go get github.com/BurntSushi/xgbutil
$ go get github.com/smartystreets/goconvey/convey
```

Build:
```
$ make GOPATH=/usr/share/gocode
```

Install:
```
sudo make install
```

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org/)
* [WiKi](http://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes.

* [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers)
