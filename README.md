# Croker

Croker is an open-source project to run periodic tasks and foreground services conveniently. And we can view their logs directly. All commands are similar to docker.

```shell
$ croker run -s "0 * * * *" echo One hour passed
eb80f0e4d262
$ croker run -s "@every 1s" echo One second passed
493572f389bd
$ croker ps
ID              COMMAND                 SCHEDULE        CREATED         STATUS
eb80f0e4d262    echo One hour passed    0 * * * *       12 seconds ago  Up 12 seconds
493572f389bd    echo One second p...    @every 1s       5 seconds ago   Up 5 seconds
$ croker logs 493 # view logs conveniently
2021-07-18 17:39:10  :
One second passed

2021-07-18 17:39:11  :
One second passed
$ croker stop 493
493572f389bd
$ croker run "nc -lk 8080" # run foreground service in the background
010741966e11
$ echo hello | nc localhost 8080
$ echo world | nc localhost 8080
$ croker ps
ID              COMMAND                 SCHEDULE        CREATED         STATUS
010741966e11    nc -lk 8080                             2 minutes ago   Up 2 minutes
eb80f0e4d262    echo One hour passed    0 * * * *       12 minutes ago  Up 12 minutes
493572f389bd    echo One second p...    @every 1s       12 minutes ago  Stopped 11 minutes
$ croker logs -f 010 # follow log output
hello
world
```

## Installation

```bash
go get github.com/urie96/croker/crokerd
go get github.com/urie96/croker/croker
crokerd # start the daemon progress
```

## Usage

```
Usage:
  croker [command]

Available Commands:
  inspect     A brief description of your command
  logs        Fetch the logs of a job
  prune       Remove all stopped jobs
  ps          List jobs
  rm          Remove one or more jobs
  run         Run a command or a cronjob in this host
  start       Start one or more stopped jobs
  stop        Stop one or more running jobs
  version     Show the Croker version information
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
```
