#!/usr/bin/env python

import logging
from flask import Flask
from unicornhatmini import UnicornHATMini

app = Flask(__name__)
unicornhatmini = UnicornHATMini()
unicornhatmini.set_brightness(0.1)
light_status = "off"

def turn_on_light():
    unicornhatmini.set_all(255, 0, 0)
    global light_status
    light_status = "on"
    unicornhatmini.show()

def turn_off_light():
    unicornhatmini.clear()
    global light_status
    light_status = "off"
    unicornhatmini.show()

# Disable logging
log = logging.getLogger('werkzeug')
log.disabled = True
app.logger.disabled = True

@app.route('/')
def home():
    return "pi-inMeeting is running"

@app.route('/status')
def get_status():
    return light_status

@app.route('/api/on', methods=['POST'])
def api_on():
    turn_on_light()
    return "light is on"

@app.route('/api/off', methods=['POST'])
def api_off():
    turn_off_light()
    return "light is off"

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=80)
