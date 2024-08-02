#!/usr/bin/env python

import logging
import time
import sys

from flask import Flask, request
from PIL import Image, ImageDraw, ImageFont
from unicornhatmini import UnicornHATMini

app = Flask(__name__)
unicornhatmini = UnicornHATMini()
unicornhatmini.set_brightness(0.1)

def turnOnLight():
    for x in range(17):
        for y in range(7):
            unicornhatmini.set_pixel(x, y, 255, 0, 0)

    unicornhatmini.show()

#Disable logging
log = logging.getLogger('werkzeug')
log.disabled = True
app.logger.disabled = True

@app.route('/')
def home():
    return "pi-inMeeting is running"

@app.route('/api/on', methods=['POST'])
def show_text():
    turnOnLight()
    return "light is on"

@app.route('/api/off', methods=['POST'])
def show_text():
    unicornhatmini.clear()
    return "light is off"


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=80)