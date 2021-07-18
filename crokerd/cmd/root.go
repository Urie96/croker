/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/urie96/croker/crokerd/consts"
	"github.com/urie96/croker/crokerd/http"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crokerd",
	Short: "A job manager that can run services and cronjobs in the background and view logs easily",
	Long:  `Crokerd is an open-source project to run periodic tasks and foreground services conveniently. And we can view their logs directly. All commands are similar to docker.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
	Args: cobra.NoArgs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}

func logFatal(err error) {
	if err != nil {
		log.New(os.Stderr, "", 0).Fatal(err.Error())
	}
}

func stop() {
	pid, err := getPrePid()
	logFatal(err)
	p, err := os.FindProcess(pid)
	logFatal(err)
	logFatal(p.Signal(syscall.SIGTERM))
}

func start() {
	if os.Getppid() != 1 { //判断当其是否是子进程，当父进程return之后，子进程会被 系统1 号进程接管
		// filePath, _ := filepath.Abs(os.Args[0]) //将命令行参数中执行文件路径转换成可用路径
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		// os.Remove(consts.DaemonLogPath)
		logfile, err := os.OpenFile(consts.DaemonLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		logFatal(err)
		cmd.Stdout = logfile
		cmd.Stderr = logfile
		logFatal(cmd.Start())
		return
	}
	// 守护进程
	if isRunning() {
		logFatal(errors.New("Croker is running"))
	}
	savePid()
	logFatal(http.Run())
}

func isRunning() bool {
	pid, _ := getPrePid()
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false
	}
	return true
}

func getPrePid() (int, error) {
	pidfile, err := os.Open(consts.DaemonPidPath)
	if err != nil {
		return 0, err
	}
	data, err := ioutil.ReadAll(pidfile)
	if err != nil {
		return 0, err
	}
	pidfile.Close()
	return strconv.Atoi(string(data))
}

func savePid() {
	pidfile, err := os.OpenFile(consts.DaemonPidPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	logFatal(err)
	_, err = pidfile.WriteString(strconv.Itoa(os.Getpid()))
	logFatal(err)
	pidfile.Sync()
	pidfile.Close()
}
