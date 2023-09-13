package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/robfig/cron/v3"

	"github.com/spf13/viper"
)

var (
	configPath = filepath.Join(os.Getenv("HOME"), ".config", "cron-app", "config.json")
	config     Config
	appname    = "app-cron"
)

func init() {
	envConfigPath := os.Getenv("CONFIG_PATH")
	if envConfigPath != "" {
		configPath = envConfigPath
	}

	err := LoadConfig(configPath)
	if err != nil {
		fmt.Printf("load config error %v\n", err)
	}
}

type CronJob struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}

func (cronJob CronJob) Run() {
	NewStdoutLogger().WithName(appname).Info("running command", cronJob.Command)

	// commands := strings.Split(cronJob.Command, " ")
	// if len(commands) > 0 {
	// 	commands = commands[1:] // Remove the first element
	// } else {
	// 	err := errors.New("command cron is empty")
	// 	NewStdoutLogger().WithName(appname).Error(err, "check validation command", cronJob.Command)
	// }

	// cmd := exec.Command(commands[0], commands...)
	cmd := exec.Command("/bin/bash", "-c", cronJob.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		NewStdoutLogger().WithName(appname).Error(err, "error running command", "command", cronJob.Command)
		return
	}
}

type Config struct {
	Jobs []CronJob `json:"jobs"`
}

func (c Config) LoadCrons(cron *cron.Cron) {
	for _, cronJob := range c.Jobs {
		NewStdoutLogger().WithName(appname).Info("list crons", "schedule", cronJob.Schedule, "command", cronJob.Command)
		_, err := cron.AddFunc(cronJob.Schedule, cronJob.Run)
		if err != nil {
			NewStdoutLogger().Error(err, "add new job cron")
		}
	}
}

func main() {
	NewStdoutLogger().WithName(appname).Info("starting app")
	schedule := cron.New(cron.WithLocation(time.UTC), cron.WithLogger(NewStdoutLogger()))
	go config.LoadCrons(schedule)
	defer schedule.Stop()

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// Reload konfigurasi
		if err := viper.ReadInConfig(); err != nil {
			return
		}

		if err := viper.Unmarshal(&config); err != nil {
			return
		}

		for _, entry := range schedule.Entries() {
			schedule.Remove(entry.ID)
		}

		NewStdoutLogger().WithName(appname).Info("the config has been change")
		config.LoadCrons(schedule)

	})

	go schedule.Start()

	// trap SIGINT untuk trigger shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
