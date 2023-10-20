package schema

import "encoding/json"

type Tags []string

func (t *Tags) Len() int {
	if t == nil {
		return 0
	}
	return len(*t)
}

func (t *Tags) Add(tag string) {
	if t == nil {
		*t = Tags{tag}
		return
	}
	if t.Contains(tag) {
		return
	}
	*t = append(*t, tag)
}

func (t *Tags) Remove(tag string) {
	if t == nil {
		return
	}
	for i, v := range *t {
		if v == tag {
			*t = append((*t)[:i], (*t)[i+1:]...)
			break
		}
	}
}

func (t *Tags) Contains(tag string) bool {
	if t == nil {
		return false
	}
	for _, v := range *t {
		if v == tag {
			return true
		}
	}
	return false
}

func (t *Tags) MarshalJSON() ([]byte, error) {
	if t == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal([]string(*t))
}

func (t *Tags) UnmarshalJSON(data []byte) error {
	if t == nil {
		return nil
	}
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil {
		return err
	}
	*t = tags
	return nil
}
