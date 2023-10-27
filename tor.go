package torgo

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/net/proxy"
)

type TorProxy struct {
	cmd     *exec.Cmd
	options *Options
	forward dialer
	proxy   ContextDialer
	url     *url.URL
}

func NewTorProxy(opts *Options) (*TorProxy, error) {
	if !CheckInstalledTor() {
		return nil, fmt.Errorf("tor not installed")
	}

	return &TorProxy{
		options: opts,
	}, nil
}

func (t *TorProxy) Start(ctx context.Context) error {
	tmpFile, err := os.CreateTemp("", "tor-config")
	if err != nil {
		return fmt.Errorf("Cannot create temp file: %v", err)
	}
	os.Chmod(tmpFile.Name(), 0600)

	// Close the file after the Tor command completes
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(t.options.Serialize())
	if err != nil {
		return fmt.Errorf("Cannot write to temp file: %v", err)
	}
	tmpFile.Close()

	cmd := exec.Command("tor", "-f", tmpFile.Name())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to create stdout pipe: %v", err)
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start Tor: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	bootstrapped := make(chan bool)

	go func() {
		for scanner.Scan() {
			if t.options.Debug {
				log.Println(scanner.Text())
			}

			line := scanner.Text()
			if strings.Contains(line, "Opened Socks listener connection (ready)") {
				bootstrapped <- true
				return
			}
			if strings.Contains(line, "error") {
				bootstrapped <- false
				cmd.Process.Signal(os.Interrupt)
				return
			}
		}
		bootstrapped <- false
	}()

	select {
	case success := <-bootstrapped:
		if success {
			t.cmd = cmd
			const socksFormat = "socks5://localhost:%d"
			socks5URL, err := url.Parse(fmt.Sprintf(socksFormat, t.options.GeneralOptions.SocksPort))
			if err != nil {
				return fmt.Errorf("failed to parse socks5 url: %v", err)
			}

			t.url = socks5URL

			return nil
		}
		return fmt.Errorf("failed to bootstrap tor")
	case <-ctx.Done():
		cmd.Process.Signal(os.Interrupt)
		return ctx.Err()
	}
}

func (t *TorProxy) GetProxy() ContextDialer {
	forward, err := proxy.FromURL(t.url, proxy.Direct)
	if err != nil {
		return nil
	}

	t.forward = forward

	t.proxy = t.forward.(ContextDialer)
	runtime.SetFinalizer(t, (*TorProxy).Close)

	return t.proxy
}

func (t *TorProxy) GetURL() *url.URL {
	return t.url
}

func (t *TorProxy) Close() error {
	if t.cmd != nil {
		return t.cmd.Process.Signal(os.Interrupt)
	}
	return nil
}
