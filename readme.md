# GoOwl - A simple command-line CI/CD Tool
![Go](https://github.com/sydneyowl/GoOwl/actions/workflows/GoOwl_Build.yml/badge.svg) [![card](https://goreportcard.com/badge/github.com/sydneyowl/GoOwl)](https://goreportcard.com/report/github.com/sydneyowl/GoOwl) [![CodeQL](https://github.com/SydneyOwl/GoOwl/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/SydneyOwl/GoOwl/actions/workflows/codeql-analysis.yml)  ![go_version](https://img.shields.io/badge/Go-1.18.1-brightgreen) [![Go Reference](https://pkg.go.dev/badge/github.com/sydneyowl/GoOwl.svg)](https://pkg.go.dev/github.com/sydneyowl/GoOwl) [![Release](https://img.shields.io/github/v/tag/SydneyOwl/GoOwl)](https://github.com/sydneyowl/GoOwl/releases/latest) 
## What is GoOwl
GoOwl is a basic CI/CD tool. By filling a simple yaml file, you can:
- [x]  clone all repos automatically at the first time
- [x]  pull spcified repos when receiving hooks
- [x]  execute script when repos are pulled
- [ ]  return a badge showing build status

## How to use GoOwl
 ***~GoOwl is unstable currently. Do not use it in an production env and should try it in docker or virtual machine instead.~ Windows and Mac releases are not reliable since they're not being tested.***

Firstly you need to fill a yaml file, which is located at config/standard.yaml and could be found in this repo. It should look like this if filled correctly:

```yml
settings:
  application:
    mode: release
    host: 0.0.0.0
    name: testApp
    port: 1234
  workspace:
    path: /home/golang-coder/Workspace
    repo:
      -
        id: gogs-test-token
        type: gogs
        trigger: ['push']
        repoaddr: https://git....
        token: abc
        buildscript: /home/golang-coder/Workspace/script/test.sh
        branch: master
      -
        id: gogs-test-up-userpass
        type: gogs
        trigger: push
        repoaddr: https://git
        username: 123
        password: 123
        buildscript: 
        branch: master
      -
        id: gogs-test-ssh
        type: gogs
        trigger: ['push']
        repoaddr: git@git...
        sshkeyaddr: /home/golang-coder/Workspace/id_rsa
        buildscript: 
        branch: master
      -
        id: github-test-token
        type: github
        trigger: push
        repoaddr: https://git
        token: ghp_....
        buildscript: /home/golang-coder/Workspace/script/test1.sh
        branch: main
      -
        id: github-test-ssh
        type: github
        trigger: push
        repoaddr: git@git...
        sshkeyaddr: /home/golang-coder/Workspace/id_rsa_2
        branch: main
```
+ `id` could be anything that could identify repos. Should be unique.
+ `workspace`.`path` defines the location repos downloaded by GoOwl storage in.
+ `buildscript` refers to the script you want to execute after specified repo is pulled. (this script will run in the directory of the repo so you don't need to use absolute addr.)
+ When using ssh, only `sshkeyaddr` and repoaddr in ssh form is needed. `username` and `password` is needed only when you need to access the repo via http(s). **However, if the repo is on github, you should use token instead of `username` and `password` since github does not supports username and password authorization via http(s).**
+ `branch` refers to the branch you want to clone/pull.

Ignore username,password or token if it is an public repo accessed via http(s) with correct settings.

*GoOwl supports webhook from gogs and github. More hooktypes will be supported in the future.*

>**Authorization methods supported:**
>
>Gogs: ssh(privatekey)/http(username&password)/http(token)
>
>Github: ssh(privatekey)/http(token)
>
>**Authorization order by dafault(both exists)**
>
>ssh->oauth->http

GoOwl reads config from `./config/settings.yaml` by default. You can also use `-c` to specify the yaml file if you don't want to put it in default location.

Run `./GoOwl --help` to get more info.

Run `./GoOwl checkenv` to check if everything works well.  

Use `-l [1,2,3]` or `--log [1,2,3]` to use log func. 1 means output to stdout, 2 means output to file only(workspace/log/...), 3 means output to both file and stdout. They all storaged in workspace.path/log specified in the yml by default.

To start the hook listener and cicd server, run `./GoOwl run`. GoOwl will automatically clone repo at the first time. You need to input "yes" if you uses ssh to clone them. If you'd like to ignore repo checking(whether config is properly filled), use `--skip-repocheck`. 

GoOwl displays hook path on start(example):
```
/gogs/1/hook---------------->Hook for repo 1,type:gogs
/gogs/2/hook---------------->Hook for repo 2,type:gogs
/github/3/hook---------------->Hook for repo 3,type:github
```
you may use `https://domain:port/gogs/1/hook` as the hook address of repo 1 for example. When GoOwl received webhook, it will start executing script automatically and print result out. However, if repo failed to clone on start, the route of the repo won't be registered, which means the hook of the repo is unavailable. 

full example:
![](./md_pic/1.png)
![](./md_pic/2.png)

## More...
`GoOwl` may be buggy currently. Issues are welcome.

**Some part (/config/command and /app/other) comes from [gogs](https://github.com/gogs/git-module) and [go-admin](https://github.com/go-admin-team/go-admin) under MIT license. Thanks!**
