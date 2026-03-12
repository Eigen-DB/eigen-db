package utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

const TEMPLATE_PATH string = "/usr/local/jails/templates/eigendb.tar"

type bsdJail struct {
	name       string
	customerId string
}

func JailFactory(customerId string) (*bsdJail, error) {
	pattern := "^[a-zA-Z0-9]{1,20}$"
	matched, err := regexp.MatchString(pattern, customerId)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("given customerId does not match this regex: %s", pattern)
	}
	return &bsdJail{
		name:       fmt.Sprintf("eigendb-%s", customerId),
		customerId: customerId,
	}, nil
}

func (j *bsdJail) jailExec(command string) ([]byte, error) {
	return exec.Command(
		"/usr/sbin/jexec",
		j.name,
		command,
	).Output()
}

func (j *bsdJail) assignPort(port uint16) error {
	out, err := j.jailExec(fmt.Sprintf("sysrc eigendb_port=%d", port))
	fmt.Println(string(out))
	return err
}

func (j *bsdJail) decompressBase() error {
	if err := os.Mkdir(fmt.Sprintf("/usr/local/jails/containers/%s", j.name), 0700); err != nil {
		return err
	}
	// decompress base env (tarball)
	return exec.Command(
		"/usr/bin/tar",
		"-xvf",
		TEMPLATE_PATH,
		"-C",
		fmt.Sprintf("/usr/local/jails/containers/%s", j.name),
		"--strip-components=1", // to get rid of the annoying top level directory
	).Run()
}

func (j *bsdJail) Start(port uint16) error {
	// check if dir already exists -> just start the jail
	if err := j.decompressBase(); err != nil {
		return err
	}
	if err := exec.Command(
		"/usr/sbin/jail",
		"-c",
		fmt.Sprintf("name=%s", j.name),
		fmt.Sprintf("host.hostname=%s", j.customerId),
		fmt.Sprintf("exec.consolelog=/var/log/jail_console_%s.log", j.name),
		"exec.start=/bin/sh /etc/rc",
		"exec.stop=/bin/sh /etc/rc.shutdown",
		"exec.clean=1",
		fmt.Sprintf("path=/usr/local/jails/containers/%s", j.name),
		"ip4=inherit",
		"interface=wlan0",
		"allow.raw_sockets=1",
		"mount.devfs=1",
	).Run(); err != nil {
		return err
	}
	if err := j.assignPort(port); err != nil {
		return err
	}
	return nil
}

func (j *bsdJail) Stop() error {
	cmd := exec.Command(
		"/usr/sbin/jail",
		"-r",
		j.name,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
