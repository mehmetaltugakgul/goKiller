![image](https://github.com/mehmetaltugakgul/goKiller/assets/10194009/b2db04a6-aeff-444e-b568-fd757cfd88a8)

# goKiller App

## Introduction

The goKiller app is a simple command-line utility written in Go that allows you to list running processes and optionally terminate them by name. It provides an easy-to-use interface for viewing and managing processes with a focus on RAM (memory) usage. The app uses the `github.com/shirou/gopsutil/process` package to interact with and retrieve information about running processes.

## Features

- **List Processes:** View a list of all running processes along with their names, Process IDs (PIDs), and RAM usage in megabytes (MB).

- **Kill a Process:** Terminate a running process by entering its name.

- **Exit the App:** Quit the goKiller App.

## Installation

Before you can use the goKiller App, make sure you have Go installed on your system. You can download and install Go from the official website: [https://golang.org/dl/](https://golang.org/dl/)

Once Go is installed, follow these steps to run the app:

1. Open a terminal or command prompt.

2. Clone the GitHub repository containing the app's source code:

   ```bash
   git clone <repository_url>

3. 
 
