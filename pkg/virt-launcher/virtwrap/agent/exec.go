package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kubevirt.io/client-go/log"

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
	Exited   *bool  `json:"exited"`
	ExitCode *int   `json:"exitcode"`
	OutData  string `json:"err-data"`
}

// ExecExitCode returned at non-zero return codes
type ExecExitCode struct {
	ExitCode int
}

func (e ExecExitCode) Error() string {
	return fmt.Sprint("exited with error code:", e.ExitCode)
}

// GuestExec sends the provided command and args to the guest agent for execution and returns an error on an unsucessful exit code
// The resulting guest stderr is returned through the existing output field.
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
	log.Log.V(4).Infof("Sending guest-exec for command %q with timeout %ds", command, timeoutSeconds)
	output, err := runCommand(cmdExec)
	if err != nil {
		log.Log.Reason(err).V(4).Infof("guest-exec failed for command %q", command)
		return "", err
	}
	log.Log.V(5).Infof("guest-exec response for command %q: %s", command, output)
	execRes := &execReturn{}
	err = json.Unmarshal([]byte(output), execRes)
	if err != nil {
		return "", fmt.Errorf("guest-exec command %q failed: invalid QEMU guest agent response: %w", command, err)
	}

	if execRes.Return.Pid <= 0 {
		return "", fmt.Errorf("guest-exec command %q failed: expected a positive PID in QEMU guest agent response: %s", command, output)
	}
	log.Log.V(4).Infof("Received guest-exec PID %d for command %q; waiting for completion", execRes.Return.Pid, command)

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
		log.Log.V(5).Infof("Sending guest-exec-status for PID %d, command %q", execRes.Return.Pid, command)
		output, err := runCommand(cmdExecStatus)
		if err != nil {
			log.Log.Reason(err).V(4).Infof("guest-exec-status failed for PID %d, command %q", execRes.Return.Pid, command)
			if timedOut() {
				return "", fmt.Errorf("Timed out waiting for guest pid [%d] for command [%s] to exit: %w", execRes.Return.Pid, command, err)
			}
			return "", err
		}
		log.Log.V(5).Infof("guest-exec-status response for PID %d, command %q: %s", execRes.Return.Pid, command, output)
		execStatusRes := &execStatusReturn{}
		err = json.Unmarshal([]byte(output), execStatusRes)
		if err != nil {
			return "", fmt.Errorf("guest-exec-status for pid [%d] failed: invalid QEMU guest agent response: %w", execRes.Return.Pid, err)
		}

		if execStatusRes.Return.Exited == nil {
			return "", fmt.Errorf("guest-exec-status for pid [%d] failed: missing exited field in QEMU guest agent response: %s", execRes.Return.Pid, output)
		}
		if *execStatusRes.Return.Exited {
			if execStatusRes.Return.ExitCode == nil || *execStatusRes.Return.ExitCode < 0 {
				return "", fmt.Errorf("guest-exec-status for pid [%d] failed: missing or invalid exit code in QEMU guest agent response: %s", execRes.Return.Pid, output)
			}
			stdOutBytes, err := base64.StdEncoding.DecodeString(execStatusRes.Return.OutData)
			if err != nil {
				return "", err
			}
			stdOut = string(stdOutBytes)
			exitCode = *execStatusRes.Return.ExitCode
			exited = true
			log.Log.V(4).Infof("Guest command %q, PID %d exited with code %d", command, execRes.Return.Pid, exitCode)
			log.Log.V(5).Infof("Guest command %q, PID %d stderr: %q", command, execRes.Return.Pid, stdOut)
			break
		}

		select {
		case <-ctx.Done():
		case <-statusCheck.C:
		}
	}

	if !exited {
		log.Log.V(4).Infof("Timed out waiting for guest PID %d, command %q", execRes.Return.Pid, command)
		return "", fmt.Errorf("Timed out waiting for guest pid [%d] for command [%s] to exit", execRes.Return.Pid, command)
	} else if exitCode != 0 {
		return stdOut, ExecExitCode{exitCode}
	}

	return stdOut, nil
}
