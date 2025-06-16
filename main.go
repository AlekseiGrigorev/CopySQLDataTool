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

// Main is the main entry point of the application.
// It reads the configuration file and processes each dataset by calling processDataset.
// The function logs the start and end of the program, config file name, and any errors encountered.
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

// prepareLogFile prepares a log file by creating a new file with the current date and time in its name.
// If the file name is empty, the function returns nil and no error.
// Otherwise, it creates a new file and returns the file object and an error.
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

// loadConfig initializes the global Config variable by loading and validating the configuration
// from the provided file path. It logs any errors encountered during the loading or validation
// process. If no datasets are found in the configuration, it logs an error and returns an error.
// Returns an error if the configuration cannot be loaded, validated, or contains no datasets.
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

// processDataset processes a single dataset by first checking its enabled status,
// table name, and query validity. If the dataset is disabled, has an empty table name,
// or an empty query, it logs a warning or error and returns without processing.
// If valid, it logs the start of processing, calls the process function to handle
// the dataset, and logs the result of the processing.
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

// process handles the data processing for a given dataset by checking its configuration
// and performing the necessary actions based on the dataset's settings. It supports
// writing data to a file or a database, or both, depending on the dataset's CopyTo
// configuration. The function initializes the data reader, manages file creation,
// connects to the destination database, and executes the data processing logic, while
// logging the progress and any errors encountered. Returns an error if any step fails.
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

// createOutputFile creates a new file with the given table name and ".sql" extension
// to write the output SQL statements. It returns the opened file and an error if
// any.
func createOutputFile(table string) (*os.File, error) {
	return os.Create(table + ".sql")
}

// processRowsAndWriteToFile processes rows from a source database and writes them to a specified file.
// It initializes a data reader using the provided database configuration and dataset information,
// and uses a RowsProcessor to manage the data transfer. The function handles opening and closing
// the data reader, logging errors, and ensuring the proper execution of the data processing logic.
// It returns an error if any step in the process fails, such as opening the data reader or processing rows.
func processRowsAndWriteToFile(src appconfig.DBConfig, file *os.File, dataset appconfig.Dataset) error {
	dataReader := createDataReader(src, dataset)
	err := dataReader.Open()
	if err != nil {
		Log.Error("Error opening data reader:", err)
		return err
	}
	defer dataReader.Close()

	processor := app.RowsProcessor{
		Processor:  &app.FileProcessor{File: file},
		DataReader: dataReader,
		Log:        &Log,
		Dataset: app.Dataset{
			InsertCommand:    dataset.InsertCommand,
			TableName:        dataset.Table,
			RowsPerCommand:   dataset.Rows,
			SqlStatementType: dataset.SqlStatement,
		},
	}

	err = processor.Process()
	if err != nil {
		Log.Error("Error processing rows:", err)
		return err
	}

	return nil
}

// processRowsAndWriteToDb processes rows from a source database and writes them to a destination database.
// It initializes a data reader using the provided source database configuration and dataset information,
// connects to the destination database using the provided configuration, and uses a RowsProcessor to manage
// the data transfer. The function handles opening and closing the data reader, connecting to the destination database,
// logging errors, and ensuring the proper execution of the data processing logic. It returns an error if any step in
// the process fails, such as opening the data reader, connecting to the destination database, or processing rows.
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
		Processor:  &app.DbProcessor{AppDb: &db, TableName: dataset.Table},
		DataReader: dataReader,
		Log:        &Log,
		Dataset: app.Dataset{
			InsertCommand:    dataset.InsertCommand,
			TableName:        dataset.Table,
			RowsPerCommand:   dataset.Rows,
			SqlStatementType: dataset.SqlStatement,
		},
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

// createDataReader creates a new DataReader instance using the provided database
// configuration and dataset information. It configures the DataReader with the
// database connection details, query, query type, execution time, and initial ID.
// The function returns a pointer to the newly created DataReader instance.
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
