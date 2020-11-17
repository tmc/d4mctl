package d4m

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Restart restarts docker-for-mac
func Restart(shouldWait bool) error {
	fmt.Println("restarting docker-for-mac")
	restartCmd := `
osascript -e 'quit app "Docker"';
open -a Docker;
while [ -z "$(docker info 2> /dev/null )" ]; do printf "."; sleep 1; done; echo ""
`
	cmd := exec.Command("sh", "-c", restartCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if !shouldWait {
		return err
	}

	var ready bool
	// Try to connect to the docker engine.
	for i := 0; i < 30; i++ {
		ready, err = dockerIsReady(false)
		if ready {
			return nil
		}
		fmt.Printf(".")
		time.Sleep(time.Second)
	}
	// Run one last time with verbose=true to hint the user.
	_, err = dockerIsReady(true)
	return err
}

func dockerIsReady(verbose bool) (bool, error) {
	cmd := exec.Command("docker", "ps")
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	return cmd.ProcessState.ExitCode() == 0, err
}
