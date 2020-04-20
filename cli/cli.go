package cli

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"

	"os"
	"os/exec"
	"strings"

	"github.com/falmotlag/aws-bulk-cli/aws_util"
	"github.com/falmotlag/aws-bulk-cli/config"
	"github.com/fatih/color"
)

var outputColor = [...]color.Attribute{
	color.FgHiYellow,
	color.FgHiCyan,
	color.FgHiMagenta,
	color.FgHiBlue,
	color.FgMagenta,
	color.FgYellow,
	color.FgBlue,
	color.FgCyan,
}

// CommandExecution encapsulates all info needed
// to run a command in an account
// and the resulting output and errror
type CommandExecution struct {
	Command     []string
	AccountID   string
	RoleName    string
	OutPutColor *color.Color
	StdOut      bytes.Buffer
	StdErr      bytes.Buffer
	Err         error
}

func (c CommandExecution) String() string {
	var scanner *bufio.Scanner
	var str strings.Builder

	if c.Err != nil {
		scanner = bufio.NewScanner(&c.StdErr)

		for scanner.Scan() {
			line := scanner.Text()
			str.WriteString(fmt.Sprintf("[%v]: %v\n", c.OutPutColor.Sprint(c.AccountID), line))
		}

		str.WriteString(fmt.Sprintf("[%v]: %v\n", c.OutPutColor.Sprint(c.AccountID), c.Err.Error()))
		return str.String()
	}

	scanner = bufio.NewScanner(&c.StdOut)
	for scanner.Scan() {
		line := scanner.Text()
		str.WriteString(fmt.Sprintf("[%v]: %v\n", c.OutPutColor.Sprint(c.AccountID), line))
	}
	return str.String()
}

// take a map of key-value and return a shell friendly
// list of key=value strings
func toEnvVarsList(envVarsAsMap map[string]string) []string {
	envVarsAsList := []string{}
	for key, value := range envVarsAsMap {
		envVarsAsList = append(envVarsAsList, fmt.Sprintf("%s=%s", key, value))
	}
	return envVarsAsList
}

func executeCommand(ch chan CommandExecution, ce CommandExecution) {
	// authenticate with target aws account by assuming role
	roleArn := fmt.Sprintf("arn:aws:iam::%v:role/%v", ce.AccountID, ce.RoleName)
	creds, err := aws_util.AssumeIAmRole(roleArn)
	if err != nil {
		ce.Err = err
		ch <- ce
		return
	}

	// execute command
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmd := exec.Command("aws", ce.Command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Env = toEnvVarsList(creds)

	fmt.Printf(
		"Running %v in account [%v]...\n",
		ce.Command,
		ce.OutPutColor.Sprint(ce.AccountID),
	)

	if err := cmd.Start(); err != nil {
		ce.Err = err
		ce.StdErr = stderrBuf
		ch <- ce
		return
	}

	// wait for command to finish executing
	err = cmd.Wait()

	if err != nil {
		fmt.Printf(
			"Running %v in account [%v]... [%v]\n",
			ce.Command,
			ce.OutPutColor.Sprint(ce.AccountID),
			color.New(color.FgRed).Sprint("fail"),
		)
		ce.Err = err
		ce.StdErr = stderrBuf
		ch <- ce
		return
	}

	fmt.Printf(
		"Running %v in account [%v]... [%v]\n",
		ce.Command,
		ce.OutPutColor.Sprint(ce.AccountID),
		color.New(color.FgGreen).Sprint("success"),
	)

	ce.StdOut = stdoutBuf
	// add execution struct to channel so it's result
	// is processed on the other side
	ch <- ce
}

func Run() error {
	cfg, err := config.ConfigureApp()
	if err != nil {
		return err
	}

	// execute command in each account in list of accounts
	ch := make(chan CommandExecution, len(cfg.AWS.Accounts))
	for i, accountID := range cfg.AWS.Accounts {
		go executeCommand(
			ch,
			CommandExecution{
				Command:     flag.Args(),
				AccountID:   accountID,
				OutPutColor: color.New(outputColor[i]),
				RoleName:    cfg.AWS.CrossAccountRoleName,
			},
		)
	}

	// collect result, wait for all accounts to return
	results := make([]CommandExecution, 0)
	for range cfg.AWS.Accounts {
		results = append(results, <-ch)
	}

	for _, r := range results {
		fmt.Println(r.String())
	}

	return nil
}
