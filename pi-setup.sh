#!/bin/bash

# Capture the current directory
currentDir=$(pwd)

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

echo "git clone https://github.com/dqle/inMeeting.git"
git clone https://github.com/dqle/inMeeting.git

echo ""
echo "----------------------"
echo "Install Piromoni Unicorn Hat Mini python package"
echo "----------------------"

cd $currentDir/inMeeting/pi/unicornhat-mini/
echo no | sudo bash ./install.sh

echo ""
echo "----------------------"
echo "Install Python package"
echo "----------------------"

cd $currentDir/inMeeting/pi
pip3 install -r requirements.txt

echo ""
echo "----------------------"
echo "Add program as a service"
echo "----------------------"

cd $currentDir
sed -i "s|PWD|$currentDir|g" $currentDir/inMeeting/pi/pi-inmeeting.service
echo "cp $currentDir/inMeeting/pi/pi-inmeeting.service /etc/systemd/system/"
cp $currentDir/inMeeting/pi/pi-inmeeting.service /etc/systemd/system/
systemctl daemon-reload


echo ""
echo "----------------------"
echo "Enable and Start service"
echo "----------------------"

echo "systemctl enable pi-inmeeting && systemctl start pi-inmeeting"
systemctl enable pi-inmeeting && systemctl start pi-inmeeting

echo ""
echo "----------------------"
echo "Reboot"
echo "----------------------"

echo "reboot"
reboot