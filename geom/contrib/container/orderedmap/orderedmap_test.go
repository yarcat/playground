package orderedmap

import "testing"

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Errorf("New() = nil, wanted non-nil")
	}
}

func TestSet(t *testing.T) {
	for _, test := range []struct {
		name string
		p    pairs
		ks   []string
		want pairs
	}{
		{
			name: "empty",
			ks:   []string{"a"},
			want: nil,
		},
		{
			name: "single pair",
			p:    pairs{{"a", "b"}},
			ks:   []string{"a", "b"},
			want: pairs{{"a", "b"}},
		},
		{
			name: "multiple pairs",
			p:    pairs{{"a", "b"}, {"b", "c"}},
			ks:   []string{"a", "b", "c"},
			want: pairs{{"a", "b"}, {"b", "c"}},
		},
		{
			name: "overwrite key",
			p:    pairs{{"a", "b"}, {"a", "c"}},
			ks:   []string{"a"},
			want: pairs{{"a", "c"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := New()
			set(m, test.p)
			actual := getinternal(m, test.ks)
			if !test.want.Equal(actual) {
				t.Errorf("Set(%v); Get(%v) = %v, want %v",
					test.p, test.ks, actual, test.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	for _, test := range []struct {
		name string
		p    pairs
		ks   []string
		want pairs
	}{
		{
			name: "empty",
			ks:   []string{"a"},
			want: nil,
		},
		{
			name: "existing",
			p:    pairs{{"a", "b"}, {"b", "c"}},
			ks:   []string{"a", "b", "c"},
			want: pairs{{"a", "b"}, {"b", "c"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := New()
			setinternal(m, test.p)
			actual := get(m, test.ks)
			if !test.want.Equal(actual) {
				t.Errorf("set(%v); Get = %v, want %v", test.p, actual, test.want)
			}
		})
	}
}

func TestIter(t *testing.T) {
	for _, test := range []struct {
		name string
		p    pairs
		want pairs
	}{
		{
			name: "empty",
		},
		{
			name: "keep order",
			p:    pairs{{"a", "b"}, {"b", "c"}},
			want: pairs{{"a", "b"}, {"b", "c"}},
		},
		{
			name: "keep original order",
			p:    pairs{{"a", "b"}, {"b", "c"}, {"a", "d"}},
			want: pairs{{"a", "d"}, {"b", "c"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := New()
			set(m, test.p)
			actual := iter(m)
			if !test.want.Equal(actual) {
				t.Errorf("Set(%v); Iter = %v, want %v", test.p, actual, test.want)
			}
		})
	}
}

type keyvalue [2]string

func (kv keyvalue) Equal(other keyvalue) bool {
	return kv[0] == other[0] && kv[1] == other[1]
}

func (kv keyvalue) Key() string {
	return kv[0]
}

func (kv keyvalue) Value() string {
	return kv[1]
}

type pairs []keyvalue

func (p pairs) Equal(other pairs) bool {
	if len(p) != len(other) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if !p[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

func setinternal(m *OrderedMap, p pairs) {
	for _, kv := range p {
		item := &mapItem{key: kv.Key(), value: kv.Value()}
		e := m.items.PushBack(item)
		m.dict[kv.Key()] = e
	}
}

func set(m *OrderedMap, p pairs) {
	for _, kv := range p {
		m.Set(kv.Key(), kv.Value())
	}
}

func getinternal(m *OrderedMap, keys []string) (p pairs) {
	for _, k := range keys {
		e, ok := m.dict[k]
		if !ok {
			continue
		}
		sv := asMapItem(e).value.(string)
		p = append(p, keyvalue{k, sv})
	}
	return
}

func get(m *OrderedMap, keys []string) (p pairs) {
	for _, k := range keys {
		v, ok := m.Get(k)
		if !ok {
			continue
		}
		sv := v.(string)
		p = append(p, keyvalue{k, sv})
	}
	return
}

func iter(m *OrderedMap) (p pairs) {
	for it := m.Iter(); it.Next(); {
		k, v := it.Value()
		p = append(p, keyvalue{k.(string), v.(string)})
	}
	return
}
