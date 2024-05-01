#### ziying - 代码沙箱

###### 主流框架 & 特性
gin-gonic/gin v1.9.1
gorm v1.25.9
redis v8.11.5
mysql v1.8.1
jwt v5.2.1
gin-swagger v1.6.0
viper v1.18.2 ...


###### 代码沙箱原理
**通过命令行执行**
写入代码 => 编译代码 => 执行代码

**示例**
```java
// 写入代码
public class Main {
	public static void main(String[] args) {
		int a =  1;
		int b =  2;
		System.out.println(a + b);
	}
}
```

```shell
# 编译代码
javac -encoding UTF-8 -cp {Main.java 文件路径}
```

```shell
# 执行代码
java -cp  {编译后class文件路径} Main
```

###### 实现思路
思路：用程序代替人工，用程序来操作命令行，去编译执行代码
1. 把用户的代码保存为文件
2. 编译代码，得到 class 文件
3. 执行代码，得到输出结果
4. 收集整理输出结果
5. 文件清理，释放空间
6. 错误处理，提升代码的健壮性


###### 代码实现
1. 把用户的代码保存为文件
    - 新建目录，将每个用户的代码都存放在独立的目录下，通过UUID随机生成目录名，便于隔离和维护

    ```go
    func SaveCodeToFile(code string) (*os.File, string, error) {
		globalCodePathName := filepath.Join(SMB_CODE_PATH, GLOBAL_CODE_DIR_NAME)
		if _, err := os.Stat(globalCodePathName); err != nil {
			if mkdirErr := os.Mkdir(globalCodePathName, 0755); mkdirErr != nil {
				return nil, "", fmt.Errorf("Failed to create tmpcode folder: %s", mkdirErr)
			}
		}

		u, err := uuid.NewRandom()
		if err != nil {
			return nil, "", fmt.Errorf("Failed to generate UUID: %v", err)
		}
		userCodeParentPath := filepath.Join(globalCodePathName, u.String())

		if _, err := os.Stat(userCodeParentPath); err != nil {
			if mkdirErr := os.Mkdir(userCodeParentPath, 0755); mkdirErr != nil {
				return nil, "", fmt.Errorf("Failed to create folder: %v", mkdirErr)
			}
		}

		userCodePath := filepath.Join(userCodeParentPath, GLOBAL_JAVA_CLASS_NAME)
		file, err := os.Create(userCodePath)
		if err != nil {
			return nil, "", fmt.Errorf("Failed to create Main.java: %v", err)
		}

		// 写入代码
		_, err = file.WriteString(code)
		if err != nil {
			file.Close()
			return nil, "", fmt.Errorf("Failed to write file: %v", err)
		}
		return file, userCodeParentPath, nil
	}
    ```

2. 编译代码
使用 `os/exec` 在终端执行命令
	- 示例代码

	```go
	// 编译文件
	func CompileFile(mainJavaDirPath string) {
		mainJavaPath := filepath.Join(mainJavaDirPath, GLOBAL_JAVA_CLASS_NAME)
		compileCmd := fmt.Sprintf("javac -encoding utf-8 %s", mainJavaPath)
		cmd := exec.Command("bash", "-c", compileCmd)
		cmdErr := cmd.Run()
		if cmdErr != nil {
			fmt.Printf("Compile error: %v", cmdErr)
			return
		}
	}
	```

3. 执行文件，获得执行结果列表
	```go
	func RunFile(mainClassDirPath string) (string, error) {
		compileCmd := fmt.Sprintf("java -cp  %s Main", mainClassDirPath)
		cmd := exec.Command("bash", "-c", compileCmd)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return outputStr, nil
			} else {
				fmt.Printf("运行命令执行失败：%v\n", err)
				return "", err
			}
		}
		return outputStr, nil
	}
	```

