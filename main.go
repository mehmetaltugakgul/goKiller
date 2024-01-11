package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
	Name      string
	PID       int32
	RamUsage  float64
	CpuUsage  float64
	User      string
	StartTime string
}

func main() {
	fmt.Println(asciiText)
	fmt.Println(color.GreenString("The goKiller app is a simple command-line utility written in Go\nthat allows you to list running processes and optionally terminate\nthem by name. It provides an easy-to-use interface for viewing and\nmanaging processes with a focus on RAM (memory) usage."))

	scanner := bufio.NewScanner(os.Stdin)

	for {
		color.Cyan("\nChoose an option:\n")
		color.Green("1. List processes\n")
		color.Red("2. Kill a process by name\n")
		color.Yellow("3. Kill a process by PID\n")
		color.White("4. Search processes by name\n")
		color.Magenta("5. Exit\n")
		fmt.Print("Enter your choice (1 or 2 or 3 or 4 or 5): ")

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
			color.Cyan("Enter PID of the process to kill: ")
			scanner.Scan()
			pidStr := scanner.Text()
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				color.Red("Invalid PID. Please enter a valid integer PID.\n")
				continue
			}
			killProcessByPID(pid)
		case "4":
			color.Cyan("Enter a search term for process names: ")
			scanner.Scan()
			searchTerm := scanner.Text()
			searchProcessesByName(searchTerm)
		case "5":
			color.Magenta("Exiting.\n")
			return
		default:
			color.Red("Invalid choice, please try again.\n")
		}
	}
}

func killProcessByPID(pid int) {
	process, err := process.NewProcess(int32(pid))
	if err != nil {
		color.Red("Error fetching process information:", err)
		return
	}

	pName, err := process.Name()
	if err != nil {
		color.Red("Error fetching process name:", err)
		return
	}

	fmt.Printf("Killing process: %s (PID: %d)\n", pName, pid)
	if err := process.Terminate(); err != nil {
		color.Red("Failed to kill process: %s\n", err)
	} else {
		color.Green("Process killed successfully\n")
	}
}

func searchProcessesByName(searchTerm string) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error fetching processes:", err)
		return
	}

	var matchingProcesses []ProcessInfo

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToLower(name), strings.ToLower(searchTerm)) {
			pid := p.Pid

			memInfo, err := p.MemoryInfo()
			if err != nil {
				continue
			}
			memUsageMB := float64(memInfo.RSS) / 1024.0 / 1024.0

			cpuPercent, err := p.Percent(0)
			if err != nil {
				continue
			}

			createTime, err := p.CreateTime()
			if err != nil {
				continue
			}
			startTime := time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05")

			username, err := p.Username()
			if err != nil {
				continue
			}

			matchingProcesses = append(matchingProcesses, ProcessInfo{
				Name:      name,
				PID:       pid,
				RamUsage:  memUsageMB,
				CpuUsage:  cpuPercent,
				User:      username,
				StartTime: startTime,
			})
		}
	}

	if len(matchingProcesses) == 0 {
		fmt.Printf("No processes found matching the search term '%s'.\n", searchTerm)
	} else {
		sort.Slice(matchingProcesses, func(i, j int) bool {
			return matchingProcesses[i].RamUsage > matchingProcesses[j].RamUsage
		})

		fmt.Printf("%-40s %-10s %-10s %-10s %-15s %-30s\n", "Process Name", "PID", "RAM (MB)", "CPU (%)", "User", "Start Time")
		fmt.Println(strings.Repeat("-", 100))

		for _, processInfo := range matchingProcesses {
			fmt.Printf("%-40s %-10d %-10.2f %-10.2f %-15s %-30s\n", processInfo.Name, processInfo.PID, processInfo.RamUsage, processInfo.CpuUsage, processInfo.User, processInfo.StartTime)
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
			if err := p.Terminate(); err != nil {
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

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		createTime, err := p.CreateTime()
		if err != nil {
			continue
		}
		startTime := time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05")

		username, err := p.Username()
		if err != nil {
			continue
		}

		processList = append(processList, ProcessInfo{
			Name:      name,
			PID:       pid,
			RamUsage:  memUsageMB,
			CpuUsage:  cpuPercent,
			User:      username,
			StartTime: startTime,
		})
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].RamUsage > processList[j].RamUsage
	})

	fmt.Printf("%-40s %-10s %-10s %-10s %-15s %-30s\n", "Process Name", "PID", "RAM (MB)", "CPU (%)", "User", "Start Time")
	fmt.Println(strings.Repeat("-", 100))

	for _, processInfo := range processList {
		fmt.Printf("%-40s %-10d %-10.2f %-10.2f %-15s %-30s\n", processInfo.Name, processInfo.PID, processInfo.RamUsage, processInfo.CpuUsage, processInfo.User, processInfo.StartTime)
	}
}
