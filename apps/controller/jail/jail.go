package jail

// TODO: ADD MORE LOGS TO MAKE EVERYTHING SUPER DEBUGGABLE

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const TEMPLATE_PATH string = "/usr/local/jails/templates/eigendb.tar"

// add more metadata
type Instance struct {
	customerId        string
	jail              *BsdJail
	userlandInstalled bool
}

type BsdJail struct {
	Jid      int    `json:"jid"`
	IPv4     string `json:"ipv4"`
	Hostname string `json:"hostname"`
	Path     string `json:"path"`
}

type Port = uint16

func InstanceFactory(customerName string, userlandInstalled bool) (*Instance, error) {
	customerId, err := generateCustomerId(customerName)
	if err != nil {
		return nil, err
	}
	return &Instance{
		customerId:        customerId,
		userlandInstalled: userlandInstalled,
		jail:              nil, // the jail pointer starts nil by default and is populated later
	}, nil
}

func BsdJailFactory(customerId string) (*BsdJail, error) {
	if err := validateCustomerId(customerId); err != nil {
		return nil, err
	}
	return getJailMetadata(customerId)
}

// use carefully
func (i *Instance) jailExec(command ...string) ([]byte, error) {
	command = append([]string{i.customerId}, command...)
	return exec.Command(
		"/usr/sbin/jexec",
		command...,
	).Output()
}

func (i *Instance) assignPort(port Port) error {
	out, err := i.jailExec(fmt.Sprintf("sysrc eigendb_port=%d", port))
	fmt.Println(string(out))
	return err
}

func (i *Instance) decompressUserland() error {
	if err := os.Mkdir(fmt.Sprintf("/usr/local/jails/containers/%s", i.customerId), 0700); err != nil {
		return err
	}
	// decompress base env (tarball)
	if err := exec.Command(
		"/usr/bin/tar",
		"-xvf",
		TEMPLATE_PATH,
		"-C",
		fmt.Sprintf("/usr/local/jails/containers/%s", i.customerId),
		"--strip-components=1", // to get rid of the annoying top level directory
	).Run(); err != nil {
		return err
	}
	i.userlandInstalled = true
	return nil
}

func (i *Instance) Create() error {
	if err := i.decompressUserland(); err != nil {
		return err
	}
	return nil
}

func (i *Instance) Start(port Port) error {
	if !i.userlandInstalled {
		return errors.New("the instance must first be created using .Create()")
	}
	if err := exec.Command(
		"/usr/sbin/jail",
		"-c",
		fmt.Sprintf("name=%s", i.customerId),
		fmt.Sprintf("host.hostname=%s", i.customerId),
		fmt.Sprintf("exec.consolelog=/var/log/jail_console_%s.log", i.customerId),
		"exec.start=/bin/sh /etc/rc",
		"exec.stop=/bin/sh /etc/rc.shutdown",
		"exec.clean=1",
		fmt.Sprintf("path=/usr/local/jails/containers/%s", i.customerId),
		"ip4=inherit",
		"interface=wlan0",
		"allow.raw_sockets=1",
		"mount.devfs=1",
	).Run(); err != nil {
		return err
	}
	if err := i.assignPort(port); err != nil {
		return err
	}

	// BUG: race condition -> BsdJailFactory is executed before the jail is done starting
	// populate i.jail
	j, err := BsdJailFactory(i.customerId)
	if err != nil {
		return err
	}
	i.jail = j

	// start the EigenDB in the jail?
	// seems like it starts automatically

	return nil
}

func (i *Instance) Stop() error {
	cmd := exec.Command(
		"/usr/sbin/jail",
		"-r",
		i.customerId,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	i.jail = nil
	return nil
}
