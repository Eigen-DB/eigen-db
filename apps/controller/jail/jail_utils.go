package jail

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type jlsOutput struct {
	Version         string `json:"__version"`
	JailInformation struct {
		Jail []BsdJail `json:"jail"`
	} `json:"jail-information"`
}

func ListInstances(customerId string) ([]Instance, error) {
	var jails jlsOutput
	outputJson, err := exec.Command(
		"/usr/sbin/jls",
		"--libxo",
		"json",
	).Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(
		outputJson,
		&jails,
	); err != nil {
		return nil, err
	}

	customerJails := make([]Instance, 0)
	for _, jail := range jails.JailInformation.Jail {
		if strings.HasSuffix(jail.Hostname, "-eigendb") && strings.TrimSuffix(jail.Hostname, "-eigendb") == customerId {
			instance, err := InstanceFactory(customerId, true)
			if err != nil {
				return nil, err
			}
			customerJails = append(customerJails, *instance)
		}
	}
	return customerJails, nil
}

func generateCustomerId(customerName string) (string, error) {
	customerId := fmt.Sprintf("%s-eigendb", customerName)
	if err := validateCustomerId(customerId); err != nil {
		return "", err
	}
	return customerId, nil
}

func validateCustomerId(customerId string) error {
	pattern := "^[a-zA-Z0-9]{1,20}-eigendb$" // to prevent OS command injection
	matched, err := regexp.MatchString(pattern, customerId)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("given customerId does not match this regex: %s", pattern)
	}
	return nil
}

func getJailMetadata(customerId string) (*BsdJail, error) {
	var jail jlsOutput
	outputJson, err := exec.Command(
		"/usr/sbin/jls",
		"-j",
		customerId,
		"--libxo",
		"json",
	).Output()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(
		outputJson,
		&jail,
	); err != nil {
		return nil, err
	}
	return &jail.JailInformation.Jail[0], nil
}
