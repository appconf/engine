package command

//NewBash 创建一个Bash的Command实例
func NewBash(name string, args ...string) *Command {
	cmd := New("bash", "-c")
	cmd.AddArguments(name)
	cmd.AddArguments(args...)
	return cmd
}
