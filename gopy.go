package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func runGopy(pkg string) {
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		nativeGopy(pkg)
		return
	}
	dockerGopy(pkg)
}

func nativeGopy(pkg string) {
	if *debug {
		log.Println("gopy bind native", pkg)
	}
	path := getExecPath("gopy")
	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{path, "bind", pkg},
		Stderr: os.Stderr,
	}
	if *debug {
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func mountPackageDir(pkg, dst string) string {
	return fmt.Sprintf("%s:%s", packageDir(pkg), dst)
}

func dockerGopy(pkg string) {
	if *debug {
		log.Println("gopy bind via Docker", pkg)
	}
	pwd, _ := os.Getwd()
	dockerPath := getExecPath("docker")
	cmd := exec.Cmd{
		Path: dockerPath,
		Args: []string{dockerPath,
			"run", "-a", "stdout", "-a", "stderr", "--rm",
			"-v", mountPackageDir(pkg, "/go/src/in"), "-v", fmt.Sprintf("%s:/out", pwd),
			"xlab/gopy", "app", "bind", "-output", "/out", "in",
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	// if *debug {
	// 	cmd.Stdout = os.Stdout
	// }
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func getDockerVersion(path string) (string, error) {
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}
