package main

import (
	"fmt"
	"os/exec"
)

func runApp(app string) error {
	err := exec.Command("/bin/bash", "-c", app).Run()
	if err != nil {
		logger.Errorf("Exec '%s' failed: %v", app, err)
	}
	return err
}

func doAction(cmd string) error {
	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		if len(out) == 0 {
			return err
		}
		return fmt.Errorf("%s", string(out))
	}
	return nil
}
