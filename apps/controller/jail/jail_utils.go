package jail

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"crypto/rand"
)

type jlsOutput struct {
	Version         string `json:"__version"`
	JailInformation struct {
		Jail []BsdJail `json:"jail"`
	} `json:"jail-information"`
}

func ListInstances(customerId string) ([]string, error) {
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

	instances := make([]string, 0)
	for _, jail := range jails.JailInformation.Jail {
		cid, instanceId, err := parseJailName(jail.Hostname)
		if cid == customerId {
			instances = append(instances, instanceId)
		}
		if err != nil {
			return nil, err
		}
	}
	return instances, nil
}

func generateJailName(customerId string) (string, error) {
	if err := validateCustomerId(customerId); err != nil {
		return "", err
	}
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	jailName := fmt.Sprintf("eigendb-%s-%s", customerId, hex.EncodeToString(bytes))
	return jailName, nil
}

func validateCustomerId(customerId string) error {
	// GitHub's valid username regex
	regex := "^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$"
	matched, err := regexp.MatchString(regex, customerId)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("given customerId does not match this regex: %s", regex)
	}
	return nil
}

// parseJailName("eigendb-john-smith-1e43f2") -> ("john-smith", "1e43f2", nil)
func parseJailName(jName string) (string, string, error) {
	return "", "", nil
}

func getJailMetadata(jName string) (*BsdJail, error) {
	var jail jlsOutput
	outputJson, err := exec.Command(
		"/usr/sbin/jls",
		"-j",
		jName,
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
