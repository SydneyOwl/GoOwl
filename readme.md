# GoOwl - A simple command-line CI/CD Tool
![Go](https://github.com/sydneyowl/GoOwl/actions/workflows/GoOwl_Build.yml/badge.svg) ![go_version](https://img.shields.io/badge/Go-1.18.1-brightgreen) [![Go Reference](https://pkg.go.dev/badge/github.com/sydneyowl/GoOwl.svg)](https://pkg.go.dev/github.com/sydneyowl/GoOwl) [![Release](https://img.shields.io/github/v/tag/SydneyOwl/GoOwl)](https://github.com/sydneyowl/GoOwl/releases/latest) [![card](https://goreportcard.com/badge/github.com/sydneyowl/GoOwl)](https://goreportcard.com/report/github.com/sydneyowl/GoOwl)
## What is GoOwl
GoOwl is a basic CI/CD tool. By filling a simple yaml file, you can:
+ Clone all repos automatically at the first time
+ pull spcified repos when receiving hooks
+ execute script when repos are pulled

## How to use GoOwl
 ***GoOwl is unstable currently. Do not use it in an production env and should try it in docker or virtual machine instead. Windows and Mac releases are not reliable since they're not being tested.***

Firstly you need to fill a yaml file, which is located at config/standard.yaml and could be found in this repo. It should look like this if filled correctly:

```yaml
settings:
  application:
    mode: release
    host: 0.0.0.0
    name: Hello GoOwl!
    port: 1234
  workspace:
    path: /home/golang-coder/Workspace
    repo:
      -
        id: 1 
        type: gogs
        trigger: ['push']
        repoaddr: https://website.com/GoOwl/hello.git
        username: someone
        password: 123
        buildscript: /home/golang-coder/Workspace/script/buildme2.sh
        branch: master
      -
        id: 2
        type: gogs
        trigger: ['push','pull']
        repoaddr: someone@website.com:someone/repo.git
        sshkeyaddr: /home/golang-coder/Workspace/ssh/repo
        buildscript: /script/buildme1.sh
        branch: dev
      -
        id: 3
        type: github
        trigger: ['push']
        repoaddr: someone@website.com:someone/repo1.git
        token: abcd1234
        buildscript: /script/buildme0.sh
        branch: main
```
+ `id` could be anything that could identify repos. Should be unique.
+ `workspace`.`path` defines the location repos downloaded by GoOwl storage in.
+ `buildscript` refers to the script you want to execute after specified repo is pulled. (this script will run in the directory of the repo so you don't need to use absolute addr.)
+ When using ssh, only `sshkeyaddr` and repoaddr in ssh form is needed. `username` and `password` is needed only when you need to access the repo via http(s). However, if the repo is on github, you should use token instead of `username` and `password` since github does not supports username and password authorization via http(s).
+ `branch` refers to the brance you want to clone/pull.
+ `token` should be used in ssh and only supports github now.

Ignore username,password or token if it is an public repo accessed via http(s) with correct settings..

*GoOwl supports webhook from gogs and github. More hooktypes will be supported in the future.*

GoOwl reads config from `./config/settings.yaml` by default. You can also use `-c` to specify the yaml file if you don't want to put it in default location.

Run `./GoOwl --help` to get more info.

Run `./GoOwl checkenv` to check if everything works well.  

To start the hook listener and cicd server, run `./GoOwl run`. GoOwl will automatically clone repo at the first time. You need to input "yes" if you uses ssh to clone them.

GoOwl displays hook path on start(example):
```
/gogs/1/hook---------------->Hook for repo 1,type:gogs
/gogs/2/hook---------------->Hook for repo 2,type:gogs
/github/3/hook---------------->Hook for repo 3,type:github
```
you may use `https://domain.com/gogs/1/hook` as the hook address of repo 1 for example. When GoOwl received webhook, it will start executing script automatically and print result out.

## More...
`GoOwl` may be buggy currently. Issues are welcome.

**Some of the code (common/command、cmd/run、app/other) comes from gogs and go-admin under mit license. Thanks!**
