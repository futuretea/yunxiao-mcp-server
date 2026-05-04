# Yunxiao MCP Risk Tools 快速参考

## 预测性工具（本 Skill 核心）

| 工具 | 用途 | 关键参数 | 返回数据 |
|------|------|---------|---------|
| `get_sprint_velocity` | 历史迭代完成率趋势 | `projectId`, `categories`, `sprintCount` | 每迭代：总任务数、完成数、完成率 |
| `get_workitem_status_timeline` | 单个任务状态变更时间线 | `workitemId` | 状态变更记录、停留时长、操作人 |
| `get_blocker_analysis` | 依赖阻塞全景 | `projectId`, `categories`, `sampleLimit` | 被阻塞任务、阻塞他人任务、分类统计 |
| `get_member_workload_trend` | 成员负载分布与趋势 | `projectId`, `assigneeIds`, `daysBack` | 每成员：任务数、状态分布、逾期数、近期活跃度 |

## 辅助工具（配合预测）

| 工具 | 用途 | 关键参数 |
|------|------|---------|
| `get_project_risk_dashboard` | 逾期/高优先级/停滞任务 | `projectId`, `categories` |
| `get_project_member_task_status` | 成员详细任务状态 | `projectId`, `assigneeIds` |
| `get_project_overview` | 项目快照（迭代、版本、成员） | `projectId` |
| `get_sprint_overview` | 单个迭代任务分布 | `projectId`, `sprintId` |
| `get_project_workitem_detail` | 单个任务详情（含活动日志） | `workitemId` |

## 常用调用顺序

**标准风险评估**：
```
1. get_project_overview -> 获取项目上下文、迭代、版本
2. get_sprint_velocity -> 历史完成率趋势
3. get_blocker_analysis -> 阻塞全景
4. get_member_workload_trend -> 容量分布
5. (可选) get_project_risk_dashboard -> 补充逾期/停滞数据
```

**深入诊断阻塞任务**：
```
1. get_blocker_analysis -> 识别关键阻塞者
2. get_workitem_status_timeline -> 分析阻塞任务历史
3. get_project_workitem_detail -> 获取任务详情和评论
```

**交付预测专项**：
```
1. get_sprint_velocity -> 历史速率
2. get_sprint_overview -> 当前迭代任务分布
3. list_versions -> 版本截止日期
4. 计算：剩余任务 / 历史速率 = 预测完成时间
```

## 注意事项

- `organizationId` 为可选，省略时自动注入用户默认组织
- `get_sprint_velocity` 默认分析最近 5 个已完成/归档迭代，可通过 `sprintCount` 调整（最大 20）
- `get_blocker_analysis` 需要逐个查询任务依赖关系，耗时与 `sampleLimit` 成正比
- `get_member_workload_trend` 默认分析最近 30 天活跃度，可通过 `daysBack` 调整
- 所有预测性工具均为只读，不会修改任何任务数据
