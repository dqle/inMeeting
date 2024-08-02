#!/bin/bash

echo "----------------------"
echo "Enable SPI on your Pi"
echo "----------------------"

echo "raspi-config nonint do_spi 0"
raspi-config nonint do_spi 0

echo ""
echo "----------------------"
echo "Install dependencies"
echo "----------------------"
apt-get update
apt-get install python3 python3-pip git python3-pil


echo ""
echo "----------------------"
echo "Allow pip packages to be installed in global context"
echo "----------------------"

echo "rm /usr/lib/python3.11/EXTERNALLY-MANAGED"
rm /usr/lib/python3.11/EXTERNALLY-MANAGED


echo ""
echo "----------------------"
echo "git clone inMeeting"
echo "----------------------"

echo "git clone https://github.com/dqle/inMeeting.git /home/pi/"
git clone https://github.com/dqle/inMeeting.git /home/pi/


echo ""
echo "----------------------"
echo "Install Piromoni Unicorn Hat Mini python package"
echo "----------------------"

cd /home/pi/inMeeting/pi
echo no | ./unicornhat-mini/install.sh

echo ""
echo "----------------------"
echo "Install Python package"
echo "----------------------"

pip3 install -r requirements.txt


echo ""
echo "----------------------"
echo "Add program as a service"
echo "----------------------"

cp pi-inmeeting.service /etc/systemd/system/
systemctl daemon-reload


echo ""
echo "----------------------"
echo "Enable and Start service"
echo "----------------------"

systemctl enable pi-inmeeting
systemctl start pi-inmeeting