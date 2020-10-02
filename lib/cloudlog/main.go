package cloudlog

// This tool transmit error with an universal single struct to
// - Google Logging
// - Google ErrorReporting (if severity is critical or emergency)
// Hence all error follow the same format accross all services
import (
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	envKeyProjectID   = "CLOUDLOG_PROJECTID"
	envKeyServiceName = "CLOUDLOG_SERVICE_NAME"
	envKeyEnvironment = "BRUSHED-CHARTS-ENVIRONMENT"
)

// Variables set accross the package to avoid duplicate inforamtion
// and to ease the log writing
var (
	projectID   string
	serviceName string
)

// Critical send report to Google ErrorReporting
// and print also on stderr
func Critical(err error) error {
	var entry LogEntry
	entry.initFromError(err)
	entry.print()
	err = report(entry)

	return err
}

// CriticalWithMessage same as `Critical` but concat the given
// `message` before the given error (`err`)
func CriticalWithMessage(err error, message string) error {
	var entry LogEntry
	entry.initFromError(err)
	entry.Error = fmt.Errorf("%s\n%v", message, entry.Error)
	entry.print()
	err = report(entry)
	return err
}

// Panic is the same as calling Critical then panic. It will
// log the error on the cloud and locally then panic with empty
// body
func Panic(err error) error {
	Critical(err)
	panic("")
}

func report(entry LogEntry) error {
	if err := tryToInit(); err != nil {
		log.Printf("%+v", err)
		return err
	}

	err := googleErrorReporting(entry)
	return err
}

func tryToInit() error {
	projectID = os.Getenv(envKeyProjectID)
	serviceName = os.Getenv(envKeyServiceName)
	if projectID == "" || serviceName == "" {
		return errors.New("cloudlogging is not initialized")
	}
	return nil
}