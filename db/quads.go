package db

import (
	"fmt"
	"strconv"
)

// Quad format:
//
// subject_id.predicate.hash(object_id) => object_id
//

func hashObject(ptr string) (string, error) {
	// return fmt.Sprintf("%s:%s:%s", ty.Name, index, val), nil

	hLock.Lock()
	defer hLock.Unlock()
	h.Reset()
	_, err := fmt.Fprint(h, ptr)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(h.Sum64(), 16), nil
}

func (w *WaifuDB) SiftQuads(ty *Type, data map[string]interface{}) (quads map[string]string, err error) {
	quads = map[string]string{}

	id := fmt.Sprintf("%s:%s", ty.Name, data["id"])
	idObj, err := hashObject(id)
	if err != nil {
		return quads, err
	}

	for k, v := range ty.Relations {
		objPtr, ok := data[k].(string)
		if ok {
			obj, err := hashObject(objPtr)
			if err != nil {
				return quads, err
			}
			qk := fmt.Sprintf("%s.%s.%s", id, k, obj)
			quads[qk] = objPtr

			if v != "" {
				iq := fmt.Sprintf("%s.%s.%s", objPtr, v, idObj)
				quads[iq] = id
			}
			delete(data, k)
		}

		if v != "" {
			objPtr, ok = data[v].(string)
			if ok {
				obj, err := hashObject(objPtr)
				if err != nil {
					return quads, err
				}
				qk := fmt.Sprintf("%s.%s.%s", id, v, obj)
				quads[qk] = objPtr

				iq := fmt.Sprintf("%s.%s.%s", objPtr, k, idObj)
				quads[iq] = id
				delete(data, v)
			}
		}
	}

	return quads, err
}

func (w *WaifuDB) SaveQuads(quads map[string]string) {
	for k, v := range quads {
		w.store.Set(bktQuads, k, []byte(v))
	}
}
