#### ziying - 代码沙箱

##### Projects
背景：xxx项目团队，它运营着一个广受学生青睐的在线coding平台，为学生提供了一个可自由部署项目代码环境。随着业务的扩展，xxx团队希望与知名大学、各大厂训练营和其他友商合作伙伴平台进行集成。此外，为了满足学生多样化的学习需求，xxx团队希望开放平台API，允许第三方开发者基于平台API开发在IDE产品上开发项目。

##### Product requirements
提供 API 接口供开发者调用的平台
- 管理员可以接入并发布接口
- 统计分析各接口调用情况
- xxx团队用户可以注册登录并开通接口调用权限
- 在线浏览接口、在线调试
- 可以使用客户端SDK在代码中调用接口

##### requests
1. 保证安全性，防止攻击
2. 不能随便被调用 （权限设置：限制 非社区用户无法使用 ）
3. 统计接口的调用次数
4. 计费
5. 流量保护
6. API接入 
7. ==

##### realization
- go version go1.22.1 linux/amd64
- github.com/gin-gonic/gin v1.9.1
- Redis
- MySQL
- ...


##### 主流框架 & 特性
- gin-gonic/gin v1.9.1
- gorm v1.25.9
- redis v8.11.5
- mysql v1.8.1
- jwt v5.2.1
- gin-swagger v1.6.0
- viper v1.18.2
...

##### 



