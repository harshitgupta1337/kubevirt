package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/cli"
)

type execReturn struct {
	Return execReturnData `json:"return"`
}
type execReturnData struct {
	Pid int `json:"pid"`
}

type execStatusReturn struct {
	Return execStatusReturnData `json:"return"`
}
type execStatusReturnData struct {
	Exited   bool   `json:"exited"`
	ExitCode int    `json:"exitcode"`
	OutData  string `json:"out-data"`
}

// ExecExitCode returned at non-zero return codes
type ExecExitCode struct {
	ExitCode int
}

func (e ExecExitCode) Error() string {
	return fmt.Sprint("exited with error code:", e.ExitCode)
}

// GuestExec sends the provided command and args to the guest agent for execution and returns an error on an unsucessful exit code
// The resulting stdout will be returned as a string
func GuestExec(virConn cli.Connection, domName string, command string, args []string, timeoutSeconds int32) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	return GuestExecWithCommand(ctx, command, args, timeoutSeconds, func(cmd string) (string, error) {
		return virConn.QemuAgentCommand(cmd, domName)
	})
}

// GuestExecWithCommand executes through the supplied guest agent transport.
func GuestExecWithCommand(ctx context.Context, command string, args []string, timeoutSeconds int32, runCommand func(string) (string, error)) (string, error) {
	if timeoutSeconds <= 0 {
		return "", fmt.Errorf("guest execution timeout must be positive")
	}
	stdOut := ""
	encodedCommand, err := json.Marshal(command)
	if err != nil {
		return "", err
	}
	encodedArgs := make([]string, 0, len(args))
	for _, arg := range args {
		encodedArg, err := json.Marshal(arg)
		if err != nil {
			return "", err
		}
		encodedArgs = append(encodedArgs, string(encodedArg))
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}
	cmdExec := fmt.Sprintf(`{"execute": "guest-exec", "arguments": { "path": %s, "arg": [ %s ], "capture-output":true } }`, encodedCommand, strings.Join(encodedArgs, ", "))
	output, err := runCommand(cmdExec)
	if err != nil {
		return "", err
	}
	execRes := &execReturn{}
	err = json.Unmarshal([]byte(output), execRes)
	if err != nil {
		return "", err
	}

	if execRes.Return.Pid <= 0 {
		return "", fmt.Errorf("Invalid pid [%d] returned from qemu agent during access credential injection: %s", execRes.Return.Pid, output)
	}

	exited := false
	exitCode := 0
	statusCheck := time.NewTicker(time.Duration(timeoutSeconds) * 100 * time.Millisecond)
	defer statusCheck.Stop()
	checkUntil, hasDeadline := ctx.Deadline()
	timedOut := func() bool {
		return ctx.Err() != nil || (hasDeadline && !time.Now().Before(checkUntil))
	}
	for {
		if timedOut() {
			break
		}
		cmdExecStatus := fmt.Sprintf(`{"execute": "guest-exec-status", "arguments": { "pid": %d } }`, execRes.Return.Pid)
		output, err := runCommand(cmdExecStatus)
		if err != nil {
			if timedOut() {
				return "", fmt.Errorf("Timed out waiting for guest pid [%d] for command [%s] to exit: %w", execRes.Return.Pid, command, err)
			}
			return "", err
		}
		execStatusRes := &execStatusReturn{}
		err = json.Unmarshal([]byte(output), execStatusRes)
		if err != nil {
			return "", err
		}

		if execStatusRes.Return.Exited {
			stdOutBytes, err := base64.StdEncoding.DecodeString(execStatusRes.Return.OutData)
			if err != nil {
				return "", err
			}
			stdOut = string(stdOutBytes)
			exitCode = execStatusRes.Return.ExitCode
			exited = true
			break
		}

		select {
		case <-ctx.Done():
		case <-statusCheck.C:
		}
	}

	if !exited {
		return "", fmt.Errorf("Timed out waiting for guest pid [%d] for command [%s] to exit", execRes.Return.Pid, command)
	} else if exitCode != 0 {
		return stdOut, ExecExitCode{exitCode}
	}

	return stdOut, nil
}
