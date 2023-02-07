package postgresql

const defaultTableName = "cloudquery_state"

type Spec struct {
	ConnectionString string   `json:"connection_string,omitempty"`
	TableName        string   `json:"table_name,omitempty"`
	PgxLogLevel      LogLevel `json:"pgx_log_level,omitempty"`
}

func (s *Spec) SetDefaults() {
	if s.TableName == "" {
		s.TableName = defaultTableName
	}
}
