# save-migrate

存档版本迁移：按注册的迁移步骤把旧存档逐版本升到目标版本，只用 Go 标准库。

```
go run ./cmd/save-migrate --sample=missing
go run ./cmd/save-migrate --sample=same
go run ./cmd/save-migrate --sample=back
go run ./cmd/save-migrate --sample=dup
go run ./cmd/save-migrate --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **登记**：每一步只能从 `from` 升到 `from+1`，否则返回 false 且不登记；同一个 `from` 重复登记也返回 false。
- **逐版本**：`Migrate` 从 `from` 一步步走到 `to`，每步并入该步的新增字段；`from` 等于 `to` 时原样返回副本。
- **缺步骤**：区间里缺任何一步都返回 `(原样副本, false)`，不做任何合并；`to` 小于 `from` 也返回 false。
- **不变量**：迁移不改动传入的存档；同一输入重复迁移结果相同。
- **代价**：迁移不许全扫注册表，`Scanned` 不随步骤数乘迁移次数增长；1 万个版本、每秒 10 万次迁移。

## 输出契约（不改格式）

```
ok=.. gold=..
gold=..
ok=..
second=..
scanned=..
```

## 目录

```
savemigrate.go        版本链、缺步骤报错与幂等登记
savemigrate_test.go   单元测试
cmd/save-migrate/     命令行入口
```
