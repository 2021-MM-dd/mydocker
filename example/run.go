package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

/*
Start方法,首先会clone出来一个namespace隔离的进程，然后在子进程中，
调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，去初始化容器的一些资源。
*/
func Run(tty bool, cmd string) {
	parentProcess := NewParentProcess(tty, cmd)
	if err := parentProcess.Start(); err != nil {
		log.Fatal(err)
	}
	_ = parentProcess.Wait()
	os.Exit(-1)
}

//启动一个新进程
/*
这里是父进程，也就是当前进程执行的内容。
1.这里的/proc/self/exe调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境。
4.如果用户指定了-it参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, cmd string) *exec.Cmd {
	args := []string{"init", cmd}
	exec := exec.Command("/proc/self/exe", args...)
	exec.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	if tty {
		exec.Stdin = os.Stdin
		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
	}
	return exec
}

// RunContainerInitProcess 启动容器的init进程
/*
这里的init函数是在容器内部执行的，也就是说，代码执行到这里后，容器所在的进程其实就已经创建出来了，
这是本容器执行的第一个进程。
使用mount先去挂载proc文件系统，以便后面通过ps等系统命令去查看当前进程资源的情况。
*/
func RunContainerInitProcess(cmd string, args []string) error {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{cmd}
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
