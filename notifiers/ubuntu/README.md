iDoneThis for Ubuntu
====================

This client is intended to be installed into /opt/idonethis.
For manual installation, run following commands:

    sudo mkdir /opt/idonethis

    sudo cp indicate.png /opt/idonethis

    sudo cp idt.glade /opt/idonethis

    go build

    sudo cp ubuntu /opt/idonethis/idonethis

    Symbolic link the idonethis binary into /usr/bin

Also I totally need to build a deb package for this
