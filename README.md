# AnnularQueue
环形定时队列（目前只有单机）

环形队列,一个圆圈划分为60段,每段拥有一个任务列表包含运行到当前分段时刻需要执行的任务,每秒向后循环一段环形队列
1 定时执行
2 失败重试（默认再次经过定时的时间后重试）
3 失败重试次数设置
4 单分段内队列执行
