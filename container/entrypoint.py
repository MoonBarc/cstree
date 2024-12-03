import cstree

print("-> Connecting to tree...")
cstree._private_connect()
print("-> Running program...")

exec(open("/sandbox/script.py").read())