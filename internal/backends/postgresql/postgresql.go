package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	pgx_zero_log "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type Backend struct {
	sourceName          string
	spec                Spec
	currentDatabaseName string
	table               schema.Table
	conn                *pgxpool.Pool
	logger              zerolog.Logger
}

func New(ctx context.Context, logger zerolog.Logger, sourceSpec specs.Source) (*Backend, error) {
	b := &Backend{
		sourceName: sourceSpec.Name,
		logger:     logger.With().Str("module", "pg-backend").Logger(),
	}
	spec := Spec{}
	err := sourceSpec.UnmarshalBackendSpec(&spec)
	if err != nil {
		return nil, err
	}
	spec.SetDefaults()
	b.spec = spec
	b.table = schema.Table{
		Name: spec.TableName,
		Columns: []schema.Column{
			{Name: "source_name", Type: schema.TypeString, CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "table_name", Type: schema.TypeString, CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "client_id", Type: schema.TypeString, CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "value", Type: schema.TypeString},
		},
	}
	logLevel, err := tracelog.LogLevelFromString(spec.PgxLogLevel.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx log level %s: %w", spec.PgxLogLevel, err)
	}
	b.logger.Info().Str("pgx_log_level", spec.PgxLogLevel.String()).Msg("Initializing postgresql backend")
	pgxConfig, err := pgxpool.ParseConfig(spec.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string %w", err)
	}
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return nil
	}

	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgx_zero_log.NewLogger(b.logger),
		LogLevel: logLevel,
	}
	pgxConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"
	b.conn, err = pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	b.currentDatabaseName, err = b.currentDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current database: %w", err)
	}
	err = b.createTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return b, nil
}

func (b *Backend) createTable(ctx context.Context) error {
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(pgx.Identifier{b.table.Name}.Sanitize())
	sb.WriteString(" (")
	for _, col := range b.table.Columns {
		sb.WriteString(pgx.Identifier{col.Name}.Sanitize() + " TEXT,")
	}
	sb.WriteString("PRIMARY KEY (")
	pks := b.table.PrimaryKeys()
	for i, col := range pks {
		sb.WriteString(pgx.Identifier{col}.Sanitize())
		if i < len(pks)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")
	sb.WriteString(")")
	_, err := b.conn.Exec(ctx, sb.String())
	return err
}

// Set sets the value for the given table and client id.
func (b *Backend) Set(ctx context.Context, table, clientID, value string) error {
	var sb strings.Builder
	sb.WriteString("insert into ")
	sb.WriteString(pgx.Identifier{b.table.Name}.Sanitize())
	sb.WriteString(" (")
	columnsLen := len(b.table.Columns)
	for i, c := range b.table.Columns {
		sb.WriteString(pgx.Identifier{c.Name}.Sanitize())
		if i < columnsLen-1 {
			sb.WriteString(",")
		} else {
			sb.WriteString(") values (")
		}
	}
	for i := range b.table.Columns {
		sb.WriteString(fmt.Sprintf("$%d", i+1))
		if i < columnsLen-1 {
			sb.WriteString(",")
		} else {
			sb.WriteString(")")
		}
	}
	_, err := b.conn.Exec(ctx, sb.String(), b.sourceName, table, clientID, value)
	return err
}

// Get returns the value for the given table and client id.
func (b *Backend) Get(ctx context.Context, table, clientID string) (string, error) {
	var sb strings.Builder
	sb.WriteString("select value from ")
	sb.WriteString(pgx.Identifier{b.table.Name}.Sanitize())
	sb.WriteString(" where source_name = $1 and table_name = $2 and client_id = $3")
	var value string
	err := b.conn.QueryRow(ctx, sb.String(), b.sourceName, table, clientID).Scan(&value)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// Close closes the backend.
func (b *Backend) Close(_ context.Context) error {
	var err error
	if b.conn == nil {
		return fmt.Errorf("backend already closed or not initialized")
	}
	if b.conn != nil {
		b.conn.Close()
		b.conn = nil
	}
	return err
}

func (b *Backend) currentDatabase(ctx context.Context) (string, error) {
	var db string
	err := b.conn.QueryRow(ctx, "select current_database()").Scan(&db)
	if err != nil {
		return "", err
	}
	return db, nil
}
