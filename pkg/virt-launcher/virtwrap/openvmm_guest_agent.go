/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors.
 *
 */

package virtwrap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"time"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/agent"
)

func (l *OpenVMMDomainManager) Exec(domainName, command string, args []string, timeoutSeconds int32) (string, error) {
	if timeoutSeconds <= 0 {
		return "", fmt.Errorf("guest execution timeout must be positive")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	select {
	case l.guestExecLock <- struct{}{}:
		defer func() { <-l.guestExecLock }()
	case <-ctx.Done():
		return "", fmt.Errorf("waiting for OpenVMM guest agent: %w", ctx.Err())
	}

	l.mutex.RLock()
	running := l.state == openVMMRunning && l.domain != nil && l.domain.Spec.Name == domainName
	socketPath := filepath.Join(l.socketsDir, openVMMCOM2Socket)
	l.mutex.RUnlock()
	if !running {
		return "", fmt.Errorf("OpenVMM domain %q is not running", domainName)
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
	if err != nil {
		return "", fmt.Errorf("connecting to OpenVMM guest agent at %s: %w", socketPath, err)
	}
	defer conn.Close()
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		return "", fmt.Errorf("setting OpenVMM guest agent deadline: %w", err)
	}

	client := &openVMMGuestAgentClient{conn: conn}
	if err := client.synchronize(); err != nil {
		return "", fmt.Errorf("synchronizing OpenVMM guest agent: %w", err)
	}
	return agent.GuestExecWithCommand(ctx, command, args, timeoutSeconds, client.command)
}

type openVMMGuestAgentClient struct {
	conn    net.Conn
	decoder *json.Decoder
}

func (c *openVMMGuestAgentClient) synchronize() error {
	id := time.Now().UnixNano()
	// The delimiter resets a partial request in the agent and marks the response
	// after any stale data left on the serial stream by a previous connection.
	request := fmt.Sprintf("\xff{\"execute\":\"guest-sync-delimited\",\"arguments\":{\"id\":%d}}\n", id)
	if _, err := io.WriteString(c.conn, request); err != nil {
		return err
	}
	reader := bufio.NewReader(c.conn)
	for {
		for {
			b, err := reader.ReadByte()
			if err != nil {
				return err
			}
			if b == 0xff {
				break
			}
		}
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		response, err := parseOpenVMMGuestAgentResponse(line)
		if err != nil {
			return err
		}
		var result struct {
			Return int64 `json:"return"`
		}
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			return err
		}
		if result.Return == id {
			c.decoder = json.NewDecoder(reader)
			return nil
		}
	}
}

func (c *openVMMGuestAgentClient) command(request string) (string, error) {
	if _, err := io.WriteString(c.conn, request+"\n"); err != nil {
		return "", fmt.Errorf("sending OpenVMM guest agent command: %w", err)
	}
	response, err := c.readResponse()
	if err != nil {
		return "", fmt.Errorf("reading OpenVMM guest agent response: %w", err)
	}
	return response, nil
}

func (c *openVMMGuestAgentClient) readResponse() (string, error) {
	var response json.RawMessage
	if err := c.decoder.Decode(&response); err != nil {
		return "", err
	}
	return parseOpenVMMGuestAgentResponse(response)
}

func parseOpenVMMGuestAgentResponse(response json.RawMessage) (string, error) {
	var envelope struct {
		Return json.RawMessage `json:"return"`
		Error  *struct {
			Class string `json:"class"`
			Desc  string `json:"desc"`
		} `json:"error"`
	}
	if err := json.Unmarshal(response, &envelope); err != nil {
		return "", err
	}
	if envelope.Error != nil {
		return "", fmt.Errorf("QEMU guest agent error %s: %s", envelope.Error.Class, envelope.Error.Desc)
	}
	if len(envelope.Return) == 0 || string(envelope.Return) == "null" {
		return "", fmt.Errorf("QEMU guest agent response has no return value: %s", response)
	}
	return string(response), nil
}
