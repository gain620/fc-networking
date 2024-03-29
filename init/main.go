// Copyright Alex Ellis 2023

package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"log"
	"os"
	"os/exec"
	"time"

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

func printLogo() {
	//dev := figure.NewColorFigure("Gain", "", "green", true)
	//dev.Print()
	app := figure.NewColorFigure("FC-VM", "", "blue", true)
	app.Print()
	ver := figure.NewColorFigure("1.0.0", "", "red", true)
	ver.Print()
	fmt.Println()
}

func main() {
	vmInitStart := time.Now()
	printLogo()
	fmt.Printf("Firecracker PoC MicroVM init booting...\nGain Chang 2024\n")

	mount("none", "/proc", "proc", 0)
	mount("none", "/dev/pts", "devpts", 0)
	mount("none", "/dev/mqueue", "mqueue", 0)
	mount("none", "/dev/shm", "tmpfs", 0)
	mount("none", "/sys", "sysfs", 0)
	mount("none", "/sys/fs/cgroup", "cgroup", 0)

	setHostname("fc-microvm")

	fmt.Printf("MicroVM started and running dotnet app ... \n")

	dotnetStart := time.Now()
	// Replace "/path/to/your/app.dll" with the path to your .NET application
	cmd := exec.Command("/usr/bin/dotnet", "/init/dotnet-hello/ConsoleApp2.dll")

	//cmd.Env = append(os.Environ(), paths)
	cmd.Env = append(cmd.Env, paths)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("could not start .NET application, error: %s", err))
	}

	fmt.Printf("Started .NET application\n")
	fmt.Printf("--> PID : %d\n", cmd.Process.Pid)

	//netSetUp := setNetwork()

	err = cmd.Wait()
	if err != nil {
		panic(fmt.Sprintf("could not wait for .NET application, error: %s", err))
	}

	trackDuration(dotnetStart, "Dotnet application")

	// Start a new shell
	shell := exec.Command("/bin/sh")
	shell.Env = append(shell.Env, paths)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	trackDuration(vmInitStart, "VM init")
	err = shell.Run()
	if err != nil {
		panic(fmt.Sprintf("could not start shell, error: %s", err))
	}
}

func trackDuration(start time.Time, desc string) {
	elapsed := time.Since(start)
	log.Printf("%s took %v", desc, elapsed.Seconds())
}
