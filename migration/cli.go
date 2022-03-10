package migration

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cloudquery/cq-provider-sdk/database"
	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
)

const defaultPath = "./resources/provider/migrations"

// Run is the main entry point for CLI usage.
func Run(ctx context.Context, p *provider.Provider, outputPath string) error {
	const defaultPrefix = "unreleased"

	if outputPath == "" {
		outputPath = defaultPath
	}

	outputPathParam := flag.String("path", outputPath, "Path to migrations directory")
	prefixParam := flag.String("prefix", defaultPrefix, "Prefix for files")
	doFullParam := flag.Bool("full", false, "Generate initial migrations (prefix will be 'init')")
	dialectParam := flag.String("dialect", "", "Dialect to generate initial migrations (empty: all)")
	dsnParam := flag.String("dsn", os.Getenv("CQ_DSN"), "DSN to compare changes against in upgrade mode")
	schemaName := flag.String("schema", "public", "Schema to compare tables from in upgrade mode")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		return fmt.Errorf("more args than necessary")
	}

	if *doFullParam && *prefixParam == defaultPrefix {
		*prefixParam = "init"
	}

	if *prefixParam != "" {
		// Add the first "." in <prefix>.up.sql, only if we have a prefix
		*prefixParam += "."
	}

	if *doFullParam {
		dialects, err := parseInputDialect(dialectParam)
		if err != nil {
			return err
		}

		if err := GenerateFull(ctx, hclog.L(), p, dialects, *outputPathParam, *prefixParam); err != nil {
			return fmt.Errorf("failed to generate migrations: %w", err)
		}
		return nil
	}

	if *dsnParam == "" {
		return fmt.Errorf("DSN not specified: Use -dsn or set CQ_DSN")
	}

	pool, dialectType, err := connect(ctx, *dsnParam)
	if err != nil {
		return err
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if *dialectType == schema.TSDB && *schemaName == "public" {
		*schemaName = "history"
	}

	if err := GenerateDiff(ctx, hclog.L(), conn, *schemaName, *dialectType, p, *outputPathParam, *prefixParam); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	return nil
}

func connect(ctx context.Context, dsn string) (*pgxpool.Pool, *schema.DialectType, error) {
	detectedDialect, newDSN, err := database.ParseDialectDSN(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse dsn: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(newDSN)
	if err != nil {
		return nil, nil, err
	}
	poolCfg.LazyConnect = true
	pool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	return pool, &detectedDialect, err
}

func parseInputDialect(inputDialect *string) ([]schema.DialectType, error) {
	defaultDialectsFullMode := []schema.DialectType{
		schema.Postgres,
		schema.TSDB,
	}

	var dialects []schema.DialectType
	if *inputDialect == "" {
		dialects = defaultDialectsFullMode
	} else {
		for _, d := range defaultDialectsFullMode {
			if string(d) == *inputDialect {
				dialects = append(dialects, d)
				break
			}
		}
		if len(dialects) == 0 {
			return nil, fmt.Errorf("invalid dialect %q", *inputDialect)
		}
	}

	return dialects, nil
}
