Run make test
go test -race ./...
panic: test timed out after 10m0s
	running tests:
		TestStopPredicatesAcrossRunStartAndGateRunCommands (33s)
		TestStopPredicatesAcrossRunStartAndGateRunCommands/run-start/resolved (7s)

goroutine 19607 [running]:
testing.(*M).startAlarm.func1()
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2802 +0x605
created by time.goFunc
	/opt/hostedtoolcache/go/1.26.2/x64/src/time/sleep.go:215 +0x45

goroutine 1 [chan receive]:
testing.(*T).Run(0xc000222488, {0x149c41d, 0x32}, 0x14d1560)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2109 +0xb3e
testing.runTests.func1(0xc000222488)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2585 +0x86
testing.tRunner(0xc000222488, 0xc0003338f0)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2036 +0x21d
testing.runTests({0x147986c, 0x18}, {0x1484835, 0x20}, 0xc000016150, {0x1aee4c0, 0x11f, 0x11f}, {0xc2a1b179a038b95c, 0x8ce70a4396, ...})
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2583 +0x9ea
testing.(*M).Run(0xc0001190e0)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2443 +0xf4c
github.com/mbell-nzl/fid/cmd/fid.runCommandTestSuite({0x14daca0, 0xc0001190e0})
	/home/runner/work/fid/fid/cmd/fid/workspace_fixture_test.go:86 +0x3a9
github.com/mbell-nzl/fid/cmd/fid.TestMain(0xc0001190e0)
	/home/runner/work/fid/fid/cmd/fid/workspace_fixture_test.go:64 +0x2e
main.main()
	_testmain.go:620 +0x172

goroutine 4006 [chan receive, 7 minutes]:
testing.(*T).Parallel(0xc0003be6c8)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:1803 +0x50c
github.com/mbell-nzl/fid/cmd/fid.TestClaimAcquireBindingMatrixRejectsWithoutMutation(0xc0003be6c8)
	/home/runner/work/fid/fid/cmd/fid/generic_claim_test.go:68 +0x3f
testing.tRunner(0xc0003be6c8, 0x14d0e60)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2036 +0x21d
created by testing.(*T).Run in goroutine 1
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2101 +0xb13

