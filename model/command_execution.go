package model

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"io"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// CommandExecution is the model class what represents command execution
type CommandExecution struct {
	Base
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *CommandExecution) Create(db *gorm.DB, parameters gin.Params, _ url.Values, _ extension.Model) (interface{}, error) {
	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	command, exists := commandIDCommandMap[id]

	if !exists {
		return nil, errors.New("record not found")
	}

	if command.Status == commandStatusRunning {
		return nil, errors.New("the command is already running")
	}

	command.Status = commandStatusRunning
	command.Killed = false
	command.StdOut = ""
	command.StdErr = ""
	command.StartedAt = time.Now().String()
	command.FinishedAt = ""
	go executeCommand(command)

	return command, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *CommandExecution) Delete(db *gorm.DB, parameters gin.Params, _ url.Values) error {
	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	command, exists := commandIDCommandMap[id]

	if !exists {
		return errors.New("record not found")
	}

	if (command.Status != commandStatusRunning) || (command.Cmd == nil) {
		return errors.New("the command is not running")
	}

	if err := command.Cmd.Process.Kill(); err != nil {
		logging.Logger().Debug(err.Error())
		return errors.New("failed to kill command")
	}

	command.Status = commandStatusFinished
	command.FinishedAt = time.Now().String()
	command.Killed = true

	return nil
}

func runCommand(command *Command) (int, error) {
	exitCode := 0

	stdOutPipe, err := command.Cmd.StdoutPipe()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return exitCode, err
	}
	defer stdOutPipe.Close()

	stdErrPipe, err := command.Cmd.StderrPipe()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return exitCode, err
	}
	defer stdErrPipe.Close()

	var bufferStdOut, bufferStdErr bytes.Buffer
	stdOutReader := io.TeeReader(stdOutPipe, &bufferStdOut)
	stdErrReader := io.TeeReader(stdErrPipe, &bufferStdErr)

	if err = command.Cmd.Start(); err != nil {
		logging.Logger().Debug(err.Error())
		return exitCode, err
	}

	go handleStdOut(command, stdOutReader)
	go handleStdErr(command, stdErrReader)

	err = command.Cmd.Wait()

	if err != nil {
		logging.Logger().Debug(err.Error())
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
				err = nil
			} else {
				err = errors.New("syscall.WaitStatus not supported")
			}
		}
	}
	return exitCode, err
}

func handleStdOut(command *Command, reader io.Reader) {
	buffer := make([]byte, 128)
	for {
		n, err := reader.Read(buffer)
		command.StdOut = command.StdOut + string(buffer[:n])
		if err != nil {
			break
		}
	}
}

func handleStdErr(command *Command, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		command.StdErr = command.StdErr + scanner.Text()
	}
}

func executeCommand(command *Command) {

	commands := strings.Split(command.CommandLine, " ")
	command.Cmd = exec.Command(commands[0], commands[1:]...)
	command.Cmd.Dir = command.WorkingDirectory
	exitCode, err := runCommand(command)

	if err != nil {
		logging.Logger().Debug(err.Error())
		command.StdErr = err.Error()
		exitCode = -1
	}

	command.ExitCode = exitCode
	command.Status = commandStatusFinished
	command.FinishedAt = time.Now().String()
}

// NewCommandExecution creates a commandExecution model instance
func NewCommandExecution() *CommandExecution {
	return CreateModel(&CommandExecution{}).(*CommandExecution)
}

func init() {
	extension.RegisterModel(NewCommandExecution(), false)
}
