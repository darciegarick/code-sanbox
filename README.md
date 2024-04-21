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
	public void main(String[] args) {
		int a =  Integer.parseInt(args[0]);
		int b =  Integer.parseInt(args[1]);
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
java -cp  {编译后class文件路径} Main 1 2
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
        
    ```

















