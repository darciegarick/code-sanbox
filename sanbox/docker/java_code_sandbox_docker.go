package docker

import (
	"fmt"
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

func DockerExecuteCode(exitValue request.ExecuteCodeRequest) (response.ExecuteCodeResponse, error) {
	code := exitValue.Code
	_, path, err := sanbox.SaveCodeToFile(code)
	if err != nil {
		return response.ExecuteCodeResponse{}, err
	}
	sanbox.CompileFile(path)
	fmt.Println("!!!!!!!!!!!", path)
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 镜像名称和标签
	imageName := "openjdk"
	tag := "latest"

	javaContaine, javaContaineErr := createJavaDockerContainer(client, imageName, tag, path, filepath.Join(DOCKER_CODE_PATH, DOCKER_JAVA_CLASS_PATH))
	if javaContaineErr != nil {
		return response.ExecuteCodeResponse{}, javaContaineErr
	}
	//容器id javaContaineId = javaContaine.ID

	fmt.Println("!!!!!!!!!!!", javaContaine.ID)

	return response.ExecuteCodeResponse{}, nil

}

func createJavaDockerContainer(client *docker.Client, imageName, tag, hostFilePath, containerFilePath string) (*docker.Container, error) {
	// 容器配置
	config := &docker.Config{
		Image: imageName + ":" + tag,
		// 设置命令为 执行 Main.class 文件
		Cmd: []string{"java", "-cp", DOCKER_CODE_PATH, "Main"},
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
	container, err := client.CreateContainer(createOpts)
	if err != nil {
		return nil, err
	}
	log.Printf("容器 %s 创建成功\n", container.ID)
	return container, nil
}