4. 删除文件夹
	```go
	// 删除文件夹
	// 删除指定目录及其所有文件的函数
	func DeleteDirectory(dirPath string) error {
		// 获取目录下的所有文件和子目录
		files, err := filepath.Glob(filepath.Join(dirPath, "*"))
		if err != nil {
			return err
		}

		// 遍历所有文件和子目录
		for _, file := range files {
			// 如果是子目录，则递归调用删除函数
			if isDir(file) {
				err := DeleteDirectory(file)
				if err != nil {
					return err
				}
			} else {
				// 如果是文件，则直接删除
				err := os.Remove(file)
				if err != nil {
					return err
				}
			}
		}

		// 删除空目录
		err = os.Remove(dirPath)
		if err != nil {
			return err
		}
		fmt.Println("Directory deleted successfully:", dirPath)
		return nil
	}

	// 判断路径是否为目录的辅助函数
	func isDir(path string) bool {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return false
		}
		return fileInfo.IsDir()
	}
	```


**Java代码沙箱完整代码示例如下：**
```go
// /code-sanbox/test/demotest.go
package main

import (
	"fmt"
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	GLOBAL_CODE_DIR_NAME   = "tmpcode"
	GLOBAL_JAVA_CLASS_NAME = "Main.java"
	TIME_OUT               = 5000
	SMB_CODE_PATH          = "/home/ubuntu/smb-code/smb-go/code-sanbox"
)

func main() {
	code := `public class Main {
	public static void main(String[] args) {
		int a =  1;
		int b =  2;
		System.out.println(a + b);
	}
}
	`
	_, path, err := SaveCodeToFile(code)
	if err != nil {
		return
	}
	CompileFile(path)
	consoleInfo, err1 := RunFile(path)
	if err1 != nil {
		fmt.Println("!!!!!!!!!!!!err1", err1)
		return
	}
	DeleteDirectory(path)
}


// 把用户的代码保存为文件
func SaveCodeToFile(code string) (*os.File, string, error) {
	globalCodePathName := filepath.Join(SMB_CODE_PATH, GLOBAL_CODE_DIR_NAME)
	if _, err := os.Stat(globalCodePathName); err != nil {
		if mkdirErr := os.Mkdir(globalCodePathName, 0755); mkdirErr != nil {
			return nil, "", fmt.Errorf("Failed to create tmpcode folder: %s", mkdirErr)
		}
	}

	u, err := uuid.NewRandom()
	if err != nil {
		return nil, "", fmt.Errorf("Failed to generate UUID: %v", err)
	}
	userCodeParentPath := filepath.Join(globalCodePathName, u.String())

	if _, err := os.Stat(userCodeParentPath); err != nil {
		if mkdirErr := os.Mkdir(userCodeParentPath, 0755); mkdirErr != nil {
			return nil, "", fmt.Errorf("Failed to create folder: %v", mkdirErr)
		}
	}

	userCodePath := filepath.Join(userCodeParentPath, GLOBAL_JAVA_CLASS_NAME)
	file, err := os.Create(userCodePath)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create Main.java: %v", err)
	}

	// 写入代码
	_, err = file.WriteString(code)
	if err != nil {
		file.Close()
		return nil, "", fmt.Errorf("Failed to write file: %v", err)
	}
	return file, userCodeParentPath, nil
}

// 编译文件
func CompileFile(mainJavaDirPath string) {
	mainJavaPath := filepath.Join(mainJavaDirPath, GLOBAL_JAVA_CLASS_NAME)
	compileCmd := fmt.Sprintf("javac -encoding utf-8 %s", mainJavaPath)
	cmd := exec.Command("bash", "-c", compileCmd)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Printf("Compile error: %v", cmdErr)
		return
	}
}

// 执行文件，获得执行结果列表
func RunFile(mainClassDirPath string) (string, error) {
	compileCmd := fmt.Sprintf("java -cp  %s Main", mainClassDirPath)
	cmd := exec.Command("bash", "-c", compileCmd)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return outputStr, nil
		} else {
			fmt.Printf("运行命令执行失败：%v\n", err)
			return "", err
		}
	}
	return outputStr, nil
}

// 删除文件夹
// 删除指定目录及其所有文件的函数
func DeleteDirectory(dirPath string) error {
	// 获取目录下的所有文件和子目录
	files, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}

	// 遍历所有文件和子目录
	for _, file := range files {
		// 如果是子目录，则递归调用删除函数
		if isDir(file) {
			err := DeleteDirectory(file)
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，则直接删除
			err := os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	// 删除空目录
	err = os.Remove(dirPath)
	if err != nil {
		return err
	}
	fmt.Println("Directory deleted successfully:", dirPath)
	return nil
}

// 判断路径是否为目录的辅助函数
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
```

