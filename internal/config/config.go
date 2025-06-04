package config

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	TYPE_UNDEFINED     = ""
	TYPE_SIMPLE        = "simple"
	TYPE_LIMIT_OFFSET  = "limitoffset"
	TYPE_ORDERBYID     = "orderbyid"
	STATEMENT_PREPARED = "prepared"
	STATEMENT_RAW      = "raw"
)

// Config represents the root configuration structure
type Config struct {
	Config   ConfigMain `json:"config"`
	Datasets []Dataset  `json:"datasets"`
}

// ConfigDetails contains configuration details for source, destination, and default dataset
type ConfigMain struct {
	Source         DBConfig      `json:"source"`
	Dest           DBConfig      `json:"dest"`
	DefaultDataset DefaultConfig `json:"default_dataset"`
}

// DBConfig contains database connection details
type DBConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

// DefaultConfig contains default dataset configuration
type DefaultConfig struct {
	InsertCommand string `json:"insert_command"`
	Rows          int    `json:"rows"`
	CopyTo        string `json:"copy_to"`
	Type          string `json:"type"`
	ExecutionTime int    `json:"execution_time"`
	Statement     string `json:"statement"`
}

// Dataset represents a query and its target table
type Dataset struct {
	Query                string `json:"query"`
	Table                string `json:"table"`
	Enabled              bool   `json:"enabled"`
	InsertCommand        string `json:"insert_command"`
	Rows                 int    `json:"rows"`
	CopyTo               string `json:"copy_to"`
	Type                 string `json:"type"`
	ExecutionTime        int    `json:"execution_time"`
	Statement            string `json:"statement"`
	InitialId            int    `json:"initial_id"`
	OnInsertSessionStart string `json:"on_insert_session_start"`
	OnInsertSessionEnd   string `json:"on_insert_session_end"`
}

func (config *Config) LoadConfig(path string) error {
	// Open the config file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a decoder and decode directly into the Config struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}

	config.fillDatasets()

	return nil
}

func (config *Config) fillDatasets() {
	for i := range config.Datasets {
		if config.Datasets[i].InsertCommand == "" {
			config.Datasets[i].InsertCommand = config.Config.DefaultDataset.InsertCommand
		}
		if config.Datasets[i].Rows == 0 {
			config.Datasets[i].Rows = config.Config.DefaultDataset.Rows
		}
		if config.Datasets[i].CopyTo == "" {
			config.Datasets[i].CopyTo = config.Config.DefaultDataset.CopyTo
		}
		if config.Datasets[i].Type == "" {
			config.Datasets[i].Type = config.Config.DefaultDataset.Type
		}
		if config.Datasets[i].ExecutionTime == 0 {
			config.Datasets[i].ExecutionTime = config.Config.DefaultDataset.ExecutionTime
		}
		if config.Datasets[i].Statement == "" {
			config.Datasets[i].Statement = config.Config.DefaultDataset.Statement
		}
	}
}

func (ds *Dataset) CopyToDbEnabled() bool {
	return strings.Contains(ds.CopyTo, "db")
}

func (ds *Dataset) CopyToFileEnabled() bool {
	return strings.Contains(ds.CopyTo, "file")
}
