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