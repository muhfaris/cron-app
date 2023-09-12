# Cron App

Cron App is a simple command-line utility written in Go that allows you to schedule and execute commands at specified intervals using the cron syntax. This README provides an overview of the project, how to build and run it, and how to configure scheduled tasks.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Building the Project](#building-the-project)
  - [Running the Application](#running-the-application)
  - [Running the Application using systemd](#running-the-applincation-using-systemd)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

Before you can build and run Cron App, you'll need the following:

- [Go](https://golang.org/dl/): Make sure you have Go installed on your system.

### Building the Project

To build the project, follow these steps:

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/muhfaris/cron-app.git
   ```

2. Change to the project directory:

   ```bash
   cd cron-app
   ```

3. Build the application:

   ```bash
   go build -o cron-app
   ```

### Running the Application

You can run Cron App using the following command:

```bash
./cron-app
```

### Running the Application using systemd

To run your Go application as a systemd service on a Linux system, you can create a systemd service unit file. This allows you to manage your application as a service, enabling automatic startup, monitoring, and graceful shutdown. Here's how you can set up a systemd service for your Go application:

1. **Move cron-app to /usr/local/bin**

   ```bash
   sudo mv cron-app /usr/local/bin/cron-app
   ```

2. **Create a systemd Service Unit File:**

   Open a terminal and create a systemd service unit file with a `.service` extension, such as `cron-app.service`. You should typically place this file in the `/etc/systemd/system/` directory, but you can also use `/etc/systemd/user/` for user-specific services. Use a text editor or a command-line tool to create the file, for example:

   ```bash
   sudo nano /etc/systemd/system/cron-app.service
   ```

3. **Define the Service Unit:**

   Add the following content to your `cron-app.service` file, replacing the placeholders with the appropriate values:

   ```plaintext
   [Unit]
   Description=Cron App Service
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/cron-app
   WorkingDirectory=/home/muhfaris/.config/cron-app
   Restart=always
   RestartSec=3
   Environment=CONFIG_PATH=/home/muhfaris/.config/cron-app/config.json
   StandardOutput=syslog
   StandardError=syslog

   [Install]
   WantedBy=multi-user.target

   ```

   - `Description`: A description for your service.
   - `ExecStart`: The path to your Go application's executable.
   - `WorkingDirectory`: The directory where your application should run.
   - `Restart`: Specifies when the service should be restarted (in this example, always).
   - `RestartSec`: The time to wait before restarting the service.
   - `Environment`: You can set environment variables if your application requires them.
   - `StandardOutput` and `StandardError`: Redirects standard output and standard error to syslog for logging.

4. **Reload systemd Configuration:**

   After creating the service unit file, reload the systemd configuration to make it aware of the new service:

   ```bash
   sudo systemctl daemon-reload
   ```

5. **Start and Enable the Service:**

   Start the service and enable it to start at boot:

   ```bash
   sudo systemctl start cron-app
   sudo systemctl enable cron-app
   ```

6. **View Service Status and Logs:**

   You can check the status of your service with:

   ```bash
   sudo systemctl status cron-app
   ```

   example response:

   ```bash
    cron-app.service - Cron App Service
         Loaded: loaded (/etc/systemd/system/cron-app.service; disabled; preset: enabled)
         Active: active (running) since Wed 2023-09-13 06:51:22 WIB; 4s ago
       Main PID: 525450 (cron-app)
          Tasks: 6 (limit: 18715)
         Memory: 1.2M
            CPU: 9ms
         CGroup: /system.slice/cron-app.service
                 └─525450 /usr/local/bin/cron-app

    Sep 13 06:51:22 ichiro systemd[1]: Started cron-app.service - Cron App Service.
    Sep 13 06:51:22 ichiro cron-app[525450]: load config error open .config/cron-app/config.json: no such file or directory
    Sep 13 06:51:22 ichiro cron-app[525450]: app-cron: "ts"="2023-09-13 06:51:22.697648" "level"=0 "msg"="starting app"
    Sep 13 06:51:22 ichiro cron-app[525450]: "ts"="2023-09-13 06:51:22.700255" "level"=0 "msg"="start"
   ```

   To view the application's logs, you can use the `journalctl` command:

   ```bash
   journalctl -u cron-app
   ```

Your Go application should now be running as a systemd service, and it will start automatically at boot and restart in case of crashes. You can manage the service using standard systemd commands like `start`, `stop`, `restart`, and `status`.

## Configuration

Cron App uses a configuration file to define the scheduled tasks. The default configuration file is located at `~/.config/cron-app/config.json`. You can customize the configuration by editing this file.

The configuration file has the following structure:

```json
{
  "jobs": [
    {
      "schedule": "*/5 * * * *", // Cron schedule expression
      "command": "echo 'Hello, World!'" // Command to be executed
    },
    {
      "schedule": "0 0 * * *",
      "command": "backup.sh"
    }
    // Add more scheduled tasks here
  ]
}
```

- `schedule`: This field specifies the cron schedule expression that determines when the command will be executed. You can use standard cron syntax to define the schedule.

- `command`: The command to be executed at the scheduled time. You can specify any valid shell command or script here.

After editing the configuration file, you can restart Cron App to apply the changes. The application will automatically reload the configuration when it detects changes to the file.

## Usage

Cron App will load the scheduled tasks from the configuration file and execute them according to the defined schedules. The application runs in the background and continues to execute tasks at their scheduled times.

You can stop the application gracefully by sending a SIGINT (Ctrl+C) or SIGTERM signal to it. This will trigger a shutdown, and the application will stop executing tasks.

## Contributing

Contributions to Cron App are welcome! If you find a bug or have a feature request, please open an issue on the GitHub repository. If you'd like to contribute code, feel free to fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
