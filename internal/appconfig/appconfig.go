// Description: This package provides configuration management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appconfig

import (
	"copysqldatatool/internal/appdb"
	"encoding/json"
	"errors"
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
	COPY_TO_FILE       = "file"
	COPY_TO_DB         = "db"
)

// Config represents the root configuration structure
type Config struct {
	Description string     `json:"description"`
	Config      ConfigMain `json:"config"`
	Datasets    []Dataset  `json:"datasets"`
}

// ConfigDetails contains configuration details for source, destination, and default dataset
type ConfigMain struct {
	Description    string        `json:"description"`
	Source         DBConfig      `json:"source"`
	Dest           DBConfig      `json:"dest"`
	DefaultDataset DefaultConfig `json:"default_dataset"`
}

// DBConfig contains database connection details
type DBConfig struct {
	Description string `json:"description"`
	Driver      string `json:"driver"`
	DSN         string `json:"dsn"`
}

// DefaultConfig contains default dataset configuration
type DefaultConfig struct {
	Description   string `json:"description"`
	InsertCommand string `json:"insert_command"`
	Rows          int    `json:"rows"`
	CopyTo        string `json:"copy_to"`
	QueryType     string `json:"query_type"`
	SqlStatement  string `json:"sql_statement"`
	ExecutionTime int    `json:"execution_time"`
}

// Dataset represents a query and its target table
type Dataset struct {
	Description          string `json:"description"`
	Query                string `json:"query"`
	Table                string `json:"table"`
	Enabled              bool   `json:"enabled"`
	InsertCommand        string `json:"insert_command"`
	Rows                 int    `json:"rows"`
	CopyTo               string `json:"copy_to"`
	QueryType            string `json:"query_type"`
	SqlStatement         string `json:"sql_statement"`
	ExecutionTime        int    `json:"execution_time"`
	InitialId            int    `json:"initial_id"`
	OnInsertSessionStart string `json:"on_insert_session_start"`
	OnInsertSessionEnd   string `json:"on_insert_session_end"`
}

func (config *Config) Validate() error {
	messages := []string{}
	if config.Config.Source.Driver == "" {
		messages = append(messages, "source database driver cannot be empty")
	}
	if config.Config.Source.DSN == "" {
		messages = append(messages, "source database DSN cannot be empty")
	}
	if config.Config.Dest.Driver == "" {
		messages = append(messages, "destination database driver cannot be empty")
	}
	if config.Config.Dest.DSN == "" {
		messages = append(messages, "destination database DSN cannot be empty")
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
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

func (config *Config) LoadConfigFromString(str string) error {
	// Create a decoder and decode directly into the Config struct
	decoder := json.NewDecoder(strings.NewReader(str))
	err := decoder.Decode(config)
	if err != nil {
		return err
	}

	config.fillDatasets()

	return nil
}

func (config *Config) fillDatasets() {
	for i := range config.Datasets {
		config.fillDataset(i)
	}
}

func (config *Config) fillDataset(i int) {
	if config.Datasets[i].Table == "" && config.Datasets[i].Query != "" {
		sqlHelper := appdb.SqlHelper{
			Sql: config.Datasets[i].Query,
		}
		config.Datasets[i].Table = sqlHelper.GetFromTableName()
	}
	if config.Datasets[i].InsertCommand == "" {
		config.Datasets[i].InsertCommand = config.Config.DefaultDataset.InsertCommand
	}
	if config.Datasets[i].Rows == 0 {
		config.Datasets[i].Rows = config.Config.DefaultDataset.Rows
	}
	if config.Datasets[i].CopyTo == "" {
		config.Datasets[i].CopyTo = config.Config.DefaultDataset.CopyTo
	}
	if config.Datasets[i].QueryType == "" {
		config.Datasets[i].QueryType = config.Config.DefaultDataset.QueryType
	}
	if config.Datasets[i].ExecutionTime == 0 {
		config.Datasets[i].ExecutionTime = config.Config.DefaultDataset.ExecutionTime
	}
	if config.Datasets[i].SqlStatement == "" {
		config.Datasets[i].SqlStatement = config.Config.DefaultDataset.SqlStatement
	}
}

func (ds *Dataset) CopyToDbEnabled() bool {
	return strings.Contains(ds.CopyTo, COPY_TO_DB)
}

func (ds *Dataset) CopyToFileEnabled() bool {
	return strings.Contains(ds.CopyTo, COPY_TO_FILE)
}
