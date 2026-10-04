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
