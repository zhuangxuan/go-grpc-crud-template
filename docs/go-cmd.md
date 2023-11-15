### 文档
- https://gorm.io/zh_CN/docs
- https://gin-gonic.com/zh-cn/docs

### 创建项目
```bash
go mod init
```

### 安装依赖
```bash
# go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/mysql
go get -u gorm.io/gorm
go get -u github.com/gin-gonic/gin

go get -u github.com/golang-jwt/jwt/v5
```

### 运行
```bash
go build main.go # 单个文件
go build # 整个文件夹
```