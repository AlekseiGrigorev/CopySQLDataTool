package main

import (
	"copysqldatatool/internal/appdb"
	"copysqldatatool/internal/applog"
	"copysqldatatool/internal/config"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const ERROR = "Error:"

var Config config.Config
var Log applog.AppLog

func main() {
	Log.Info("Program started")
	configFile := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()
	Log.Info("Config file:", *configFile)
	Config = config.Config{}
	err := Config.LoadConfig(*configFile)
	if err != nil {
		Log.Error("Error loading config:", err)
		return
	}

	if (Config.Datasets == nil) || (len(Config.Datasets) == 0) {
		Log.Error("No datasets found in the config")
		return
	}

	for _, dataset := range Config.Datasets {
		if !dataset.Enabled {
			Log.Warn("Skipping disabled table:", dataset.Table)
			continue
		}
		if dataset.Table == "" {
			Log.Error("Skipping wrong table:", dataset.Table, "Table name is empty")
			continue
		}
		if dataset.Query == "" {
			Log.Error("Skipping wrong table:", dataset.Table, "Query is empty")
			continue
		}
		Log.Info("Processing table:", dataset.Table)
		err = Process(Config.Config.Source, Config.Config.Dest, dataset)
		if err == nil {
			Log.Ok("Processing completed for table:", dataset.Table)
		} else {
			Log.Error("Error processing table:", dataset.Table, ERROR, err)
		}
	}
	Log.Ok("Program ended")
}

func Process(src config.DBConfig, dst config.DBConfig, dataset config.Dataset) error {
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

func processRowsAndWriteToDb(src config.DBConfig, dst config.DBConfig, dataset config.Dataset) error {
	dataReader := createDataReader(src, dataset)
	dataReader.Open()
	defer dataReader.Close()

	// Connect to the destination database
	db := appdb.AppDb{
		Driver: dst.Driver,
		Dsn:    dst.DSN,
	}
	err := db.Open()
	if err != nil {
		Log.Error("Error connecting to the database:", err)
		return err
	}
	defer db.Close()

	var buffer []string
	data := make([]any, 0)
	count := 0
	rowsCount := 0
	columns := make([]string, 0)

	if dataset.OnInsertSessionStart != "" {
		err = db.ExecMultiple(dataset.OnInsertSessionStart)
		if err != nil {
			Log.Error("Error executing on_insert_session_start:", err)
			return err
		}
	}

	for dataReader.Next() {
		if rowsCount == 0 {
			columns = dataReader.WrappedColumns()
		}

		values, err := dataReader.Scan()
		if err != nil {
			Log.Error("Error scanning row:", err)
			return err
		}

		insertStatement := getInsertStatement(values, dataset)
		buffer = appendRowToBuffer(buffer, dataset, columns, insertStatement, count)
		if dataset.Statement == config.STATEMENT_PREPARED {
			data = append(data, values...)
		}
		count++
		rowsCount++

		if count == dataset.Rows {
			if err := writeBufferToDb(&db, buffer, data); err != nil {
				Log.Error("Error writing buffer to database:", err)
				return err
			}
			buffer = nil
			data = make([]any, 0)
			count = 0
			Log.Info("Rows processed to table", dataset.Table, "...:", rowsCount)
		}
	}

	// Handle any remaining rows
	if len(buffer) > 0 {
		if err := writeBufferToDb(&db, buffer, data); err != nil {
			Log.Error("Error writing buffer to database:", err)
			return err
		}
	}

	if dataset.OnInsertSessionEnd != "" {
		err = db.ExecMultiple(dataset.OnInsertSessionEnd)
		if err != nil {
			Log.Error("Error executing on_insert_session_end:", err)
			return err
		}
	}

	Log.Ok("Rows processed to table", dataset.Table, ":", rowsCount)

	return nil
}

func createDataReader(dbConf config.DBConfig, dataset config.Dataset) *appdb.DataReader {
	return &appdb.DataReader{
		AppDb: appdb.AppDb{
			Driver: dbConf.Driver,
			Dsn:    dbConf.DSN,
		},
		Query:         dataset.Query,
		Type:          dataset.Type,
		ExecutionTime: dataset.ExecutionTime,
		InitialId:     dataset.InitialId,
	}
}

func processRowsAndWriteToFile(src config.DBConfig, file *os.File, dataset config.Dataset) error {
	dataReader := createDataReader(src, dataset)
	dataReader.Open()
	defer dataReader.Close()

	var buffer []string
	count := 0
	rowsCount := 0
	columns := make([]string, 0)

	for dataReader.Next() {
		if rowsCount == 0 {
			columns = dataReader.WrappedColumns()
		}

		values, err := dataReader.Scan()
		if err != nil {
			Log.Error("Error scanning row:", err)
			return err
		}

		insertStatement := getInsertStatement(values, dataset)
		buffer = appendRowToBuffer(buffer, dataset, columns, insertStatement, count)
		count++
		rowsCount++

		if count == dataset.Rows {
			if err := writeBufferToFile(file, buffer); err != nil {
				Log.Error("Error writing buffer to file:", err)
				return err
			}
			buffer = nil
			count = 0
			Log.Info("Rows processed to file", file.Name(), "...:", rowsCount)
		}
	}

	if len(buffer) > 0 {
		if err := writeBufferToFile(file, buffer); err != nil {
			Log.Error("Error writing buffer to file:", err)
			return err
		}
	}

	Log.Ok("Rows processed to file", file.Name(), ":", rowsCount)
	return nil
}

func appendRowToBuffer(buffer []string, dataset config.Dataset, columns []string, insertStatement string, count int) []string {
	if count == 0 {
		buffer = appendInitialInsert(buffer, dataset.InsertCommand, dataset.Table, columns, insertStatement)
	} else {
		buffer = append(buffer, fmt.Sprintf(", (%s)", insertStatement))
	}
	return buffer
}

func getInsertStatement(values []any, dataset config.Dataset) string {
	if dataset.Statement == config.STATEMENT_PREPARED {
		return buildInsertPlaceholders(len(values))
	}
	return formatRowValues(values)
}

func buildInsertPlaceholders(columnCount int) string {
	return strings.Repeat("?, ", columnCount-1) + "?"
}

func formatRowValues(values []interface{}) string {
	var formattedValues []string

	for _, val := range values {
		formattedValues = append(formattedValues, formatSingleValue(val))
	}

	return strings.Join(formattedValues, ", ")
}

func formatSingleValue(val any) string {
	switch v := val.(type) {
	case []byte:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(string(v), "'", "''"))
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(string(v), "'", "''"))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func appendInitialInsert(buffer []string, command string, table string, columns []string, insertStatement string) []string {
	columnsStr := strings.Join(columns, ", ")
	insertCommand := fmt.Sprintf("%s %s (%s) VALUES", command, table, columnsStr)
	buffer = append(buffer, insertCommand)
	buffer = append(buffer, fmt.Sprintf("(%s)", insertStatement))
	return buffer
}

func writeBufferToFile(file *os.File, buffer []string) error {
	buffer = append(buffer, ";")
	for _, stmt := range buffer {
		_, err := file.WriteString(stmt + "\n")
		if err != nil {
			Log.Error("Error writing to file:", err)
			return err
		}
	}
	return nil
}

func writeBufferToDb(db *appdb.AppDb, buffer []string, data []any) error {
	buffer = append(buffer, ";")
	_, err := db.Exec(strings.Join(buffer, ""), data...)
	if err != nil {
		Log.Error("Error writing to database:", err)
		return err
	}
	return nil
}
