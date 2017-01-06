package main

import (
	"fmt"
	"os/exec"
)

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
