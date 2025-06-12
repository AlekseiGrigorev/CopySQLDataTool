// Description: This package provides main entry point for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package main

import (
	"copysqldatatool/internal/app"
	"copysqldatatool/internal/appconfig"
	"copysqldatatool/internal/appdb"
	"copysqldatatool/internal/appfilepath"
	"copysqldatatool/internal/applog"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const ERROR = "Error:"

var Config appconfig.Config
var Log applog.AppLog
var Formatter appdb.Formatter

func prepareLogFile(logFile string) (*os.File, error) {
	if logFile == "" {
		return nil, nil
	}

	fp := appfilepath.AppFilePath{
		Path: logFile,
	}
	logFile = fp.GetWithDateTime()
	file, err := os.Create(logFile)
	if err != nil {
		Log.Error("Error creating file:", err)
		return nil, err
	}
	return file, nil
}

func main() {
	Log = applog.AppLog{}

	configFileName := flag.String("config", "config.json", "Path to the configuration file")
	logFileName := flag.String("log", "", "Path to the log file")
	flag.Parse()

	logFile, err := prepareLogFile(*logFileName)
	if logFile != nil && err == nil {
		Log.File = logFile
		defer logFile.Close()
	}

	Log.Info("Program started")
	Log.Info("Config file:", *configFileName)

	if loadConfig(*configFileName) != nil {
		Log.Error("Error loading config:", err)
		return
	}

	Formatter = appdb.Formatter{}

	for _, dataset := range Config.Datasets {
		processDataset(dataset)
	}
	Log.Ok("Program ended")
}

func loadConfig(configPath string) error {
	Config = appconfig.Config{}
	err := Config.LoadConfig(configPath)
	if err != nil {
		Log.Error("Error loading config:", err)
		return err
	}

	err = Config.Validate()
	if err != nil {
		Log.Error("Error validating config:", err)
		return err
	}

	if (Config.Datasets == nil) || (len(Config.Datasets) == 0) {
		Log.Error("No datasets found in the config")
		return fmt.Errorf("no datasets found in the config")
	}

	return nil
}

func processDataset(dataset appconfig.Dataset) {
	if !dataset.Enabled {
		Log.Warn("Skipping disabled table:", dataset.Table)
		return
	}
	if dataset.Table == "" {
		Log.Error("Skipping wrong table:", dataset.Table, "Table name is empty")
		return
	}
	if dataset.Query == "" {
		Log.Error("Skipping wrong table:", dataset.Table, "Query is empty")
		return
	}
	Log.Info("Processing table:", dataset.Table)
	err := process(Config.Config.Source, Config.Config.Dest, dataset)
	if err == nil {
		Log.Ok("Processing completed for table:", dataset.Table)
	} else {
		Log.Error("Error processing table:", dataset.Table, ERROR, err)
	}
}

func process(src appconfig.DBConfig, dst appconfig.DBConfig, dataset appconfig.Dataset) error {
	if dataset.CopyToFileEnabled() {
		Log.Info("Write to file started for table:", dataset.Table)

		file, err := createOutputFile(dataset.Table)
		if err != nil {
			Log.Error("Error creating file:", err)
			return err
		}
		defer file.Close()

		err = processRowsAndWriteToFile(src, file, dataset)
		if err != nil {
			Log.Error("Error processing rows to file for table:", dataset.Table, ERROR, err)
			return err
		}
		Log.Ok("Write to file completed for table:", dataset.Table)
	}

	if dataset.CopyToDbEnabled() {
		Log.Info("Write to db started for table:", dataset.Table)
		err := processRowsAndWriteToDb(src, dst, dataset)
		if err != nil {
			Log.Error("Error processing rows to db for table:", dataset.Table, ERROR, err)
			return err
		}
		Log.Ok("Write to db completed for table:", dataset.Table)
	}

	return nil
}

func createOutputFile(table string) (*os.File, error) {
	return os.Create(table + ".sql")
}

func processRowsAndWriteToDb(src appconfig.DBConfig, dst appconfig.DBConfig, dataset appconfig.Dataset) error {
	dataReader := createDataReader(src, dataset)
	dataReader.Open()
	err := dataReader.Open()
	if err != nil {
		Log.Error("Error opening data reader:", err)
		return err
	}
	defer dataReader.Close()

	// Connect to the destination database
	db := appdb.AppDb{
		Driver: dst.Driver,
		Dsn:    dst.DSN,
	}
	err = db.Open()
	if err != nil {
		Log.Error("Error connecting to the database:", err)
		return err
	}
	defer db.Close()

	if dataset.OnInsertSessionStart != "" {
		err = db.ExecMultiple(dataset.OnInsertSessionStart)
		if err != nil {
			Log.Error("Error executing on_insert_session_start:", err)
			return err
		}
	}

	processor := app.RowsProcessor{
		Processor:     &app.DbProcessor{AppDb: &db, Table: dataset.Table},
		DataReader:    dataReader,
		Log:           &Log,
		InsertCommand: dataset.InsertCommand,
		Table:         dataset.Table,
		Rows:          dataset.Rows,
		SqlStatement:  dataset.SqlStatement,
	}

	err = processor.Process()
	if err != nil {
		Log.Error("Error processing rows:", err)
		return err
	}

	if dataset.OnInsertSessionEnd != "" {
		err = db.ExecMultiple(dataset.OnInsertSessionEnd)
		if err != nil {
			Log.Error("Error executing on_insert_session_end:", err)
			return err
		}
	}

	return nil
}

func createDataReader(dbConf appconfig.DBConfig, dataset appconfig.Dataset) *appdb.DataReader {
	return &appdb.DataReader{
		AppDb: &appdb.AppDb{
			Driver: dbConf.Driver,
			Dsn:    dbConf.DSN,
		},
		Query:         dataset.Query,
		Type:          dataset.QueryType,
		ExecutionTime: dataset.ExecutionTime,
		InitialId:     dataset.InitialId,
	}
}

func processRowsAndWriteToFile(src appconfig.DBConfig, file *os.File, dataset appconfig.Dataset) error {
	dataReader := createDataReader(src, dataset)
	err := dataReader.Open()
	if err != nil {
		Log.Error("Error opening data reader:", err)
		return err
	}
	defer dataReader.Close()

	processor := app.RowsProcessor{
		Processor:     &app.FileProcessor{File: file},
		DataReader:    dataReader,
		Log:           &Log,
		InsertCommand: dataset.InsertCommand,
		Table:         dataset.Table,
		Rows:          dataset.Rows,
		SqlStatement:  dataset.SqlStatement,
	}

	err = processor.Process()
	if err != nil {
		Log.Error("Error processing rows:", err)
		return err
	}

	return nil
}
