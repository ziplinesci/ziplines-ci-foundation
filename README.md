# Ziplines CI Foundation
The `ziplines-ci-foundation` repository provides foundational components and configurations essential for the next-generation continuous integration (CI) infrastructure. This repository includes tools and scripts for logging configuration, monitoring, server graceful shutdown, and other infrastructure utilities.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [Logging](#logging)
    - [Graceful Shutdown](#graceful-shutdown)
    - [Utilities](#utilities)
      - [File Utilities](#file-utilities)
      - [String Utilities](#string-utilities)
      - [Array Utilities](#array-utilities)
      - [Jitter Utility](#jitter-utility)
- [Contributing](#contributing)
- [License](#license)

## Features
The `ziplines-ci-foundation` repository provides the following features:

- **Logging**: Flexible and configurable logging system with support for multiple log formats (plain text, console JSON, Cloud Logging, etc.).
- **Graceful Shutdown**: Ensure services can shut down gracefully, preserving state and avoiding data loss.
- **Utilities**: Various utility functions for working with files, strings, arrays, and jitter.

## Installation

To use the `ziplines-ci-foundation` library in your Go project, add it as a dependency in your `go.mod` file:

```sh
go get github.com/ziplinesci/ziplines-ci-foundation
```

## Usage

### Logging

The `ziplines-ci-foundation` library provides a flexible logging system that can be configured to log to different destinations (e.g., console, file, Cloud Logging). To use the logging system, import the `logging` package and create a new logger:

```go
package main

import (
	"github.com/ziplinesci/ziplines-ci-foundation/foundation/domain"
	foundation "github.com/ziplinesci/ziplines-ci-foundation/foundation/logging"
)

func main() {
	appInfo := domain.NewApplicationInfo("appGroup", "appName", "1.0.0", "main", "abc123", "2023-01-01")
	foundation.InitLoggingFromEnv(appInfo)
}
```

### Graceful Shutdown

The graceful package provides a simple and effective way to manage the shutdown of an application by:
- Tracking and waiting for all active goroutines to finish before exiting. 
- Handling OS signals (like SIGINT and SIGTERM) that indicate the application should terminate.
- Observing and logging the specific shutdown signal that triggered the termination.

Usage:
- Use NewShutdownObserver to create a new observer goroutine that will be notified when an OS signal is received via the shutdown channel (the first return value). This goroutine should handle its own cleanup as well as the cleanup of any other spawned goroutines. To inform the graceful shutdown process that cleanup is complete, it should call done() (the second return value).
- Use HandleSignals to block the main goroutine until an OS signal is received by the process.
- Use HandleSignalsWithContext to block the main goroutine until either an OS signal is received or the provided context is canceled.
- Use Shutdown to manually trigger a shutdown signal to all observers. This can be useful when you need to initiate a shutdown based on an API call or from a goroutine other than the main one.

```go
package main

import (
    "github.com/ziplinesci/ziplines-ci-foundation/foundation/graceful"
)

func main() {
    go someGoroutine()
    // if INT or TERM signal is received, go-shutdown-graceful will trigger shutdown signal to all observers.
    // Observers can do cleanup and call done() to notify go-shutdown-graceful that they are done.
    // Default timeout for cleanup is 30 seconds. This can be changed by calling HandleOsSignals with a time.Duration value.
    graceful.HandleSignals(0)
}

func someGoroutine() {
    // do something in separate goroutine
    shutdown, done := graceful.NewShutdownObserver()
    <-shutdown
    // close the background goroutine started before
    done()
}
``` 

### Utilities

The `ziplines-ci-foundation` library provides various utility functions for working with files, strings, arrays, and jitter.

#### File Utilities

The file utilities package provides functions for working with files, such as reading and writing files, watching files for changes, and creating temporary files.

#### String Utilities

The string utilities package provides functions for working with strings, such upper snake case, lowercase snake case, and more.

#### Array Utilities

The array utilities package provides functions for working with arrays, such as checking if a array contains a specific string/int, and more.

#### Jitter Utility

The jitter utility package provides functions for adding jitter to time durations.

# Contributing

Contributions are welcome! Please read our [Contributing Guidelines]() for more details on how to contribute.


# License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
