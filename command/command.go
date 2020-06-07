package command

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"
)

//Command 运行的Shell指令需要参数
type Command struct {
	name    string
	args    []string
	envs    []string
	dir     string
	timeout time.Duration
}

//New 新建一个Command实例。name为命令路径，args为命令参数
func New(name string, args ...string) *Command {
	return &Command{
		name: name,
		args: args,
	}
}

//SetTimeout 设置指令运行超时时间
func (cmd *Command) SetTimeout(timeout time.Duration) *Command {
	cmd.timeout = timeout
	return cmd
}

//RunDir 设置运行指令时的工作目录
func (cmd *Command) RunDir(dir string) *Command {
	cmd.dir = dir
	return cmd
}

//AddEnvs 添加环境变量，如: PATH="/bin:/sbin"
func (cmd *Command) AddEnvs(envs ...string) *Command {
	cmd.envs = append(cmd.envs, envs...)
	return cmd
}

//AddArguments 为运行的指令添加参数
func (cmd *Command) AddArguments(args ...string) *Command {
	cmd.args = append(cmd.args, args...)
	return cmd
}

//RunWithPipe 运行指令，并将命令运行时的标准输出，和标准错误分别输出到stdout, stderr
//如果指令运行失败，会返回err
func (cmd *Command) RunWithPipe(stdout, stderr io.Writer) (err error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if cmd.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), cmd.timeout)
	} else {
		ctx, cancel = context.WithCancel(context.TODO())
	}
	defer cancel()

	command := exec.CommandContext(ctx, cmd.name, cmd.args...)
	command.Stderr = stderr
	command.Stdout = stdout
	command.Env = append(os.Environ(), cmd.envs...)
	command.Dir = cmd.dir

	return command.Run()
}

//Run 运行命令，并将命令的标准输出，标准错误内容合并到output返回
func (cmd *Command) Run() (output []byte, err error) {
	buffer := new(bytes.Buffer)
	err = cmd.RunWithPipe(buffer, buffer)
	return buffer.Bytes(), err
}
