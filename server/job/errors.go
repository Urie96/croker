package job

import "fmt"

func errMultipleIDs(id string) error {
	return fmt.Errorf("Multiple IDs found with provided prefix: %s", id)
}

func errNoSuchJob(id string) error {
	return fmt.Errorf("No such job: %s", id)
}

func errRemoveRunningJob(id string) error {
	return fmt.Errorf("You cannot remove a running container %s. Stop the container before attempting removal or force remove", id)
}
