import cstree

lights = [0xed75fa] * 150
i = 0
while True:
    def j(i):
        return i % 150

    lights[j(i)] = cstree.BLUE
    lights[j(i + 1)] = cstree.GREEN
    lights[j(i + 2)] = cstree.GREEN
    lights[j(i + 3)] = cstree.RED

    cstree.render(lights)
    cstree.wait()

    i += 1