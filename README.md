
## Demo 项目说明

### 目录结构

- api 对外接口目录
  - proto proto文件暂存目录
- client 业务端代码目录
- data 数据接口目录
- server 服务端文件目录

### 说明

本项目为示例项目，结构和功能简单，主要作为介绍说明及学习使用

项目中使用到的一些技术说明：
1. 依赖管理：https://github.com/uber-go/fx ， 对应项目文件： **client/main.go** **server/main.go**
2. 单元测试：https://github.com/stretchr/testify , 对应项目文件: ****_test.go**
3. 数据相关测试 https://github.com/go-testfixtures/testfixtures ， 对应文件 **data/fixtures/**