goroutine 3989 [chan receive, 8 minutes]:
testing.(*T).Parallel(0xc00034a908)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:1803 +0x50c
github.com/mbell-nzl/fid/cmd/fid.TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild(0xc00034a908)
bufio.(*Reader).writeBuf(0xc000435f00, {0x7fe079dc9000, 0xc0002112c0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000435f00, {0x7fe079dc9000, 0xc0002112c0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac900)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24676 [select]:
io.(*pipe).write(0xc000210f00, {0xc0005c6000, 0xb5c, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc000210f00, {0xc0005c6000, 0xb5c, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc000484f00, {0x7fe079dc9000, 0xc000210f00})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000484f00, {0x7fe079dc9000, 0xc000210f00})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac540)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24693 [select]:
io.(*pipe).write(0xc0002100c0, {0xc000662000, 0x3185, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc0002100c0, {0xc000662000, 0x3185, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc000480f00, {0x7fe079dc9000, 0xc0002100c0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000480f00, {0x7fe079dc9000, 0xc0002100c0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac000)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24680 [select]:
io.(*pipe).write(0xc000211200, {0xc000614000, 0x25c, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc000211200, {0xc000614000, 0x25c, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc000432f00, {0x7fe079dc9000, 0xc000211200})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000432f00, {0x7fe079dc9000, 0xc000211200})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac840)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24678 [select]:
io.(*pipe).write(0xc000211080, {0xc000578000, 0x1ac9, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc000211080, {0xc000578000, 0x1ac9, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc0000f9f00, {0x7fe079dc9000, 0xc000211080})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc0000f9f00, {0x7fe079dc9000, 0xc000211080})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac6c0)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24677 [select]:
io.(*pipe).write(0xc000210fc0, {0xc0005e0000, 0xf30, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc000210fc0, {0xc0005e0000, 0xf30, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc000482f00, {0x7fe079dc9000, 0xc000210fc0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000482f00, {0x7fe079dc9000, 0xc000210fc0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0000ac600)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24707 [chan receive]:
database/sql.(*Tx).awaitDone(0xc0000c2480)
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:2212 +0x49
created by database/sql.(*DB).beginDC in goroutine 24666
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:1925 +0x41f

goroutine 24666 [runnable]:
modernc.org/libc.Xfree(0x4cca31?, 0x7fe030713aa0)
	/home/runner/go/pkg/mod/modernc.org/libc@v1.72.3/mem_musl.go:100 +0x91
modernc.org/sqlite/lib._sqlite3MemFree(0xc0000a6420, 0x7fe030713aa8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:12911 +0x38
modernc.org/sqlite/lib.Xsqlite3_free(0xc0000a6420, 0x7fe030713aa8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:14138 +0x3e2
modernc.org/sqlite/lib._sqlite3DbNNFreeNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe030713aa8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:14206 +0x27a
modernc.org/sqlite/lib._sqlite3ExprDeleteNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe030713aa8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79839 +0x545
modernc.org/sqlite/lib._sqlite3ExprDelete(...)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79845
modernc.org/sqlite/lib._exprListDeleteNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe079f8e428)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:80793 +0xe5
modernc.org/sqlite/lib._sqlite3ExprListDelete(...)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:80812
modernc.org/sqlite/lib._sqlite3ExprDeleteNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe030712ba8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79820 +0x26d
modernc.org/sqlite/lib._sqlite3ExprDeleteNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe0307125a8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79815 +0x154
modernc.org/sqlite/lib._sqlite3ExprDeleteNN(0xc0000a6420, 0x7fe079f8ec28, 0x7fe0307126a8)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79815 +0x154
modernc.org/sqlite/lib._sqlite3ExprDelete(...)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:79845
modernc.org/sqlite/lib._sqlite3BeginTrigger(0xc0000a6420, 0x7fe079f8b420, 0x7fe0303a10c8, 0x7fe0303a10e0, 0x21, 0x80, 0x0, 0x7fe030715e28, 0x7fe0307126a8, 0x0, ...)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:125027 +0x1648
modernc.org/sqlite/lib._yy_reduce(0xc0000a6420, 0x7fe0303a1028, 0x105, 0x0?, {0xdb?, 0x153d260?}, 0x7fe079f8b420)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:157300 +0x15145
modernc.org/sqlite/lib._sqlite3Parser(0xc0000a6420, 0x7fe0303a1028, 0x5, {0x4cb72c?, 0x4cce55?})
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:157790 +0x1fd
modernc.org/sqlite/lib._sqlite3RunParser(0xc0000a6420, 0x7fe079f8b420, 0x7fe030322514)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:160472 +0xfaa
modernc.org/sqlite/lib._sqlite3Prepare(0xc0000a6420, 0x7fe079f8ec28, 0x7fe030322514, 0xffffffff, 0x80, 0x0, 0x7fe07a11fa10, 0x7fe07a11fd20)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:115375 +0xb92
modernc.org/sqlite/lib._sqlite3LockAndPrepare(0xc0000a6420, 0x7fe079f8ec28, 0x7fe030322514, 0xffffffff, 0x80, 0x0, 0x7fe07a11fa10, 0x7fe07a11fd20)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:115434 +0x2b1
modernc.org/sqlite/lib.Xsqlite3_prepare_v2(...)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/lib/sqlite_linux_amd64.go:115521
modernc.org/sqlite.(*conn).prepareV2(0xc0002109c0, 0xc0001d47e0)
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/conn.go:628 +0x2d9
modernc.org/sqlite.(*stmt).exec(0xc00049c168, {0x10bb2a0, 0x15942c0}, {0x15942c0, 0x0, 0x0})
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/stmt.go:191 +0x56f
modernc.org/sqlite.(*conn).exec(0xc0002109c0, {0x10bb2a0, 0x15942c0}, {0xc0001ff000, 0xe37}, {0x15942c0, 0x0, 0x0})
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/conn.go:951 +0x17a
modernc.org/sqlite.(*conn).ExecContext(0xc0002109c0, {0x10bb2a0, 0x15942c0}, {0xc0001ff000, 0xe37}, {0x15942c0, 0x0, 0x0})
	/home/runner/go/pkg/mod/modernc.org/sqlite@v1.51.0/conn.go:1208 +0x95
database/sql.ctxDriverExec({0x10bb2a0, 0x15942c0}, {0x7fe07a1e8ea0, 0xc0002109c0}, {0x0, 0x0}, {0xc0001ff000, 0xe37}, {0x15942c0, 0x0, ...})
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/ctxutil.go:31 +0xeb
database/sql.(*DB).execDC.func2()
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:1713 +0x266
database/sql.withLock({0x10baab8, 0xc0000c2c00}, 0xc0001d4e90)
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:3572 +0xa3
database/sql.(*DB).execDC(0x50c528?, {0x10bb2a0, 0x15942c0}, 0xc0000c2c00, 0xc000306620, {0xc0001ff000, 0xe37}, {0x0, 0x0, 0x0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:1708 +0x405
database/sql.(*Tx).ExecContext(0xc0000c2480, {0x10bb2a0, 0x15942c0}, {0xc0001ff000, 0xe37}, {0x0, 0x0, 0x0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:2516 +0x129
database/sql.(*Tx).Exec(...)
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:2525
github.com/golang-migrate/migrate/v4/database/sqlite.(*Sqlite).executeQuery(0xc00049d518, {0xc0001ff000, 0xe37})
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/database/sqlite/sqlite.go:207 +0xb1
github.com/golang-migrate/migrate/v4/database/sqlite.(*Sqlite).Run(0xc00049d518, {0x10b99c0, 0xc000210de0})
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/database/sqlite/sqlite.go:199 +0xdb
github.com/golang-migrate/migrate/v4.(*Migrate).runMigrations(0xc0000c2c80, 0xc0000710a0)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:744 +0x30d
github.com/golang-migrate/migrate/v4.(*Migrate).Up(0xc0000c2c80)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:282 +0x16c
github.com/mbell-nzl/fid/internal/bootstrap.initializeStagedWorkspace.MigrateUp.func1(0xc0000c2c80)
	/home/runner/work/fid/fid/internal/db/migrate.go:15 +0x2f
github.com/mbell-nzl/fid/internal/db.runMigration({0xc000024a10, 0x6e}, {0xc0001db6c0, 0x10}, 0x10b31c0)
	/home/runner/work/fid/fid/internal/db/migrate.go:54 +0x3de
github.com/mbell-nzl/fid/internal/db.MigrateUp(...)
	/home/runner/work/fid/fid/internal/db/migrate.go:14
github.com/mbell-nzl/fid/internal/bootstrap.initializeStagedWorkspace({_, _}, {_, _}, {{0xc000375720, 0x47}, {{0xc000375810, 0x4c}, {0xc00020c360, 0x54}, ...}, ...}, ...)
	/home/runner/work/fid/fid/internal/bootstrap/init.go:175 +0xa8
github.com/mbell-nzl/fid/internal/bootstrap.initializeWorkspaceWithHooks({{0xc000375720, 0x47}, {0xc0001db6c0, 0x10}, {0x1080613, 0x3}}, {0x0?})
	/home/runner/work/fid/fid/internal/bootstrap/init.go:98 +0x607
github.com/mbell-nzl/fid/internal/bootstrap.InitializeWorkspace(...)
	/home/runner/work/fid/fid/internal/bootstrap/init.go:71
github.com/mbell-nzl/fid/internal/export.TestImportAfterPublicationFaultReplaysInstalledResult(0xc0004998c8)
	/home/runner/work/fid/fid/internal/export/import_test.go:108 +0x114
testing.tRunner(0xc0004998c8, 0x10b24d8)
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2036 +0x21d
created by testing.(*T).Run in goroutine 1
	/opt/hostedtoolcache/go/1.26.2/x64/src/testing/testing.go:2101 +0xb13

goroutine 24697 [select]:
io.(*pipe).write(0xc00007c840, {0xc000696000, 0x14ad, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:86 +0x2e5
io.(*PipeWriter).Write(0xc00007c840, {0xc000696000, 0x14ad, 0x186a0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/io/pipe.go:161 +0x47
bufio.(*Reader).writeBuf(0xc000483f00, {0x7fe079dc9000, 0xc00007c840})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:563 +0xc8
bufio.(*Reader).WriteTo(0xc000483f00, {0x7fe079dc9000, 0xc00007c840})
	/opt/hostedtoolcache/go/1.26.2/x64/src/bufio/bufio.go:520 +0xdf
github.com/golang-migrate/migrate/v4.(*Migration).Buffer(0xc0001fa0c0)
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migration.go:156 +0x591
github.com/golang-migrate/migrate/v4.(*Migrate).readUp.func2()
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:618 +0x38
created by github.com/golang-migrate/migrate/v4.(*Migrate).readUp in goroutine 24670
	/home/runner/go/pkg/mod/github.com/golang-migrate/migrate/v4@v4.19.1/migrate.go:617 +0x49d

goroutine 24667 [select]:
database/sql.(*DB).connectionOpener(0xc0001be000, {0x10bb3f8, 0xc00045e5f0})
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:1261 +0xed
created by database/sql.OpenDB in goroutine 24666
	/opt/hostedtoolcache/go/1.26.2/x64/src/database/sql/sql.go:841 +0x287
FAIL	github.com/mbell-nzl/fid/internal/export	600.033s
ok  	github.com/mbell-nzl/fid/internal/issues	401.093s
ok  	github.com/mbell-nzl/fid/internal/journal	382.237s
ok  	github.com/mbell-nzl/fid/internal/labels	82.794s
ok  	github.com/mbell-nzl/fid/internal/links	117.636s
ok  	github.com/mbell-nzl/fid/internal/maintenance	47.904s
?   	github.com/mbell-nzl/fid/internal/policy	[no test files]
ok  	github.com/mbell-nzl/fid/internal/project	2.213s
ok  	github.com/mbell-nzl/fid/internal/requirements	35.798s
ok  	github.com/mbell-nzl/fid/internal/sqlquery	9.558s
ok  	github.com/mbell-nzl/fid/internal/stats	16.059s
ok  	github.com/mbell-nzl/fid/internal/validation	1.011s
ok  	github.com/mbell-nzl/fid/internal/workflow	602.082s
FAIL
make: *** [Makefile:10: test] Error 1
Error: Process completed with exit code 2.