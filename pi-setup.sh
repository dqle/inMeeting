#!/bin/bash

# Enable SPI on your Pi
raspi-config nonint do_spi 0

# Install dependencies
apt-get update
apt-get install python3 python3-pip git python3-pil

# Allow pip packages to be installed in global context
sudo rm /usr/lib/python3.11/EXTERNALLY-MANAGED

# git clone inMeeting
git clone https://github.com/dqle/inMeeting.git /home/pi/

# Install Piromoni Unicorn Hat Mini python package
cd /home/pi/inMeeting/pi
echo no | ./unicornhat-mini/install.sh

# Install Python package
pip3 install -r requirements.txt

# Add program as a service
cp pi-inmeeting.service /etc/systemd/system/
systemctl daemon-reload

# Enable and Start service
systemctl enable pi-attention
systemctl start pi-attention