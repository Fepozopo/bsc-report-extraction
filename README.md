# Commission Report

## Description

This is a Go program that takes an Excel document as input. The document comes from a report that is generated from Sage 100. It takes the data from the report and creates a new file for each sales rep group on the report.

## Motivation

A co-worker came to me and asked if I could automate this task for her. Normally, she would just manually copy and paste the data from the report into their own Excel file to send to each sales rep. It was a lot of work. I wanted to make it easier for her to do that.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- Go programming language (version 1.16 or later)
- A C compiler for your target platform (e.g., GCC for Linux, Clang for macOS, or MinGW for Windows) because this program uses the Fyne GUI toolkit, which requires C bindings.
- You can optionally use zig as a C cross-compiler for all OS targets. The Makefile uses zig for Windows and Linux by default, but you can also use it for macOS if you aren't on that platform.

## Quick Start

1. Clone the repository and navigate to the project root in a terminal.
2. Run `make <target>` to build and run the program. Replace `<target>` with one of the following targets: `windows-amd`, `windows-arm`, `linux-amd`, `linux-arm`, `macos-amd`, or `macos-arm`.

## Usage

1. The program will open a GUI window where you can select the input Excel file.
2. The program will process the Excel file and create a new file for each sales rep group.
3. These new files will be saved in the "output" directory in the same directory as the input file. If the "output" directory does not exist, it will be created.
4. The program will also create a log file in the "logs-bsc" directory in the temporary directory of the current user. If the "logs-bsc" directory does not exist, it will be created.

## Building the Program

To build the program, I've included a Makefile. You can run `make <target>` to build the program for different platforms. You can also run `make clean` to remove the `bin` folder that contains the compiled binaries.
The targets are:
```bash
make windows-amd
make windows-arm
make linux-amd
make linux-arm
make macos-amd
make macos-arm
```
You can also run `make all` to build all the targets.

These commands will build the program for the specified platform and output the binary to the `bin` folder.

## 🤝 Contributing

### Clone the repository

```bash
git clone https://github.com/Fepozopo/bsc-commissions.git
cd bsc-commissions
```

### Submit a Pull Request

Sorry, I'm not accepting any pull requests at the moment.