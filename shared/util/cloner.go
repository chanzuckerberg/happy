package util

import (
	"encoding/json"
	"reflect"
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

func DeepMerge(dst, src map[string]any) {
	for k, v := range src {
		t := reflect.TypeOf(v)
		if t == nil {
			continue
		}
		switch v.(type) {
		case map[string]any:
			subSrc := v.(map[string]any)
			subDst := dst[k].(map[string]any)
			DeepMerge(subDst, subSrc)
			dst[k] = subDst
		default:
			if _, ok := dst[k]; ok {
				if v != nil {
					dst[k] = v
				}
			}
		}
	}
}
