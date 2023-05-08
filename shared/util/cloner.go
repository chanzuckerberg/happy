package util

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

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

func DeepMerge(dst, src map[string]any) error {
	for k, v := range src {
		t := reflect.TypeOf(v)
		if t == nil {
			continue
		}
		switch v.(type) {
		case map[string]any:
			subSrc := v.(map[string]any)
			subDst := dst[k].(map[string]any)
			err := DeepMerge(subDst, subSrc)
			if err != nil {
				return err
			}
			dst[k] = subDst
		default:
			if _, ok := dst[k]; ok {
				if v != nil {
					if o, ok := dst[k]; ok {
						t1 := reflect.TypeOf(o)
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
		}
	}
	return nil
}
