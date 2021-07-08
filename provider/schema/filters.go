package schema

func DeleteParentFieldsFilter(fields ...string) func(meta ClientMeta, parent *Resource) []interface{} {
	return func(meta ClientMeta, parent *Resource) []interface{} {
		if parent == nil {
			return nil
		}
		var filters []interface{}
		for _, f := range fields {
			filters = append(filters, f, parent.Get(f))
		}
		return filters
	}
}
