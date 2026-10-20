# A printer of the raw port 9100 for test/print/ci.sh: each connection
# that sends something is a job, kept in a directory as <n>.prn.
#
#   python3 raw9100.py <port> <dir>
import os
import socket
import sys

srv = socket.socket()
srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
srv.bind(("127.0.0.1", int(sys.argv[1])))
srv.listen()
n = 0
while True:
    c, _ = srv.accept()
    data = b""
    while True:
        b = c.recv(65536)
        if not b:
            break
        data += b
    c.close()
    if data:
        n += 1
        with open(os.path.join(sys.argv[2], "%d.prn" % n), "wb") as f:
            f.write(data)
