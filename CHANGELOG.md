# Changelog

All notable changes to the MCSRVR project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Support for additional server types (Spigot, Bukkit, Velocity, Forge, BungeeCord, Cuberite)

## [0.6.0] - 2025-03-03

### Changed
- Refactored server.go into multiple smaller, more manageable files
- Improved server process management to ensure servers continue running after CLI window is closed
- Enhanced process detachment on Windows using CREATE_NEW_PROCESS_GROUP and CREATE_NO_WINDOW flags
- Enhanced process detachment on Unix-like systems with improved SysProcAttr settings

## [0.5.0] - 2025-03-02

### Added
- Persistent server process tracking across command invocations
- JSON-based storage of active server processes
- Automatic server status refresh on startup
- PID-based server status detection
- Improved error handling for server management

## [0.4.0] - 2025-03-02

### Added
- RCON support for server management
- Server console access via RCON
- Command execution via RCON
- Process ID (PID) tracking for server processes
- Improved hidden process management using syscall.SysProcAttr
- Graceful server shutdown using RCON

## [0.3.0] - 2025-03-02

### Added
- Support for Fabric server type
- Fabric loader version configuration
- Automatic mods directory creation for Fabric servers
- Automatic backup functionality
- Improved error handling
- Unit and integration tests

## [0.2.0] - 2025-03-02

### Added
- Hidden process server management (systemctl-like behavior)
- Server log viewing with follow capability
- RCON configuration for server management
- Default configuration for memory and Java arguments
- Configuration editor for server settings
- Improved server status detection in list command

## [0.1.0] - 2025-03-02

### Added
- Project initialization
- Basic directory structure
- CLI command structure using Cobra
- Server initialization for PaperMC and vanilla servers
- Server management commands (start, stop, restart)
- Server console access
- Server command execution
- Server listing
- Server deletion
- Server backup and restore functionality
- Configuration management
