package consts

import (
	"errors"
)

var (
	HasStart   = errors.New("the job has start")
	HasStop    = errors.New("the job has stop")
	IDNotExist = errors.New("id not exist")
)

// func IDNotExist(id string) error {
// 	return fmt.Errorf("No such job: %s", id)
// }
