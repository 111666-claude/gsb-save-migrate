// Package savemigrate 是存档版本迁移：按注册的迁移步骤把旧存档逐版本升到目标版本。
package savemigrate

type step struct {
	from int
	to   int
	add  map[string]int
}

// Migrator 是迁移注册表。
type Migrator struct {
	steps   map[int]step
	scanned int
}

// New 建空注册表。
func New() *Migrator { return &Migrator{steps: map[int]step{}} }

// Register 注册一步迁移，返回是否成功。
// 只接受 from 升到 from+1 的相邻步骤；同一 from 重复登记返回 false 且不生效。
func (m *Migrator) Register(from, to int, add map[string]int) bool {
	if to != from+1 {
		return false
	}
	if _, exists := m.steps[from]; exists {
		return false
	}
	m.steps[from] = step{from: from, to: to, add: add}
	return true
}

// Migrate 把存档从 from 版本迁到 to 版本，返回迁移后的数据与是否成功。
// 逐版本并入每步的新增字段；from == to 时原样返回副本；to < from 或区间缺步时
// 返回原样副本与 false，不做任何合并。
func (m *Migrator) Migrate(data map[string]int, from, to int) (map[string]int, bool) {
	out := cloneData(data)
	if from == to {
		return out, true
	}
	if to < from {
		return out, false
	}
	for version := from; version < to; version++ {
		if _, ok := m.steps[version]; !ok {
			return out, false
		}
	}
	for version := from; version < to; version++ {
		item := m.steps[version]
		m.scanned++
		for key, value := range item.add {
			out[key] += value
		}
	}
	return out, true
}

func cloneData(data map[string]int) map[string]int {
	out := make(map[string]int, len(data))
	for key, value := range data {
		out[key] = value
	}
	return out
}

// Scanned 是累计执行的迁移步骤数（规模观测）。
func (m *Migrator) Scanned() int { return m.scanned }
