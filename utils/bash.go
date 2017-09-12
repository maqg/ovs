package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"text/template"

	"github.com/Sirupsen/logrus"

	"github.com/pkg/errors"
)

// Bash Command for shell
type Bash struct {
	Command   string
	PipeFail  bool
	Arguments map[string]string
	NoLog     bool

	retCode int
	stdout  string
	stderr  string
	err     error
}

func (b *Bash) build() error {
	Assert(b.Command != "", "Command cannot be emptry string")

	if b.Arguments != nil {
		tmpl, err := template.New("script").Parse(b.Command)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, b.Arguments)
		if err != nil {
			return err
		}

		b.Command = buf.String()
	}

	if b.PipeFail {
		b.Command = fmt.Sprintf("set -o pipefail; %s", b.Command)
	}

	return nil
}

// Run for bash command
func (b *Bash) Run() error {
	ret, so, se, err := b.RunWithReturn()
	if err != nil {
		return fmt.Errorf("failed to execute the command[%s] because of an internal errro", b.Command)
	}

	if ret != 0 {
		return fmt.Errorf("failed to exectue the command[%s],return code:%d,stdout:%s,stderr:%s",
			b.Command, ret, so, se)
	}
	return nil
}

// RunWithReturn Run Command with Reture value
func (b *Bash) RunWithReturn() (retCode int, stdout, stderr string, err error) {
	if err = b.build(); err != nil {
		b.err = err
		return -1, "", "", err
	}

	if !b.NoLog {
		logrus.Debugf("shell start: %s", b.Command)
	}

	var so, se bytes.Buffer
	var cmd *exec.Cmd

	if len(b.Command) > 1024*4 {
		content := []byte(b.Command)
		tmpfile, err := ioutil.TempFile("/home/vyos", "zvrcommand")
		PanicOnError(err)
		err = tmpfile.Chmod(0777)
		PanicOnError(err)
		_, err = tmpfile.Write(content)
		PanicOnError(err)
		err = tmpfile.Close()
		PanicOnError(err)
		cmd = exec.Command("bash", "-c", tmpfile.Name())
		defer os.Remove(tmpfile.Name())
	} else {
		cmd = exec.Command("bash", "-c", b.Command)
	}

	cmd.Stdout = &so
	cmd.Stderr = &se

	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			retCode = waitStatus.ExitStatus()
		} else {
			panic(errors.Errorf("unable to get return code, %s", err))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		retCode = waitStatus.ExitStatus()
	}

	stdout = string(so.Bytes())
	stderr = string(se.Bytes())

	b.retCode = retCode
	b.stdout = stdout
	b.stderr = stderr

	if !b.NoLog {
		logrus.WithFields(logrus.Fields{
			"return code": fmt.Sprintf("%v", retCode),
			"stdout":      stdout,
			"stderr":      stderr,
		}).Debugf("shell done: %s", b.Command)
	}

	return
}

// PanicIfError for bashing
func (b *Bash) PanicIfError() {
	if b.err != nil {
		panic(fmt.Errorf("shell failure[command: %v], internal error: %v",
			b.Command, b.err))
	}

	if b.retCode != 0 {
		panic(fmt.Errorf("shell failure[command: %v, return code: %v, stdout: %v, stderr: %v",
			b.Command, b.retCode, b.stdout, b.stderr))
	}
}

// NewBash for Bash Comand
func NewBash() *Bash {
	return &Bash{}
}
