# GoTakeOut 外卖点餐系统

GoTakeOut 是一个基于 Golang 开发的外卖点餐系统后端服务，提供完整的餐厅菜品管理、用户订单处理、支付集成等功能。

## 功能特性

- **用户模块**：用户注册、登录、个人信息管理
- **地址管理**：用户收货地址的增删改查
- **菜品管理**：菜品分类、菜品信息维护、口味管理
- **套餐管理**：套餐组合创建与管理
- **购物车**：商品添加、修改、删除
- **订单系统**：下单、支付、订单状态跟踪
- **员工管理**：员工账户管理、权限控制
- **支付集成**：微信支付接口对接

## 技术栈

- **框架**：Gin Web 框架
- **数据库**：MySQL + GORM
- **缓存**：Redis
- **认证**：JWT
- **日志**：Zap + Lumberjack
- **配置**：Viper
- **支付**：微信支付 API
- **存储**：阿里云 OSS

## 项目结构

```
GoTakeOut/
├── common/           # 公共模块
├── config.yaml       # 应用配置文件
├── config-env.yaml   # 环境变量配置
├── internal/         # 主要模块
├── logs/             # 日志文件
├── main.go           # 主程序入口
├── model/            # 数据模型
│   ├── entity/       # 数据库实体
│   ├── dto/          # 数据传输对象
│   ├── vo/           # 视图对象
│   └── wrap/         # 包装对象
├── router/           # 路由定义
│   ├── admin/        # 管理员路由
│   ├── user/         # 用户路由
│   └── notify/       # 通知路由
└── template/         # 模板文件
```

## 快速开始

### 环境要求

- Go 1.22+
- MySQL
- Redis

### 安装与配置

1. 克隆仓库

```bash
git clone https://github.com/bdly-mrxh/GoTakeOut.git
cd GoTakeOut
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置环境变量

编辑 `config-env.yaml` 文件，配置数据库、Redis、JWT 密钥等信息。

4. 运行应用

```bash
go run main.go
```

应用将在配置的端口（默认 8080）启动。


## 许可证

[MIT License](LICENSE) 