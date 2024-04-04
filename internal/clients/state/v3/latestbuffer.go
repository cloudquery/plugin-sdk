package state

import (
	"strconv"
	"time"
)

type LatestBuffer struct {
	data map[string][]formatMap
}

type format int

const (
	fOriginal format = iota
	fInt
	fFloat
	fTime
)

type formatMap map[format]any

func NewLatestBuffer() *LatestBuffer {
	return &LatestBuffer{
		data: make(map[string][]formatMap),
	}
}

func (l *LatestBuffer) Add(key, value string) {
	l.data[key] = append(l.data[key], possibleFormats(value))
}

func (l *LatestBuffer) All() map[string]string {
	ret := make(map[string]string, len(l.data))
	for key := range l.data {
		ret[key] = l.Get(key)
	}
	return ret
}

func (l *LatestBuffer) Get(key string) string {
	vals := l.data[key]
	switch len(vals) {
	case 0: // unknown key
		return ""
	case 1: // single value
		return vals[0][fOriginal].(string)
	}

	var common *format
	for _, fmts := range vals {
		f := fmts.Type()
		if f == nil {
			break
		}
		if common == nil {
			common = f
			continue
		}
		if (*f == fInt && *common == fFloat) || (*f == fFloat && *common == fInt) {
			*common = fFloat
			continue
		}
		if *common != *f {
			common = nil
			break
		}
	}
	if common == nil || *common == fOriginal {
		return vals[0][fOriginal].(string) // no known common format, return first value
	}

	valIndex := -1

	// Depending on our common format, find the largest value
	switch *common {
	case fInt:
		var maxVal int64
		for i, fVal := range vals {
			v, ok := fVal.Get(fInt)
			if !ok {
				panic("wanted to get fInt but not found")
			}
			vt := v.(int64)
			if valIndex == -1 || vt > maxVal {
				maxVal = vt
				valIndex = i
			}
		}
	case fFloat:
		var maxVal float64
		for i, fVal := range vals {
			v, ok := fVal.Get(fFloat)
			if !ok {
				v, ok = fVal.Get(fInt)
				if ok {
					v = float64(v.(int64)) // cast int to float as other values are floats
				}
			}
			if !ok {
				panic("wanted to get fFloat but not found")
			}
			vt := v.(float64)
			if valIndex == -1 || vt > maxVal {
				maxVal = vt
				valIndex = i
			}
		}
	case fTime:
		var maxVal time.Time
		for i, fVal := range vals {
			v, ok := fVal.Get(fTime)
			if !ok {
				panic("wanted to get fTime but not found")
			}
			vt := v.(time.Time)
			if valIndex == -1 || vt.After(maxVal) {
				maxVal = vt
				valIndex = i
			}
		}
	}

	if valIndex == -1 {
		// should not happen
		panic("valIndex is -1, common is " + strconv.FormatInt(int64(*common), 10))
		//return vals[0][fOriginal].(string) // return first value
	}

	return vals[valIndex][fOriginal].(string) // Return original value (of the largest value)
}

func possibleFormats(s string) formatMap {
	fmts := map[format]any{
		fOriginal: s,
	}

	// ints also parse as floats, so check int and if it passes ignore float
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		fmts[fInt] = v
	} else if v, err := strconv.ParseFloat(s, 64); err == nil {
		fmts[fFloat] = v
	}
	if v, err := time.Parse(time.RFC3339, s); err == nil {
		fmts[fTime] = v
	}
	if v, err := time.Parse(time.RFC3339Nano, s); err == nil {
		fmts[fTime] = v
	}

	return fmts
}

func (m formatMap) String() string {
	return m[fOriginal].(string)
}

func (m formatMap) Get(f format) (any, bool) {
	if v, ok := m[f]; ok {
		return v, true
	}
	return nil, false
}

func (m formatMap) Type() *format {
	formats := make([]format, 0, len(m))
	for f := range m {
		formats = append(formats, f)
	}
	switch len(formats) {
	case 0: // empty
		return nil
	case 1: // this would be fOriginal
		return &formats[0]
	case 2: // fOriginal and one other
		if formats[0] == fOriginal {
			return &formats[1]
		}
		return &formats[0]
	}

	return nil // no common format
}
