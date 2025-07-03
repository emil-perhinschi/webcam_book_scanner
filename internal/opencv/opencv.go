package opencv

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func ListDevices() ([]string, error) {
	var result []string
	devices, err := listFilesByPattern("/dev/v4l/by-id/usb*")
	if err != nil {
		return nil, err
	}

	for _, device := range devices {

		var out bytes.Buffer
		var outErr bytes.Buffer
		cmd := exec.Command("/bin/sh", "-c", "v4l2-ctl -D --device="+device)
		cmd.Stdout = &out
		cmd.Stderr = &outErr
		err := cmd.Run()

		if err != nil {
			fmt.Println("ERROR: ", outErr.String())
			continue
		}

		r, errRegex := regexp.Compile(`(?s)\tDevice Caps.*\n\t\tVideo Capture`)
		if errRegex != nil {
			fmt.Println("ERROR REGEX: ", errRegex)
			continue
		}

		match := r.MatchString(out.String())
		fmt.Println("MATCH: ", match)
		if match {
			fmt.Println("OUTPUT: ", out.String())
			result = append(result, device)
		}
	}
	return result, nil
}

func listFilesByPattern(pattern string) ([]string, error) {
	var result []string
	var out bytes.Buffer
	var outErr bytes.Buffer

	// exec does not expand *
	cmd := exec.Command("/bin/sh", "-c", "ls "+pattern)
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()

	if err != nil {
		return nil, errors.New(outErr.String())
	}
	fmt.Println(out.String())
	devices := strings.Split(out.String(), "\n")

	for i, d := range devices {
		if d != "" {
			result = append(result, d)
		}
		fmt.Println(">>>> ", i, " ", d)
	}

	return result, nil
}
