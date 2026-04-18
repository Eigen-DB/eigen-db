package utils

import (
	"controller/jail"
)

// decompresses the jail's userland
func CreateInstance(customerName string) error {
	instance, err := jail.InstanceFactory(customerName, false)
	if err != nil {
		return err
	}
	return instance.Create()
}

// instance must be created before started
func StartInstance(customerName string) error {
	// verify that the customer Id is valid. try not using an external database and instead get data from the BSD system to avoid bugs.
	instance, err := jail.InstanceFactory(customerName, true)
	if err != nil {
		return err
	}

	// jail network configuration logic goes here...
	var port jail.Port = 6789

	return instance.Start(port)
}

func StopInstance(customerName string) error {
	instance, err := jail.InstanceFactory(customerName, true)
	if err != nil {
		return err
	}
	return instance.Stop()
}

func ListInstances(customerName string) ([]jail.Instance, error) {
	instances, err := jail.ListInstances(customerName)
	if err != nil {
		return []jail.Instance{}, err
	}
	return instances, nil
}
