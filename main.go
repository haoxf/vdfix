package main

import (
	"path"
	"os"
	"os/exec"
	"fmt"
	"bufio"
)

var drives = []string{
	`C:`,
	`D:`,
	`E:`,
	`F:`,
	`G:`,
}

var vmwarePaths = []string{
	`Program Files (x86)\VMware\VMware Workstation`,
	`Program Files\VMware\VMware Workstation`,
}

const diskManagerName = "vmware-vdiskmanager.exe"

func findDiskManager() string {
	for _, vp := range vmwarePaths {
		for _, d := range drives {
			p := path.Join(d, vp, diskManagerName)
			stat, err := os.Stat(p)
			if err == nil && !stat.IsDir() {
				return p
			}
		}
	}

	return ""
}

func fix() bool {
	var vmdk string
	if len(os.Args) < 2 {
		fmt.Print(`Usage:
  vdfix  *.vmdk

Or input the vmdk file: `)

		reader := bufio.NewReader(os.Stdin)
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			return false
		}

		vmdk = string(strBytes)
	} else {
		vmdk = os.Args[1]
	}

	dm := findDiskManager()
	if dm == "" {
		fmt.Printf("error: cannot found '%s'\n", diskManagerName)
		return false
	}

	cmd := exec.Command(dm, "-R", vmdk)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	return err == nil && len(out) == 0
}

func main() {
	ret := fix()
	if !ret {
		fmt.Print("Press 'Enter' to continue..")
		fmt.Scanln()
	}
}
