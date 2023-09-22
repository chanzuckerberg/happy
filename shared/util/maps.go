package util

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// Create a deep copy of src into dst
func DeepClone(dst, src any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &dst); err != nil {
		return err
	}
	return nil
}

// Given two nested maps, merge the src into the dst
func DeepMerge(dst, src map[string]any) error {
	for k, v := range src {
		t := reflect.TypeOf(v)
		if t == nil {
			continue
		}
		if subSrc, ok := v.(map[string]any); ok {
			if dst[k] == nil {
				dst[k] = make(map[string]any)
			}
			subDst := dst[k].(map[string]any)
			err := DeepMerge(subDst, subSrc)
			if err != nil {
				return err
			}
			dst[k] = subDst
			continue
		}

		if v != nil {
			if o, ok := dst[k]; ok {
				t1 := reflect.TypeOf(o)
				// Source and destination types must be the same.
				if t1 != t {
					dstTypeName := "unknown"
					if t1 != nil {
						dstTypeName = t1.Name()
					}
					return errors.Errorf("Destination type (%s) is not the same as source type (%s) for field '%s'.", dstTypeName, t.Name(), k)
				}
			}
			dst[k] = v
		}
	}
	return nil
}

// Based on two nested maps, calculate the intersection of the two.
func DeepIntersect(m1, m2 map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}

		m1, ok1 := v1.(map[string]any)
		m2, ok2 := v2.(map[string]any)
		if ok1 && ok2 {
			res[k] = DeepIntersect(m1, m2)
			continue
		}

		if reflect.DeepEqual(v1, v2) {
			res[k] = v2
		}
	}

	return res
}

// Determine if two nested maps are equal.
func DeepEquals(m1, m2 map[string]any) bool {
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}

		m1, ok1 := v1.(map[string]any)
		m2, ok2 := v2.(map[string]any)
		if ok1 && ok2 {
			if !DeepEquals(m1, m2) {
				return false
			}
			continue
		}

		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}

	return true
}

// Build a nested map based on the passed argument, ignoring null values and empty maps.
func DeepCleanup(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		if v != nil {
			if v1, ok := v.(map[string]any); ok {
				res1 := DeepCleanup(v1)
				if len(res1) > 0 {
					res[k] = res1
				}
				continue
			}
			res[k] = v
		}
	}
	return res
}

// Calculate a difference between the maps and present it as a nested map.
func DeepDiff(base, overlay map[string]any) map[string]any {
	if len(base) == 0 && len(overlay) == 0 {
		return nil
	}
	res := make(map[string]any)
	for k, v1 := range base {
		v2, ok := overlay[k]
		if !ok {
			continue
		}

		m1, ok1 := v1.(map[string]any)
		m2, ok2 := v2.(map[string]any)
		if ok1 && ok2 {
			res[k] = DeepDiff(m1, m2)
			continue
		}

		if !reflect.DeepEqual(v1, v2) {
			res[k] = v2
		}
	}

	for k, v1 := range overlay {
		v2, ok := base[k]
		if !ok {
			res[k] = v1
			continue
		}

		_, ok1 := v1.(map[string]any)
		_, ok2 := v2.(map[string]any)
		if ok1 && ok2 {
			continue
		}

		if !reflect.DeepEqual(v1, v2) {
			res[k] = v2
		}
	}

	return res
}
