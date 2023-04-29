package schema

func TestSourceTable(name string) *Table {
	return &Table{
		Name:        name,
		Description: "Test table",
		Columns: ColumnList{
			CqIDColumn,
			CqParentIDColumn,
			{
				Name:            "uuid_pk",
				Type:            TypeUUID,
				CreationOptions: ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name:            "string_pk",
				Type:            TypeString,
				CreationOptions: ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name: "bool",
				Type: TypeBool,
			},
			{
				Name: "int",
				Type: TypeInt,
			},
			{
				Name: "float",
				Type: TypeFloat,
			},
			{
				Name: "uuid",
				Type: TypeUUID,
			},
			{
				Name: "text",
				Type: TypeString,
			},
			{
				Name: "text_with_null",
				Type: TypeString,
			},
			{
				Name: "bytea",
				Type: TypeByteArray,
			},
			{
				Name: "text_array",
				Type: TypeStringArray,
			},
			{
				Name: "text_array_with_null",
				Type: TypeStringArray,
			},
			{
				Name: "int_array",
				Type: TypeIntArray,
			},
			{
				Name: "timestamp",
				Type: TypeTimestamp,
			},
			{
				Name: "json",
				Type: TypeJSON,
			},
			{
				Name: "uuid_array",
				Type: TypeUUIDArray,
			},
			{
				Name: "inet",
				Type: TypeInet,
			},
			{
				Name: "inet_array",
				Type: TypeInetArray,
			},
			{
				Name: "cidr",
				Type: TypeCIDR,
			},
			{
				Name: "cidr_array",
				Type: TypeCIDRArray,
			},
			{
				Name: "macaddr",
				Type: TypeMacAddr,
			},
			{
				Name: "macaddr_array",
				Type: TypeMacAddrArray,
			},
		},
	}
}

func TestTable(name string) *Table {
	sourceTable := TestSourceTable(name)
	sourceTable.Columns = append(ColumnList{
		CqSourceNameColumn,
		CqSyncTimeColumn,
	}, sourceTable.Columns...)
	return sourceTable
}