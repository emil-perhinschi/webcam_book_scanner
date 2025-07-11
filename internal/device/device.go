package device

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"gocv.io/x/gocv"
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
		// fmt.Println("MATCH: ", match)
		if match {
			// fmt.Println("OUTPUT: ", out.String())
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

func MatToPixbuf(mat *gocv.Mat) (*gdk.Pixbuf, error) {
	// Convert Mat to bytes (assuming RGB image)
	data, err := gocv.IMEncode(".png", *mat)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Mat to PNG: %v", err)
	}
	defer data.Close()

	// Create a PixbufLoader to load the image data
	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, fmt.Errorf("failed to create PixbufLoader: %v", err)
	}
	defer loader.Close()

	// Write the image data to the loader
	_, err = loader.Write(data.GetBytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write to PixbufLoader: %v", err)
	}

	// Get the Pixbuf from the loader
	pixbuf, err := loader.GetPixbuf()

	if err != nil {
		return nil, fmt.Errorf("failed to get Pixbuf: %v", err)
	}

	return pixbuf, nil
}

/*
Driver Info:
	Driver name      : uvcvideo
	Card type        : Logitech BRIO
	Bus info         : usb-0000:13:00.0-4.3
	Driver version   : 6.12.3
	Capabilities     : 0x84a00001
		Video Capture
		Metadata Capture
		Streaming
		Extended Pix Format
		Device Capabilities
	Device Caps      : 0x04200001
		Video Capture
		Streaming
		Extended Pix Format
Media Driver Info:
	Driver name      : uvcvideo
	Model            : Logitech BRIO
	Serial           : 50316219
	Bus info         : usb-0000:13:00.0-4.3
	Media version    : 6.12.3
	Hardware revision: 0x00000317 (791)
	Driver version   : 6.12.3
Interface Info:
	ID               : 0x03000002
	Type             : V4L Video
Entity Info:
	ID               : 0x00000001 (1)
	Name             : Logitech BRIO
	Function         : V4L2 I/O
	Flags            : default
	Pad 0x01000007   : 0: Sink
	  Link 0x0200001f: from remote pad 0x100000a of entity 'Processing 3' (Video Pixel Formatter): Data, Enabled, Immutable
*/

func ParseV4L2DeviceDetails(deviceDescription string) map[string]string {
	var result map[string]string

	return result
}
