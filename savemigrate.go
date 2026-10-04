// Package savemigrate 是存档版本迁移：按注册的迁移步骤把旧存档逐版本升到目标版本。
// 缺陷：同版本也会重跑步骤、缺步骤时静默放过、目标版本比来源还低也照过、重复登记照收、
// 每次迁移全扫注册表。
package savemigrate

type step struct {
	from int
	to   int
	add  map[string]int
}

// Migrator 是迁移注册表。
type Migrator struct {
	steps   []step
	scanned int
}

// New 建空注册表。
func New() *Migrator { return &Migrator{} }

// Register 注册一步迁移，返回是否成功。缺陷：不校验相邻版本、不查重。
func (m *Migrator) Register(from, to int, add map[string]int) bool {
	m.steps = append(m.steps, step{from: from, to: to, add: add})
	return true
}

// Migrate 把存档从 from 版本迁到 to 版本，返回迁移后的数据与是否成功。
// 缺陷：只要步骤的起点落在区间里就并入、同版本也重跑、缺步骤不报错、全扫注册表。
func (m *Migrator) Migrate(data map[string]int, from, to int) (map[string]int, bool) {
	out := map[string]int{}
	for key, value := range data {
		out[key] = value
	}
	m.scanned += len(m.steps)
	for _, item := range m.steps {
		if item.from >= from && item.from <= to {
			for key, value := range item.add {
				out[key] += value
			}
		}
	}
	return out, true
}

// Scanned 是累计扫描的步骤数（规模观测）。
func (m *Migrator) Scanned() int { return m.scanned }
