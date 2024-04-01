package dbdef

type TableDefinition struct {
	TableName string
	HashKey   string
	RangeKey  string
	TtlKey    string

	NewEntity func() interface{}
}

func (t *TableDefinition) NewTableName(prefix string) string {
	return prefix + t.TableName
}
