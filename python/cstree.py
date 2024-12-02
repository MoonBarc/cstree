import socket

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect("localhost:1999")

def render(lights):
    if not isinstance(lights, list):
        raise "the lights parameter must be a list of colors"
    
    