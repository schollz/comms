# comms


A simple tool to communicate with serial devices or midi devices (via sysex). Useful for embedded hardware development.

## windows instructions

Download MSYS: https://www.msys2.org/

Run MSYS and install the following:

```powershell
pacman -S mingw-w64-x86_64-rtmidi
pacman -S mingw-w64-x86_64-toolchain
```

Before building with Go, update the environmental variables:

```powershell
$env:Path += ";C:\msys64\mingw64\bin"
$env:CGO_ENABLED = "1"
$env:CC="x86_64-w64-mingw32-gcc"
$env:CGO_LDFLAGS="-static"
go build -v -x
```

## install

install with

```
curl https://getcomms.schollz.com 
```


## usage

To communicate with a midi devices:

```
comms --midi zeptocore
```

To communicate with a serial device:

```
comms --serial
```