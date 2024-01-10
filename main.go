package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/process"
)

var asciiText = `
                ____  __.__.__  .__                
   ____   ____ |    |/ _|__|  | |  |   ___________ 
  / ___\ /  _ \|      < |  |  | |  | _/ __ \_  __ \
 / /_/  >  <_> )    |  \|  |  |_|  |_\  ___/|  | \/
 \___  / \____/|____|__ \__|____/____/\___  >__|   
/_____/                \/                 \/       
                                                   
`

type ProcessInfo struct {
	Name     string
	PID      int32
	RamUsage float64
}

func main() {

	fmt.Println(asciiText)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		color.Cyan("\nChoose an option:\n")
		color.Green("1. List processes\n")
		color.Red("2. Kill a process\n")
		color.Yellow("3. Exit\n")
		fmt.Print("Enter your choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			printProcesses()
		case "2":
			color.Cyan("Enter process name to kill: ")
			scanner.Scan()
			killProcessByName(scanner.Text())
		case "3":
			color.Magenta("Exiting.\n")
			return
		default:
			color.Red("Invalid choice, please try again.\n")
		}
	}
}

func killProcessByName(name string) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error fetching processes:", err)
		return
	}

	found := false

	for _, p := range processes {
		pName, err := p.Name()
		if err != nil {
			continue
		}

		if strings.EqualFold(pName, name) {
			found = true
			fmt.Printf("Killing process: %s (PID: %d)\n", pName, p.Pid)
			if err := p.Kill(); err != nil {
				color.Red("Failed to kill process: %s\n", err)
			} else {
				color.Green("Process killed successfully\n")
			}

		}
	}

	if !found {
		fmt.Println("No process found with the specified name.")
	}
}

func printProcesses() {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error fetching processes:", err)
		return
	}

	var processList []ProcessInfo

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		pid := p.Pid

		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		memUsageMB := float64(memInfo.RSS) / 1024.0 / 1024.0

		processList = append(processList, ProcessInfo{
			Name:     name,
			PID:      pid,
			RamUsage: memUsageMB,
		})
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].RamUsage > processList[j].RamUsage
	})

	fmt.Printf("%-40s %-10s %-10s\n", "Process Name", "PID", "RAM (MB)")
	fmt.Println(strings.Repeat("-", 60))

	for _, processName := range processList {
		fmt.Printf("%-40s %-10d %-10.2f\n", processName.Name, processName.PID, processName.RamUsage)
	}
}
