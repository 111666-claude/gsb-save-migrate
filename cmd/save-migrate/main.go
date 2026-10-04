// Command save-migrate 跑存档版本迁移样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/savemigrate"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("save-migrate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "missing", "missing / same / back / dup / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "missing":
		migrator := savemigrate.New()
		migrator.Register(1, 2, map[string]int{"gold": 10})
		got, ok := migrator.Migrate(map[string]int{"gold": 100}, 1, 3)
		fmt.Fprintf(stdout, "ok=%v gold=%d\n", ok, got["gold"])
	case "same":
		migrator := savemigrate.New()
		migrator.Register(1, 2, map[string]int{"gold": 10})
		got, _ := migrator.Migrate(map[string]int{"gold": 100}, 1, 1)
		fmt.Fprintf(stdout, "gold=%d\n", got["gold"])
	case "back":
		migrator := savemigrate.New()
		migrator.Register(1, 2, map[string]int{"gold": 10})
		migrator.Register(2, 3, map[string]int{"gold": 5})
		_, ok := migrator.Migrate(map[string]int{"gold": 100}, 3, 2)
		fmt.Fprintf(stdout, "ok=%v\n", ok)
	case "dup":
		migrator := savemigrate.New()
		migrator.Register(1, 2, map[string]int{"gold": 10})
		second := migrator.Register(1, 2, map[string]int{"gold": 20})
		fmt.Fprintf(stdout, "second=%v\n", second)
	case "scan":
		migrator := savemigrate.New()
		for index := 0; index < 3000; index++ {
			migrator.Register(index, index+1, map[string]int{"gold": 1})
		}
		for index := 0; index < 3000; index++ {
			migrator.Migrate(map[string]int{"gold": 0}, 0, 0)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", migrator.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample missing|same|back|dup|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
