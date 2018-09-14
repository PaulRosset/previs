package api

import (
	"fmt"
	"os"
	"os/exec"
)

// IsDockerInstall check for the presence of docker on the current machine
func IsDockerInstall() {
	cmd := exec.Command("php", "-v")
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Docker is not installed on your computer.\n%+v", err)
		os.Exit(1)
	}
}

func buildImage(imgDocker string) {
	cwd, _ := os.Getwd()
	cmd := exec.Command("docker", "build", "-t", "treevis", "-f", imgDocker, cwd)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Docker is not installed on your computer.\n%+v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", output)
}

func execContainer() error {

}

func cleaningAfterExec(imgDocker string) {

}

// DoDocker is launching the whole process in order to launch test in a sandboxed environment (Build, Exec, Clean)
func DoDocker(imgDocker string) {
	//buildImage(imgDocker)
	// execContainer()
	// cleaningAfterExec(imgDocker)
}
