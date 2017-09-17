// +build !windows,!solaris

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/docker/libnetwork/portallocator"
	"github.com/maliceio/engine/cmd/maliced/hack"
	"golang.org/x/sys/unix"
)

const defaultDaemonConfigFile = "/etc/malice/daemon.json"

// setDefaultUmask sets the umask to 0022 to avoid problems
// caused by custom umask
func setDefaultUmask() error {
	desiredUmask := 0022
	unix.Umask(desiredUmask)
	if umask := unix.Umask(desiredUmask); umask != desiredUmask {
		return fmt.Errorf("failed to set umask: expected %#o, got %#o", desiredUmask, umask)
	}

	return nil
}

func getDaemonConfDir(_ string) string {
	return "/etc/malice"
}

// setupConfigReloadTrap configures the USR2 signal to reload the configuration.
func (cli *DaemonCli) setupConfigReloadTrap() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGHUP)
	go func() {
		for range c {
			cli.reloadConfig()
		}
	}()
}

// func (cli *DaemonCli) getPlatformRemoteOptions() []libcontainerd.RemoteOption {
// 	opts := []libcontainerd.RemoteOption{
// 		libcontainerd.WithDebugLog(cli.Config.Debug),
// 		libcontainerd.WithOOMScore(cli.Config.OOMScoreAdjust),
// 	}
// 	if cli.Config.ContainerdAddr != "" {
// 		opts = append(opts, libcontainerd.WithRemoteAddr(cli.Config.ContainerdAddr))
// 	} else {
// 		opts = append(opts, libcontainerd.WithStartDaemon(true))
// 	}
// 	if daemon.UsingSystemd(cli.Config) {
// 		args := []string{"--systemd-cgroup=true"}
// 		opts = append(opts, libcontainerd.WithRuntimeArgs(args))
// 	}
// 	if cli.Config.LiveRestoreEnabled {
// 		opts = append(opts, libcontainerd.WithLiveRestore(true))
// 	}
// 	opts = append(opts, libcontainerd.WithRuntimePath(daemon.DefaultRuntimeBinary))
// 	return opts
// }

// allocateDaemonPort ensures that there are no containers
// that try to use any port allocated for the malice server.
func allocateDaemonPort(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	var hostIPs []net.IP
	if parsedIP := net.ParseIP(host); parsedIP != nil {
		hostIPs = append(hostIPs, parsedIP)
	} else if hostIPs, err = net.LookupIP(host); err != nil {
		return fmt.Errorf("failed to lookup %s address in host specification", host)
	}

	pa := portallocator.Get()
	for _, hostIP := range hostIPs {
		if _, err := pa.RequestPort(hostIP, "tcp", intPort); err != nil {
			return fmt.Errorf("failed to allocate daemon listening port %d (err: %v)", intPort, err)
		}
	}
	return nil
}

// notifyShutdown is called after the daemon shuts down but before the process exits.
func notifyShutdown(err error) {
}

func wrapListeners(proto string, ls []net.Listener) []net.Listener {
	switch proto {
	case "unix":
		ls[0] = &hack.MalformedHostHeaderOverride{ls[0]}
	case "fd":
		for i := range ls {
			ls[i] = &hack.MalformedHostHeaderOverride{ls[i]}
		}
	}
	return ls
}