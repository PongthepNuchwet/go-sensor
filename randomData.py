import random
from datetime import datetime
import time
import requests


headers = {'Accept': 'application/json'}
pump = ['ON', 'OFF']

for i in range(0, 24):

    dt = datetime(2023, 4, 15, i, 0, 0, 0)
    print('Input Datetime:', dt)

    # convert datetime to ISO date
    iso_date = dt.isoformat()
    print('ISO Date:', iso_date)
    json = {
        "Value": pump[random.randint(0, 51) % 2],
        "CreatedAt": iso_date+"+07:00"
    }
    res = requests.post('http://127.0.0.1:8080/api/pumps',
                        headers=headers, json=json)
    print(res.json())
