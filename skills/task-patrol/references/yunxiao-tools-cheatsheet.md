# Yunxiao MCP Tools 快速参考

## 用户相关

| 工具 | 用途 |
|------|------|
| `get_current_user` | 获取当前用户 ID、姓名、默认组织 |

## 任务查询

| 工具 | 用途 | 关键参数 |
|------|------|---------|
| `get_my_project_workitems` | 聚合查询某用户的任务（按分类） | `projectId`, `userId`, `relation` |
| `search_workitems` | 灵活搜索任务 | `projectId`, `category`, `assignedTo`, `status`, `subject` |
| `get_project_workitem_board` | Kanban 视图 | `projectId` |
| `get_project_workitem_detail` | 单个任务详情 | `workitemId` |

## 项目上下文

| 工具 | 用途 | 关键参数 |
|------|------|---------|
| `search_projects` | 搜索用户参与的项目列表 | `organizationId`, `perPage` |
| `get_project_overview` | 项目快照 | `projectId` |
| `list_labels` | 项目标签列表 | `projectId` |
| `get_project_risk_dashboard` | 风险任务识别 | `projectId` |

## 常用调用顺序

**单项目模式**：
```
1. get_current_user -> 获取 userId
2. get_my_project_workitems -> 获取任务列表
3. (可选) get_project_workitem_detail -> 深入具体任务
4. (可选) list_labels -> 标签一致性检查
```

**跨项目模式**：
```
1. get_current_user -> 获取 userId
2. search_projects -> 获取用户参与的项目列表
3. 对每个项目调用 get_my_project_workitems -> 获取各项目任务
4. 聚合数据 -> 统一视图
```

## 状态值参考

> ⚠️ **重要**：状态 ID 因项目配置而异，下表仅为示例格式，**实际使用时必须通过 API 获取当前项目的状态映射**，不可硬编码。
>
> 获取方式：`get_project_workitem_board` 或观察 `search_workitems` 返回数据中的 `status` 字段。

示例格式（非实际值）：

| 状态语义 | 示例状态名 | 说明 |
|---------|-----------|------|
| 未排期 | `backlog` / `待排期` | 尚未安排迭代的任务 |
| 待开始 | `待处理` / `Todo` / `待办` | 已排期但未开始 |
| 进行中 | `处理中` / `Doing` / `开发中` | 当前活跃处理的任务 |
| 验证中 | `测试中` / `验收中` / `In Review` | 等待验证或评审 |
| 已完成 | `已完成` / `Done` / `已关闭` | 已结束的任务 |
| 已取消 | `已取消` / `作废` / ` wontfix` | 不再进行的任务 |
| 挂起/阻塞 | `pending` / `阻塞中` / `等待中` | 因外部依赖暂停的任务 |

## 注意事项

- `organizationId` 为可选，省略时自动注入用户默认组织
- `perPage` 默认较小，查询大量任务时建议显式设置（如 50）
- `get_my_project_workitems` 返回按 category 分组的数据，`search_workitems` 返回扁平列表
