package savemigrate

import "testing"

func TestMigrateSameVersionKeepsData(t *testing.T) {
	migrator := New()
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 1)
	if !ok || got["gold"] != 100 {
		t.Fatalf("同版本迁移应该原样返回，得到 %v ok=%v", got, ok)
	}
}

func TestMigrateUnknownVersionKeepsData(t *testing.T) {
	got, _ := New().Migrate(map[string]int{"gold": 5}, 1, 2)
	if got["gold"] != 5 {
		t.Fatalf("没有注册步骤时数据不该变，得到 %v", got)
	}
}

func TestMigrateDoesNotTouchInput(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	original := map[string]int{"gold": 100}
	migrator.Migrate(original, 1, 2)
	if original["gold"] != 100 {
		t.Fatalf("迁移不该改动传入存档，得到 %d", original["gold"])
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestRegisterRejectsNonAdjacent(t *testing.T) {
	migrator := New()
	if migrator.Register(1, 3, map[string]int{"gold": 10}) {
		t.Fatal("跳版本的登记应该返回 false")
	}
	if migrator.Register(2, 2, map[string]int{"gold": 10}) {
		t.Fatal("原地不动的登记应该返回 false")
	}
	if migrator.Register(3, 2, map[string]int{"gold": 10}) {
		t.Fatal("倒退的登记应该返回 false")
	}
}

func TestRegisterRejectsDuplicateFrom(t *testing.T) {
	migrator := New()
	if !migrator.Register(1, 2, map[string]int{"gold": 10}) {
		t.Fatal("首次登记应该成功")
	}
	if migrator.Register(1, 2, map[string]int{"gold": 20}) {
		t.Fatal("同一个 from 重复登记应该返回 false")
	}
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 2)
	if !ok || got["gold"] != 110 {
		t.Fatalf("重复登记不该覆盖首次步骤，得到 %v ok=%v", got, ok)
	}
}

func TestMigrateAppliesEachStep(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"gold": 5, "gem": 1})
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
	if !ok || got["gold"] != 115 || got["gem"] != 1 {
		t.Fatalf("逐版本迁移应该并入每步新增字段，得到 %v ok=%v", got, ok)
	}
}

func TestMigrateMissingStepFails(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
	if ok {
		t.Fatal("区间缺步骤应该返回 false")
	}
	if got["gold"] != 100 {
		t.Fatalf("缺步骤时不该做任何合并，得到 %v", got)
	}
}

func TestMigrateBackwardFails(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"gold": 5})
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 3, 2)
	if ok {
		t.Fatal("目标版本低于来源应该返回 false")
	}
	if got["gold"] != 100 {
		t.Fatalf("回退失败时应该返回原样副本，得到 %v", got)
	}
}

func TestMigrateIsDeterministic(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"gold": 5})
	first, ok1 := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
	second, ok2 := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
	if !ok1 || !ok2 || first["gold"] != second["gold"] {
		t.Fatalf("同一输入重复迁移结果应该相同，得到 %v 与 %v", first, second)
	}
}

func TestScannedGrowsWithStepsNotRegistry(t *testing.T) {
	migrator := New()
	for index := 0; index < 3000; index++ {
		migrator.Register(index, index+1, map[string]int{"gold": 1})
	}
	for index := 0; index < 3000; index++ {
		migrator.Migrate(map[string]int{"gold": 0}, 0, 0)
	}
	if migrator.Scanned() > 6000 {
		t.Fatalf("同版本迁移不该扫描注册表，scanned=%d", migrator.Scanned())
	}
}
