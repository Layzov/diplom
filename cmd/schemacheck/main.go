// One-off: inspect DB schema before/after AutoMigrate (go run ./cmd/schemacheck).
package main

import (
	"fmt"
	"os"

	"diplom/internal/config"
	"diplom/internal/db"

	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()

	gdb, err := db.Open(cfg.DB.DatabaseDSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}

	if err := db.Ping(gdb); err != nil {
		fmt.Fprintf(os.Stderr, "ping: %v\n", err)
		fmt.Fprintln(os.Stderr, "Postgres недоступен (Docker/служба не запущены?).")
		os.Exit(2)
	}

	fmt.Println("=== BEFORE AutoMigrate ===")
	dump(gdb)

	fmt.Println("\n--- AutoMigrate ---")
	if err := db.Migrate(gdb); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(3)
	}

	fmt.Println("\n=== AFTER AutoMigrate ===")
	dump(gdb)
}

func dump(gdb *gorm.DB) {
	type col struct {
		Table  string
		Column string
		Type   string
	}
	want := []col{
		{"subjects", "user_id", ""},
		{"topics", "subject_id", ""},
		{"tasks", "topic_id", ""},
		{"attachments", "task_id", ""},
		{"task_attempts", "session_id", ""},
		{"attempt_stats", "session_id", ""},
		{"repetitions", "repeat_at", ""},
	}

	for _, w := range want {
		var exists bool
		q := `
SELECT EXISTS (
  SELECT 1 FROM information_schema.columns
  WHERE table_schema = 'public' AND table_name = ? AND column_name = ?
)`
		if err := gdb.Raw(q, w.Table, w.Column).Scan(&exists).Error; err != nil {
			fmt.Printf("  %s.%s: error %v\n", w.Table, w.Column, err)
			continue
		}
		fmt.Printf("  %s.%s exists = %v\n", w.Table, w.Column, exists)
	}

	type fk struct {
		Table      string
		Column     string
		RefTable   string
		RefColumn  string
		Constraint string
	}
	var fks []fk
	err := gdb.Raw(`
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS ref_table,
    ccu.column_name AS ref_column,
    tc.constraint_name
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu
  ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage ccu
  ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = 'public'
  AND tc.table_name IN ('subjects','topics','tasks','attachments','sessions','task_attempts','repetitions')
ORDER BY 1, 2`).Scan(&fks).Error
	if err != nil {
		fmt.Printf("  FK query error: %v\n", err)
		return
	}
	fmt.Printf("  FK count on core tables: %d\n", len(fks))
	for _, f := range fks {
		fmt.Printf("    %s.%s -> %s.%s (%s)\n", f.Table, f.Column, f.RefTable, f.RefColumn, f.Constraint)
	}

	var viewOK bool
	_ = gdb.Raw(`
SELECT EXISTS (
  SELECT 1 FROM information_schema.views
  WHERE table_schema = 'public' AND table_name = 'calendar'
)`).Scan(&viewOK)
	fmt.Printf("  view calendar exists = %v\n", viewOK)

	type idx struct {
		Table string
		Name  string
	}
	var indexes []idx
	err = gdb.Raw(`
SELECT tablename, indexname
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename IN (
    'subjects','topics','tasks','sessions',
    'task_attempts','repetitions','attempt_stats','users'
  )
ORDER BY tablename, indexname`).Scan(&indexes).Error
	if err != nil {
		fmt.Printf("  indexes query error: %v\n", err)
		return
	}
	fmt.Printf("  indexes on core tables: %d\n", len(indexes))
	wantIdx := []string{
		"idx_repetitions_planned_user_due",
		"idx_sessions_user_started_at",
		"idx_task_attempts_user_task",
	}
	for _, w := range wantIdx {
		found := false
		for _, ix := range indexes {
			if ix.Name == w {
				found = true
				break
			}
		}
		fmt.Printf("    %s = %v\n", w, found)
	}
}
