package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"song-library/internal/config"
	"song-library/internal/constants"
	"song-library/internal/utils"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	directionUp    = "up"
	directionDown  = "down"
	migrationsPath = "migrations"
	queriesPath    = "internal/repository/queries/migrations"

	// SQL файлы
	sqlCreateMigrationsTable = "create_migrations_table.sql"
	sqlCheckMigrationExists  = "check_migration_exists.sql"
	sqlGetAppliedMigrations  = "get_applied_migrations.sql"
	sqlInsertMigration       = "insert_migration.sql"
	sqlDeleteMigration       = "delete_migration.sql"

	actionApply    = "применения"
	actionRollback = "отката"
	resultApplied  = "применены"
	resultRolled   = "отменены"
)

type Migrator struct {
	db     *sql.DB
	logger *log.Logger
}

type Migration struct {
	Version  string
	FileName string
	Content  []byte
}

func NewMigrator(ctx context.Context, dbConfig config.DatabaseConfig, logger *log.Logger) (*Migrator, error) {
	if logger == nil {
		return nil, fmt.Errorf(constants.ErrLoggerNil)
	}

	logger.Printf(constants.LogDBConnecting, dbConfig.Host, dbConfig.Port)

	db, err := sql.Open(constants.PostgresDriver, fmt.Sprintf(
		constants.PostgresConnectionString,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.SSLMode,
	))
	if err != nil {
		logger.Printf(constants.LogError, constants.ErrDBConnection, err.Error())
		return nil, fmt.Errorf(constants.ErrDBConnection, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf(constants.ErrDBConnection, err)
	}

	logger.Println(constants.LogDBConnected)
	return &Migrator{
		db:     db,
		logger: logger,
	}, nil
}

// getMigrationFiles читает и сортирует файлы миграций из директории
// suffixWord определяет постфикс файла миграции: "up" - мигарция, "down" - откат
func (m *Migrator) getMigrationFiles(suffixWord string) ([]string, error) {
	projectRoot := utils.GetProjectRoot(0)
	migrationsDir := filepath.Join(projectRoot, migrationsPath)
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrReadingMigrationDir, err)
	}

	suffix := fmt.Sprintf(constants.SQLSuffix, suffixWord)
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			fullPath := filepath.Join(migrationsDir, file.Name())
			migrationFiles = append(migrationFiles, fullPath)
		}
	}

	// Сортировка файлов миграции по имени для гарантированного порядка выполнения
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// executeMigrations выполняет миграции в указанном направлении
// direction - направление миграции ("up" или "down")
func (m *Migrator) executeMigrations(ctx context.Context, direction string) error {
	if ctx == nil {
		return fmt.Errorf(constants.ErrContextNil)
	}

	startTime := time.Now()

	actionNames := map[string]string{
		directionUp:   actionApply,
		directionDown: actionRollback,
	}
	actionResults := map[string]string{
		directionUp:   resultApplied,
		directionDown: resultRolled,
	}

	actionName, ok := actionNames[direction]
	if !ok {
		return fmt.Errorf(constants.ErrMigrationDiraction, direction)
	}

	m.logger.Printf(constants.LogMigrationStart, actionName)

	if err := m.ensureMigrationsTable(ctx); err != nil {
		return fmt.Errorf(constants.ErrMigrationTableCheck, err)
	}

	files, err := m.getMigrationFiles(direction)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		m.logger.Println(constants.LogMigrationNotFound)
		return nil
	}

	for _, file := range files {
		version := extractVersionFromFileName(file)

		if direction == directionUp {
			exists, err := m.isMigrationApplied(ctx, version)
			if err != nil {
				return err
			}
			if exists {
				m.logger.Printf(constants.LogMigrationSkipped, version)
				continue
			}
		}

		err = m.executeInTransaction(ctx, func(tx *sql.Tx) error {
			if err := m.executeSingleMigration(ctx, tx, file); err != nil {
				return err
			}

			var queryFile string
			if direction == directionUp {
				queryFile = sqlInsertMigration
			} else {
				queryFile = sqlDeleteMigration
			}
			query, err := utils.ReadQueryFile(queriesPath, queryFile)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, query, version)
			return err
		})

		if err != nil {
			return err
		}

		m.logger.Printf(constants.LogMigrationProcess, cases.Title(language.Russian).String(actionName), file)
	}

	duration := time.Since(startTime)
	m.logger.Printf(constants.LogMigrationPerTime, actionResults[direction], duration)
	return nil
}

func (m *Migrator) Up(ctx context.Context) error {
	return m.executeMigrations(ctx, directionUp)
}

func (m *Migrator) Down(ctx context.Context) error {
	return m.executeMigrations(ctx, directionDown)
}

func (m *Migrator) Close() error {
	return m.db.Close()
}

func (m *Migrator) ensureMigrationsTable(ctx context.Context) error {
	query, err := utils.ReadQueryFile(queriesPath, sqlCreateMigrationsTable)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, query)
	return err
}

func (m *Migrator) executeInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(constants.ErrTransactionStart, err)
	}
	defer tx.Rollback() // откатится только если не было commit

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(constants.ErrTransactionCommit, err)
	}
	return nil
}

func (m *Migrator) isMigrationApplied(ctx context.Context, version string) (bool, error) {
	query, err := utils.ReadQueryFile(queriesPath, sqlCheckMigrationExists)
	if err != nil {
		return false, err
	}
	var exists bool
	err = m.db.QueryRowContext(ctx, query, version).Scan(&exists)
	return exists, err
}

func (m *Migrator) GetAppliedMigrations(ctx context.Context) ([]string, error) {
	query, err := utils.ReadQueryFile(queriesPath, sqlGetAppliedMigrations)
	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

// extractVersionFromFileName извлекает версию из имени файла миграции
// Пример: "20240315_create_users_up.sql" -> "20240315"
func extractVersionFromFileName(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// executeSingleMigration выполняет одну миграцию
// file - путь к файлу миграции
func (m *Migrator) executeSingleMigration(ctx context.Context, tx *sql.Tx, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf(constants.ErrReadingMigrationFile, file, err)
	}

	_, err = tx.ExecContext(ctx, string(content))
	if err != nil {
		return fmt.Errorf(constants.ErrExecutingMigration, file, err)
	}

	return nil
}
