# go-schedule

go-schedule

## db 模拟 shedlock-spring 实现作业调度。

```mysql
CREATE TABLE shedlock
(
    name       VARCHAR(64),
    lock_until TIMESTAMP(3) NULL,
    locked_at  TIMESTAMP(3) NULL,
    locked_by  VARCHAR(255),
    PRIMARY KEY (name)
)
```