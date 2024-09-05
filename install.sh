#!/bin/bash 

ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
else 
    # exit if not supported
    echo "The architecture $ARCH is not supported."
    exit
fi
if [ "$OS" != "linux" ]; then
    echo "The OS $OS is not supported."
    exit
fi
echo "Downloading comms for $OS $ARCH..."
curl -s https://api.github.com/repos/schollz/comms/releases/latest | \
    grep 'comms_'$OS'_'$ARCH | \
    grep 'browser' | cut -d : -f 2,3 | tr -d \" | wget --show-progress -O comms -qi -
chmod +x comms
if [ ! -f comms ]; then
    echo "Failed to download comms."
    exit
else
    sudo mv comms /usr/local/bin/comms
    echo "Downloaded comms to /usr/local/bin/comms"
    comms --version
fi
