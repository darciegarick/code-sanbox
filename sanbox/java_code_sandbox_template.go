package sanbox

import (
	"fmt"
	"jassue-gin/entity/request"
	"jassue-gin/entity/response"
	"os"
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
	return file, userCodePath, nil
}

// 编译文件
func CompileFile(string) (string) {

	return ""

}

// func CompileFile(file *os.File) (*os.File, error) {
// 	mainJavaPath, err := getFilePath(file)

// 	fmt.Println("!!!!!mainJavaPath:", mainJavaPath)

// 	if err != nil {
// 		return nil, err
// 	}
// 	compileCmd := fmt.Sprintf("javac -encoding utf-8 %s", mainJavaPath)
// 	cmd := exec.Command("bash", "-c", compileCmd)

// 	cmdErr := cmd.Run()
// 	if cmdErr != nil {
// 		return nil, cmdErr
// 	}

// 	fmt.Println("Compile success!!!!!!!!!!!!!!!!!!!")
// 	return file, nil

// }
