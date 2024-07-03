# HikHello Project

## Overview
HikHello is a Go-based application designed to interface with HIKAX devices, providing a comprehensive solution for monitoring and managing device statuses, including sirens. It utilizes a polling mechanism to fetch data from devices and updates the status in real-time. The application supports both HTTP and MQTT protocols for data polling, ensuring wide compatibility and flexibility in deployment.

## Features
- Real-time device status monitoring
- Support for multiple device types, with a focus on sirens
- Data fetching using a configurable polling mechanism
- HTTP and MQTT support for versatile data polling
- Secure authentication with HIKAX devices

## Prerequisites
- Go (version 1.15 or later)
- Access to HIKAX devices
- MQTT broker (if using MQTT polling)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/hikhello.git
   cd hikhello
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build Mac OS binary for the project:
   ```bash
   make
   ```
4. Build docker image for the project:
   ```bash
   make build-docker
    ```

## Configuration
Before running HikHello, configure the application by editing the `options` struct in the main.go file or by providing command-line arguments.

- `HIKAX.Host`: Specify the host of the HIKAX device.
- `HIKAX.Username`: Username for device authentication.
- `HIKAX.Password`: Password for device authentication.
- `PollingTime`: Interval for polling device status (in seconds).

## Running the Application
To start the application, simply run:
```bash
./hikhello --hikax.host=hikax_ip --hikax.username=hikax_username --hikax.password=hikax_password
```

## Usage
Once running, HikHello will begin polling connected HIKAX devices based on the specified interval. The application logs will provide real-time feedback on the polling process and any changes in device statuses.
You can wiew devices status by accessing the following URL: http://localhost:8080/ by default or specify the host and port using the --listen flag.

## Contributing
Contributions to HikHello are welcome! Please feel free to submit pull requests or open issues to discuss proposed changes or report bugs.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
```

This README provides a basic structure for your project, including sections for an overview, features, prerequisites, installation instructions, configuration details, running the application, usage, contributing, and licensing. Adjust the content as necessary to fit the specifics of your project and repository URL.