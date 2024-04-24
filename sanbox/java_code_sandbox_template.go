package sanbox

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

func ExecuteCode(exitValue request.ExecuteCodeRequest) (response.ExecuteCodeResponse, error) {
	code := exitValue.Code
	_, path, err := saveCodeToFile(code)
	if err != nil {
		return response.ExecuteCodeResponse{}, err
	}
	compileFile(path)
	result, runFileErr := runFile(path)
	if runFileErr != nil {
		return response.ExecuteCodeResponse{}, runFileErr
	}
	deleteDirectory(path)
	return response.ExecuteCodeResponse{Result: result}, nil
}

// 把用户的代码保存为文件
func saveCodeToFile(code string) (*os.File, string, error) {
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
func compileFile(mainJavaDirPath string) {
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
func runFile(mainClassDirPath string) (string, error) {
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
func deleteDirectory(dirPath string) error {
	// 获取目录下的所有文件和子目录
	files, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}

	// 遍历所有文件和子目录
	for _, file := range files {
		// 如果是子目录，则递归调用删除函数
		if isDir(file) {
			err := deleteDirectory(file)
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
