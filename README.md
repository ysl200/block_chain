# 区块链项目

## 项目结构

```
blockchain/
├── cmd/                    # 应用程序入口点
│   └── main.go            # 主程序入口
├── internal/              # 内部包，不对外暴露
│   ├── blockchain/        # 区块链核心逻辑
│   │   ├── block.go       # 区块定义
│   │   ├── chain.go       # 区块链管理
│   │   ├── transaction.go # 交易定义
│   │   └── pool.go        # 交易池
│   ├── consensus/         # 共识机制
│   │   └── raft.go        # Raft共识
│   ├── network/           # 网络层
│   │   └── node.go        # 节点管理
│   ├── storage/           # 存储层
│   │   └── store.go       # 数据存储
│   └── hash/              # 哈希算法
│       └── consistent_hash.go
├── pkg/                   # 可重用的公共包
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── metrics/           # 监控指标
│   └── scheduler/         # 调度器
├── api/                   # API层
│   ├── handler/           # HTTP处理器
│   │   └── node_handler.go
│   ├── middleware/        # 中间件
│   └── router/            # 路由
├── web/                   # Web界面
│   ├── static/            # 静态资源
│   │   ├── index.html
│   │   ├── style.css
│   │   └── app.js
│   └── server.go          # Web服务器
├── scripts/               # 脚本文件
├── docs/                  # 文档
├── tests/                 # 测试文件
├── go.mod
├── go.sum
└── README.md
```

## 主要改进

1. **清晰的层次结构**：
   - `cmd/`: 应用程序入口
   - `internal/`: 内部业务逻辑
   - `pkg/`: 可重用的公共包
   - `api/`: API接口层
   - `web/`: Web界面

2. **模块化设计**：
   - 区块链核心逻辑集中在 `internal/blockchain/`
   - 共识机制独立在 `internal/consensus/`
   - 网络层独立在 `internal/network/`

3. **标准化的命名**：
   - 使用一致的命名规范
   - 避免下划线命名

4. **测试结构优化**：
   - 测试文件与源码放在同一目录
   - 集成测试放在 `tests/` 目录

## 迁移步骤

1. 创建新的目录结构
2. 移动现有文件到对应目录
3. 更新导入路径
4. 重构代码以适应新的结构 