package utils

import (
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

type FollowedFile struct {
	*os.File
	watcher *fsnotify.Watcher
}

func FollowFile(path string) (file *FollowedFile, err error) {
	file = &FollowedFile{}
	if file.File, err = os.OpenFile(path, os.O_RDONLY, 0777); err != nil {
		return nil, err
	}
	if file.watcher, err = fsnotify.NewWatcher(); err != nil {
		return nil, err
	}
	if err = file.watcher.Add(path); err != nil {
		return nil, err
	}
	return file, nil
}

func (f FollowedFile) Read(b []byte) (int, error) {
	for {
		n, err := f.File.Read(b)
		if err == io.EOF { // 现有文件已经读取完毕
			select {
			case _, ok := <-f.watcher.Events: // 事件到来时唤醒，重新尝试读
				if !ok { // 通道已关闭
					return 0, io.EOF
				}
			case err, ok := <-f.watcher.Errors:
				if !ok { // 通道已关闭
					return 0, io.EOF
				}
				return 0, err
			}
		} else {
			return n, err
		}
	}
}

func (f FollowedFile) Close() error {
	f.File.Close()
	return f.watcher.Close()
}
