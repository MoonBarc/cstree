import cstree

on = False
loff = [0] * 150
lon = [0xFFFFFF] * 150

while True:
    cstree.wait()
    if on:
        cstree.render(lon)
    else:
        cstree.render(loff)

    on = not on

