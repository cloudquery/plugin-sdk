package schema

// DeleteParentIdFilter is mostly used for table relations to delete table data based on parent's cq_id
func DeleteParentIdFilter(id string) func(meta ClientMeta, parent *Resource) []interface{} {
	return func(meta ClientMeta, parent *Resource) []interface{} {
		if parent == nil {
			return nil
		}
		return []interface{}{id, parent.Id()}
	}
}
