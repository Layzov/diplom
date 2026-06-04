// One-off: EXPLAIN hot queries (go run ./cmd/explain [-user=<uuid>] [-analyze]).
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"diplom/internal/config"
	"diplom/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	analyze := flag.Bool("analyze", false, "use EXPLAIN ANALYZE (executes queries; run after seed data)")
	userFlag := flag.String("user", "", "user UUID for filtered queries (default: first user in DB)")
	flag.Parse()

	cfg := config.MustLoad()
	gdb, err := db.Open(cfg.DB.DatabaseDSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	if err := db.Ping(gdb); err != nil {
		fmt.Fprintf(os.Stderr, "ping: %v\n", err)
		os.Exit(2)
	}
	if err := db.Migrate(gdb); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(3)
	}

	userID, err := resolveUserID(gdb, *userFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "user: %v\n", err)
		os.Exit(4)
	}
	fmt.Printf("user_id = %s\n", userID)
	if *analyze {
		fmt.Println("mode: EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)")
	} else {
		fmt.Println("mode: EXPLAIN (FORMAT TEXT) — add -analyze after seed for timing")
	}
	fmt.Println()

	now := time.Now().UTC()
	limit := 20

	queries := []struct {
		name string
		sql  string
		args []any
	}{
		{
			name: "calendar upcoming (stats repo)",
			sql: `SELECT * FROM calendar
			      WHERE user_id = ? AND due_at >= ?
			      ORDER BY due_at ASC LIMIT ?`,
			args: []any{userID, now, limit},
		},
		{
			name: "sessions by user",
			sql: `SELECT * FROM sessions
			      WHERE user_id = ?
			      ORDER BY started_at DESC`,
			args: []any{userID},
		},
		{
			name: "task_attempts by session (latest session)",
			sql: `SELECT ta.* FROM task_attempts ta
			      WHERE ta.session_id = (
			        SELECT s.id FROM sessions s
			        WHERE s.user_id = ?
			        ORDER BY s.started_at DESC LIMIT 1
			      )
			      ORDER BY ta.created_at ASC`,
			args: []any{userID},
		},
		{
			name: "repetitions list planned",
			sql: `SELECT * FROM repetitions
			      WHERE user_id = ? AND status = 'planned'
			      ORDER BY repeat_at ASC`,
			args: []any{userID},
		},
		{
			name: "topic_learning_stats",
			sql: `SELECT * FROM topic_learning_stats WHERE user_id = ?`,
			args: []any{userID},
		},
		{
			name: "user_learning_stats",
			sql: `SELECT * FROM user_learning_stats WHERE user_id = ?`,
			args: []any{userID},
		},
		{
			name: "session_learning_stats",
			sql: `SELECT * FROM session_learning_stats WHERE user_id = ? ORDER BY started_at DESC`,
			args: []any{userID},
		},
	}

	for _, q := range queries {
		fmt.Printf("=== %s ===\n", q.name)
		plan, err := runExplain(gdb, *analyze, q.sql, q.args...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error: %v\n\n", err)
			continue
		}
		fmt.Println(plan)
		fmt.Println()
	}
}

func resolveUserID(gdb *gorm.DB, flagValue string) (uuid.UUID, error) {
	if flagValue != "" {
		return uuid.Parse(flagValue)
	}
	var id uuid.UUID
	err := gdb.Table("users").Select("id").Order("created_at ASC").Limit(1).Scan(&id).Error
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("no users in database (register or run seed first)")
	}
	return id, nil
}

func runExplain(gdb *gorm.DB, analyze bool, sql string, args ...any) (string, error) {
	prefix := "EXPLAIN (FORMAT TEXT)"
	if analyze {
		prefix = "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)"
	}
	sql = strings.TrimSpace(sql)
	rows, err := gdb.Raw(prefix+" "+sql, args...).Rows()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var parts []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		parts = append(parts, line)
	}
	return strings.Join(parts, "\n"), rows.Err()
}
