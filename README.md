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

	fmt.Println("!!!!!!!!!!!!path", path)

	if err != nil {
		return
	}
	CompileFile(path)
	consoleInfo, err1 := RunFile(path)
	if err1 != nil {
		fmt.Println("!!!!!!!!!!!!err1", err1)
		return
	}
	fmt.Println("!!!!!!!!!!!!consoleInfo", consoleInfo)

	DeleteDirectory(path)

}

func ExecuteCode(exitValue request.ExecuteCodeRequest) (response.ExecuteCodeResponse, error) {
	// intputList := exitValue.InputList
	// code := exitValue.Code
	// language := exitValue.Language

	return response.ExecuteCodeResponse{}, nil // 返回结果
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




















