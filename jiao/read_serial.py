import serial
import websocket

ser = serial.Serial()
ser.baudrate = 9600
ser.port = "COM3"
ser.open()

buffer = []
line = ser.readline()
count = 0
C = 10

class WebSocketSender:
    def __init__(self, url):
        self.ws = websocket.create_connection(url)

    def send(self, buffer):
        print("sending: ", str(buffer))
        self.ws.send(str(buffer))

sender = WebSocketSender("ws://localhost:8080/sonar")
while line:
    #print(int(line))
    line = ser.readline()
    line = line.strip(b'\n')
    buffer.append(int(line))
    count += 1
    if count == C:
        try:
            sender.send(buffer)
            buffer.clear()
            count = 0
        except Exception as e:
            sender.ws.close()
            print(e)
        



