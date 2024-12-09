import socket
import struct
import time

MAGIC_NUMBER = 0x1999
# MAGIC, 150 32-bit light states
PACKET_FORMAT = "< H 150I"

connection = socket.socket(socket.AddressFamily.AF_UNIX, socket.SocketKind.SOCK_STREAM)

def _private_connect():
    connection.connect("/run/cstree.sock")

def render(lights: list[int]):
    """
    Render takes a list of 150 lights onto the LED strip
    """
    if len(lights) != 150:
        raise ValueError(f"the lights list must be exactly 150 in length! got {len(lights)}")
    header_data = struct.pack(PACKET_FORMAT, MAGIC_NUMBER, *lights)

    connection.sendall(header_data)

def wait():
    """
    Waits a 15th of a second to maintain 15 FPS. If you don't use this, you're
    probably going to see some weird side-effects.
    """
    time.sleep(1/15)

# convenience color definitions
RED = 0xff0000
GREEN = 0x00ff00
BLUE = 0x0000ff