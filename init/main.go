// Copyright Alex Ellis 2023

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"syscall"
)

const paths = "PATH=/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin"

// main starts an init process that can prepare an environment and start a shell
// after the Kernel has started.
func main2() {
	fmt.Printf("Firecracker PoC MicroVM init booting...\nCopyright Alex Ellis 2024\n")

	mount("none", "/proc", "proc", 0)
	mount("none", "/dev/pts", "devpts", 0)
	mount("none", "/dev/mqueue", "mqueue", 0)
	mount("none", "/dev/shm", "tmpfs", 0)
	mount("none", "/sys", "sysfs", 0)
	mount("none", "/sys/fs/cgroup", "cgroup", 0)

	setHostname("fc-microvm")

	fmt.Printf("Sandbox MicroVM starting... /bin/sh\n")

	cmd := exec.Command("/bin/sh")

	cmd.Env = append(cmd.Env, paths)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("could not start /bin/sh, error: %s", err))
	}

	//netSetUp := setNetwork()

	err = cmd.Wait()
	if err != nil {
		panic(fmt.Sprintf("could not wait for /bin/sh, error: %s", err))
	}
}

func setHostname(hostname string) {
	err := syscall.Sethostname([]byte(hostname))
	if err != nil {
		panic(fmt.Sprintf("cannot set hostname to %s, error: %s", hostname, err))
	}
}

func setNetwork() {
	exec.Command("/bin/sh", "-c", "ip link set dev eth0 up")
	// ip link set dev eth0 up
	// ip addr add
	// ip route add default via
}

func mount(source, target, filesystemtype string, flags uintptr) {

	if _, err := os.Stat(target); os.IsNotExist(err) {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			panic(fmt.Sprintf("error creating target folder: %s %s", target, err))
		}
	}

	err := syscall.Mount(source, target, filesystemtype, flags, "")
	if err != nil {
		log.Printf("%s", fmt.Errorf("error mounting %s to %s, error: %s", source, target, err))
	}
}

func main() {
	fmt.Printf("Firecracker PoC MicroVM init booting...\nCopyright Alex Ellis 2023\n")

	mount("none", "/proc", "proc", 0)
	mount("none", "/dev/pts", "devpts", 0)
	mount("none", "/dev/mqueue", "mqueue", 0)
	mount("none", "/dev/shm", "tmpfs", 0)
	mount("none", "/sys", "sysfs", 0)
	mount("none", "/sys/fs/cgroup", "cgroup", 0)

	setHostname("fc-microvm")

	fmt.Printf("MicroVM starting and running dotnet app ... \n")

	// Replace "/path/to/your/app.dll" with the path to your .NET application
	cmd := exec.Command("dotnet", "/init/dotnet-hello/ConsoleApp2.dll")

	cmd.Env = append(cmd.Env, paths)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("could not start .NET application, error: %s", err))
	}

	fmt.Printf("Started .NET application with PID %d\n", cmd.Process.Pid)

	//netSetUp := setNetwork()

	err = cmd.Wait()
	if err != nil {
		panic(fmt.Sprintf("could not wait for .NET application, error: %s", err))
	}
}
