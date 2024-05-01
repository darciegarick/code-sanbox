package docker

import (
	"bytes"
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"jassue-gin/sanbox"
	"log"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
)

const (
	GLOBAL_CODE_DIR_NAME   = "tmpcode"
	GLOBAL_JAVA_CLASS_NAME = "Main.java"
	TIME_OUT               = 5000
	SMB_CODE_PATH          = "/home/ubuntu/smb-code/smb-go/code-sanbox"
	DOCKER_CODE_PATH       = "/home/ubuntu/app"
	DOCKER_JAVA_CLASS_PATH = "Main.class"
)

func ExecuteCode(exitValue request.ExecuteCodeRequest) (response.ExecuteCodeResponse, error) {
	code := exitValue.Code
	_, path, err := sanbox.SaveCodeToFile(code)
	if err != nil {
		return response.ExecuteCodeResponse{}, err
	}
	sanbox.CompileFile(path)
	// 创建 Java 容器
	javaContainerId, javaContainerErr := CreateJavaDockerContainer(filepath.Join(path, DOCKER_JAVA_CLASS_PATH), filepath.Join(DOCKER_CODE_PATH, DOCKER_JAVA_CLASS_PATH))
	if javaContainerErr != nil {
		return response.ExecuteCodeResponse{}, javaContainerErr
	}
	// 启动 Java 容器
	StartJavaDockerContainer(javaContainerId)
	// 执行容器命令
	output, execErr := ExecuteContainerCommand(javaContainerId)
	if execErr != nil {
		return response.ExecuteCodeResponse{}, execErr
	}
	// 停止并删除容器
	StopAndRemoveContainer(javaContainerId)
	// 清理文件，释放空间
	sanbox.DeleteDirectory(path)
	return response.ExecuteCodeResponse{Result: output}, nil
}

// 创建可交互的 Java 运行容器
// hostFilePath: 主机文件路径 containerFilePath: 容器文件路径
func CreateJavaDockerContainer(hostFilePath string, containerFilePath string) (string, error) {
	// 创建 Docker 客户端
	endpoint := "unix:///var/run/docker.sock"
	client, clientErr := docker.NewClient(endpoint)
	// client, clientErr := docker.NewClientFromEnv()
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	// 镜像名称和标签
	imageName := "openjdk"
	tag := "latest"

	// 容器配置
	config := &docker.Config{
		Image:        imageName + ":" + tag,
		Tty:          true, // 开启 TTY
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash"}, // 启动 bash 以保持容器运行
	}

	// 主机配置
	hostConfig := &docker.HostConfig{
		Binds: []string{
			hostFilePath + ":" + containerFilePath, // 将本地 class 文件映射到容器中
		},
	}

	// 创建容器请求
	createOpts := docker.CreateContainerOptions{
		Config:     config,
		HostConfig: hostConfig,
	}

	// 创建容器
	container, createErr := client.CreateContainer(createOpts)
	if createErr != nil {
		return "", createErr
	}
	log.Printf("容器 %s 创建成功\n", container.ID)
	return container.ID, nil
}

// 启动 Java 容器
func StartJavaDockerContainer(containerId string) {
	// 创建 Docker 客户端
	client, clientErr := docker.NewClientFromEnv()
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	// 启动容器
	startErr := client.StartContainer(containerId, nil)
	if startErr != nil {
		log.Fatal(startErr)
	}
	log.Printf("容器 %s 启动成功\n", containerId)
}

// 执行容器命令
func ExecuteContainerCommand(containerId string) (string, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, clientErr := docker.NewClient(endpoint)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	execOpts := docker.CreateExecOptions{
		// AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"java", "-cp", "/home/ubuntu/app", "Main"},
		Container:    containerId,
	}
	// 创建 exec 实例
	exec, execErr := client.CreateExec(execOpts)
	if execErr != nil {
		log.Fatal(execErr)
	}
	// 创建 buffer 用于捕获输出
	var outBuf, errBuf bytes.Buffer

	// 执行命令
	commandErr := client.StartExec(exec.ID, docker.StartExecOptions{
		// InputStream:  nil, // 如果不需要输入可以设置为 nil
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
		RawTerminal:  true,
	})
	if commandErr != nil {
		log.Fatal(commandErr)
	}
	return outBuf.String(), nil
}

func StopAndRemoveContainer(containerId string) {
	endpoint := "unix:///var/run/docker.sock"
	client, clientErr := docker.NewClient(endpoint)
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	// 检查容器状态
	container, inspectErr := client.InspectContainer(containerId)
	if inspectErr != nil {
		log.Fatal(inspectErr)
	}

	// 如果容器正在运行，则先停止容器
	if container.State.Running {
		stopErr := client.StopContainer(containerId, 10) // 等待超时时间，单位为秒
		if stopErr != nil {
			log.Fatal(stopErr)
		}
		log.Printf("容器 %s 已停止\n", containerId)
	}

	// 删除容器
	removeErr := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    containerId,
		Force: true, // 设置为 true 表示强制删除，即使容器在运行中也会被删除
	})
	if removeErr != nil {
		log.Fatal(removeErr)
	}
	log.Printf("容器 %s 已删除\n", containerId)
}
