package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEquals(t *testing.T) {
	r := require.New(t)

	m1 := map[string]any{}
	m2 := map[string]any{}

	r.True(DeepEquals(m1, m2))

	m1 = map[string]any{
		"foo": "bar",
	}
	m2 = map[string]any{
		"foo": "bar",
	}
	r.True(DeepEquals(m1, m2))

	m1 = map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
			},
		},
	}
	m2 = map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
			},
		},
	}
	r.True(DeepEquals(m1, m2))
}

func TestDeepClone(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1.0,
			"quz": map[string]any{
				"quy": true,
			},
		},
	}
	var m2 map[string]any
	err := DeepClone(&m2, m1)
	r.NoError(err)
	r.True(DeepEquals(m1, m2))
}

func TestDeepMerge(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
			},
		},
	}
	m2 := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
				"c":   2,
			},
		},
	}

	res := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
				"c":   2,
			},
		},
	}

	err := DeepMerge(m2, m1)

	r.NoError(err)
	r.True(DeepEquals(res, m2))
}

func TestDeepMerge2(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
			},
		},
	}
	m2 := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
				"c":   2,
				"a":   1,
			},
		},
	}

	err := DeepMerge(m2, m1)

	r.Error(err)
}

func TestDeepIntersect(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
			},
		},
	}
	m2 := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
				"c":   2,
			},
		},
	}

	res := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
			},
		},
	}

	m3 := DeepIntersect(m2, m1)
	r.True(DeepEquals(res, m3))
}

func TestDeepDiff(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
			},
		},
	}
	m2 := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
				"c":   2,
			},
		},
	}

	res := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"quz": map[string]any{
				"quy": true,
			},
		},
	}

	res1 := map[string]any{
		"r": true,
		"baz": map[string]any{
			"e": "f",
			"quz": map[string]any{
				"a": "b",
			},
		},
	}

	res2 := map[string]any{
		"baz": map[string]any{
			"quz": map[string]any{
				"c": 2,
			},
		},
	}

	intersection := DeepIntersect(m2, m1)
	r.True(DeepEquals(res, intersection))

	m3 := DeepDiff(intersection, m1)
	r.True(DeepEquals(res1, m3))

	m3 = DeepDiff(intersection, m2)
	r.True(DeepEquals(res2, m3))
}

func TestDeepCleanup(t *testing.T) {
	r := require.New(t)
	m1 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
			},
		},
	}
	m2 := map[string]any{
		"foo": "bar",
		"r":   true,
		"baz": map[string]any{
			"qux": "quux",
			"quu": 1,
			"e":   "f",
			"d":   map[string]any{},
			"quz": map[string]any{
				"quy": true,
				"a":   "b",
				"c":   nil,
			},
		},
	}
	m3 := DeepCleanup(m1)
	r.True(DeepEquals(m1, m3))
	m3 = DeepCleanup(m2)
	r.True(DeepEquals(m1, m3))
}