###### 优化: 代码沙箱docker实现

**前置准备**
使用 `go-dockerclient:` github.com/fsouza/go-dockerclient

**常规操作**
1. 拉去镜像
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	// 镜像名称
	imageName := "openjdk:21-jdk"
	// 镜像标签
	tag := "latest"
	// 配置拉取选项
	options := docker.PullImageOptions{
		Repository: imageName,
		Tag:        tag,
	}
	// 拉取镜像
	err = client.PullImage(options, docker.AuthConfiguration{})
	if err != nil {
		log.Fatal(err)
	}
	// 输出拉取日志
	log.Printf("镜像拉取完成：%s:%s\n", imageName, tag)
}
```

2. 创建容器
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 镜像名称和标签
	imageName := "openjdk"
	tag := "latest"
	// 容器配置
	config := &docker.Config{
		Image: imageName + ":" + tag,
		// 设置命令为查看 Java 版本的命令
		Cmd: []string{"java", "-version"},
	}
	// 创建容器请求
	createOpts := docker.CreateContainerOptions{
		Config: config,
	}
	// 创建容器
	container, err := client.CreateContainer(createOpts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("容器 %s 创建成功\n", container.ID)
}
```

3. 查看容器状态
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// 容器 ID
	containerID := "325aa5d63cf6"

	// 获取容器信息
	containerInfo, err := client.InspectContainer(containerID)
	if err != nil {
		log.Fatal(err)
	}
	// 打印容器状态信息
	log.Printf("容器状态: %s\n", containerInfo.State.StateString())
}
```

4. 启动容器
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 容器 ID
	containerID := "325aa5d63cf6"
	// 启动容器
	err = client.StartContainer(containerID, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("容器 %s 启动成功\n", containerID)
}
```

5. 查看日志
```go
package main

import (
	"log"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// 容器ID
	containerID := "325aa5d63cf6"

	// 读取容器日志
	logOptions := docker.LogsOptions{
		Container:    containerID,
		OutputStream: os.Stdout, // 将日志输出到标准输出
		ErrorStream:  os.Stderr, // 将错误信息输出到标准错误输出
		Stdout:       true,      // 获取标准输出日志
		Stderr:       true,      // 获取标准错误输出日志
		Timestamps:   true,      // 包含时间戳信息
		Follow:       true,      // 实时跟踪日志输出
	}
	err = client.Logs(logOptions)
	if err != nil {
		log.Fatal(err)
	}
}
```

6. 停止容器
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 容器ID
	containerID := "325aa5d63cf6"

	// 停止容器
	err = client.StopContainer(containerID, 10) // 第二个参数是等待超时时间，单位为秒
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("容器 %s 已停止\n", containerID)
}
```

7. 删除容器
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 容器ID
	containerID := "325aa5d63cf6"
	// 删除容器
	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    containerID,
		Force: true, // 设置为 true 表示强制删除，即使容器在运行中也会被删除
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("容器 %s 已删除\n", containerID)
}
```

8. 删除镜像
```go
package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	// 创建 Docker 客户端
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 镜像名称或ID
	imageName := "05af4ac0cbe3"
	// 删除镜像
	err = client.RemoveImage(imageName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("镜像 %s 已删除\n", imageName)
}
```


#### Docker 实现代码沙箱

1. 创建可交互Java运行容器，指定本地文件同步到容器中，让容器可以访问。

```go
// 创建可交互的 Java 运行容器
// hostFilePath: 主机文件路径 containerFilePath: 容器文件路径
func CreateJavaDockerContainer(hostFilePath string, containerFilePath string) (string, error) {
	// 创建 Docker 客户端
	endpoint := "unix:///var/run/docker.sock"
	client, clientErr := docker.NewClient(endpoint)
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
```

2. 启动容器

```go
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
```

3. 执行容器命令

```go
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
```

4. 停止并和删除容器

```go
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
```
#### Docker 代码沙箱 完整实现如下:

```go
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
```
















