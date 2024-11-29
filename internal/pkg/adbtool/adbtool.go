package adbtool

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func IsAppActive(adbPath, serialNumber, targetPackage string) (bool, error) {
	cmd := exec.Command(adbPath, "-s", serialNumber, "shell", "dumpsys", "activity", "activities")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error executing adb command: %v", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	var currentPackage string
	for _, line := range lines {
		if strings.Contains(line, "topResumedActivity") {
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.Contains(part, "/") {
					currentPackage = strings.Split(part, "/")[0]
					break
				}
			}
		}
	}

	return currentPackage == targetPackage, nil
}
