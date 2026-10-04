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
		t.Fatal("跨版本登记应该返回 false")
	}
	if migrator.Register(1, 1, map[string]int{"gold": 10}) {
		t.Fatal("原地登记应该返回 false")
	}
	if _, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 2); ok {
		t.Fatal("被拒绝的步骤不该登记")
	}
}

func TestRegisterRejectsDuplicateFrom(t *testing.T) {
	migrator := New()
	if !migrator.Register(1, 2, map[string]int{"gold": 10}) {
		t.Fatal("首次登记应该成功")
	}
	if migrator.Register(1, 2, map[string]int{"gold": 20}) {
		t.Fatal("同一 from 重复登记应该返回 false")
	}
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 2)
	if !ok || got["gold"] != 110 {
		t.Fatalf("重复登记不该生效，得到 %v ok=%v", got, ok)
	}
}

func TestMigrateAppliesEachStep(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"silver": 5})
	got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
	if !ok || got["gold"] != 110 || got["silver"] != 5 {
		t.Fatalf("应该逐版本并入新增字段，得到 %v ok=%v", got, ok)
	}
}

func TestMigrateMissingStepReturnsCopyAndFalse(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	original := map[string]int{"gold": 100}
	got, ok := migrator.Migrate(original, 1, 3)
	if ok {
		t.Fatal("缺步骤时应该返回 false")
	}
	if got["gold"] != 100 {
		t.Fatalf("缺步骤时不该做任何合并，得到 %v", got)
	}
	got["gold"] = 999
	if original["gold"] != 100 {
		t.Fatalf("失败时也该返回副本，原存档被改成 %d", original["gold"])
	}
}

func TestMigrateBackwardReturnsFalse(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"gold": 5})
	original := map[string]int{"gold": 100}
	got, ok := migrator.Migrate(original, 3, 2)
	if ok {
		t.Fatal("目标版本低于来源应该返回 false")
	}
	if got["gold"] != 100 {
		t.Fatalf("回退时数据不该变，得到 %v", got)
	}
}

func TestMigrateRepeatedIsDeterministic(t *testing.T) {
	migrator := New()
	migrator.Register(1, 2, map[string]int{"gold": 10})
	migrator.Register(2, 3, map[string]int{"gold": 5})
	input := map[string]int{"gold": 100}
	first, ok1 := migrator.Migrate(input, 1, 3)
	second, ok2 := migrator.Migrate(input, 1, 3)
	if !ok1 || !ok2 || first["gold"] != second["gold"] || first["gold"] != 115 {
		t.Fatalf("重复迁移结果应该一致，得到 %v 与 %v", first, second)
	}
	if input["gold"] != 100 {
		t.Fatalf("重复迁移不该改动传入存档，得到 %d", input["gold"])
	}
}

func TestScannedCountsAppliedStepsOnly(t *testing.T) {
	migrator := New()
	for index := 0; index < 100; index++ {
		migrator.Register(index, index+1, map[string]int{"gold": 1})
	}
	for index := 0; index < 50; index++ {
		migrator.Migrate(map[string]int{"gold": 0}, 0, 0)
	}
	if scanned := migrator.Scanned(); scanned != 0 {
		t.Fatalf("同版本迁移不该扫描步骤，得到 %d", scanned)
	}
	migrator.Migrate(map[string]int{"gold": 0}, 0, 10)
	if scanned := migrator.Scanned(); scanned != 10 {
		t.Fatalf("应该只扫描走过的 10 步，得到 %d", scanned)
	}
	migrator.Migrate(map[string]int{"gold": 0}, 0, 50)
	if scanned := migrator.Scanned(); scanned != 60 {
		t.Fatalf("缺步失败不该计入扫描，累计应为 60，得到 %d", scanned)
	}
}
