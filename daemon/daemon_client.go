package daemon

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"

	u "github.com/jbenet/go-ipfs/util"
)

// ErrDaemonNotRunning is an error signaling that the daemon is offline
var ErrDaemonNotRunning = errors.New("daemon not running")

func getDaemonAddr(confdir string) (string, error) {
	var err error
	confdir, err = u.TildeExpansion(confdir)
	if err != nil {
		return "", err
	}
	fi, err := os.Open(confdir + "/rpcaddress")
	if err != nil {
		log.Debug("getDaemonAddr failed: %s", err)
		if err == os.ErrNotExist {
			return "", ErrDaemonNotRunning
		}
		return "", err
	}

	read := bufio.NewReader(fi)

	// TODO: operating system agostic line delim
	line, err := read.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(line), nil
}

// SendCommand attempts to run the command over a currently-running daemon.
// If there is no running daemon, returns ErrDaemonNotRunning
func SendCommand(command *Command, confdir string) error {
	//check if daemon is running
	log.Info("Checking if daemon is running...")
	var err error
	confdir, err = u.TildeExpansion(confdir)
	if err != nil {
		return err
	}
	lk, err := daemonLock(confdir)
	if err == nil {
		lk.Close()
		return ErrDaemonNotRunning
	}

	log.Info("Daemon is running! %s", err)

	server, err := getDaemonAddr(confdir)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(conn)
	err = enc.Encode(command)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, conn)

	return nil
}
