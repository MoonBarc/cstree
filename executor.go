package main

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var ActiveProgram *Program
var SkipProgram chan bool = make(chan bool, 5)

// can't believe go doesn't have a copy function
func copy(in, out string) {
	r, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	w, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	w.ReadFrom(r)
}

func ExecuteProgram(program *Program) {
	ctx := context.Background()
	prog := program.program

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalln("Couldn't connect to Docker", err)
	}

	cfg := container.Config{
		Image: "python:3-slim",
		Cmd:   []string{"python3", "/sandbox/entrypoint.py"},
	}

	containerDir, err := os.MkdirTemp("", "cstree-container")
	log.Println("tmp=", containerDir)
	if err != nil {
		log.Fatalln("failed to make temp dir", err)
	}
	defer os.RemoveAll(containerDir)

	err = os.WriteFile(path.Join(containerDir, "script.py"), []byte(prog), os.ModePerm)
	copy("./container/cstree.py", path.Join(containerDir, "cstree.py"))
	copy("./container/entrypoint.py", path.Join(containerDir, "entrypoint.py"))

	if err != nil {
		log.Fatalln("failed to write progfile")
	}

	hostCfg := container.HostConfig{
		Binds: []string{
			containerDir + ":/sandbox",
			sockAddr + ":/run/cstree.sock",
		},
		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"max-size": "10m",
				"max-file": "1",
			},
		},
	}

	c, err := cli.ContainerCreate(ctx, &cfg, &hostCfg, nil, nil, "")
	if err != nil {
		log.Fatalln("failed to make container", err)
	}

	err = cli.ContainerStart(ctx, c.ID, container.StartOptions{})
	if err != nil {
		log.Fatalln("failed to start container", err)
	}

	started := time.Now()
	done, err_chan := cli.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)
	timer := time.NewTimer(2*time.Minute + 15*time.Second)

	kill := func() {
		// time has expired
		timeout := 1
		err = cli.ContainerStop(ctx, c.ID, container.StopOptions{
			Timeout: &timeout,
		})
		if err != nil {
			log.Println("warn: couldn't stop container", err)
		}
	}

	select {
	case <-done:
		if time.Since(started).Seconds() < 15 {
			// program likely crashed :(
			// let's disable it for now:
			TrDB.SetDisabled(program.id, true)
		}
	case err = <-err_chan:
		log.Println("error watching container", err)
		kill()
	case <-timer.C:
		kill()
	case <-SkipProgram:
		log.Println("skipping current program...")
		kill()
	}

	// salvage logs
	logs, err := cli.ContainerLogs(ctx, c.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})

	if err != nil {
		log.Println("warn, couldn't get logs for container", c.ID, err)
	} else {
		// persist logs
		defer logs.Close()
		os.Stdout.ReadFrom(logs)
		// 	betterFile, err := os.OpenFile("/tmp/cstree-"+c.ID+"-log.json", os.O_CREATE|os.O_WRONLY, 0o644)

		// 	if err != nil {
		// 		log.Println("persisting logs failed", err)
		// 	}

		// 	defer betterFile.Close()
		// 	betterFile.ReadFrom(logs)
	}

	// remove container
	err = cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{})
	if err != nil {
		log.Fatalln("couldn't remove container", err)
	}
}
