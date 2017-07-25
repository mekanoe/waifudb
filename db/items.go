package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/kayteh/waifudb/datastore"
	"github.com/segmentio/ksuid"
)

func (w *WaifuDB) PutItem(t string, data map[string]interface{}) (dat map[string]interface{}, err error) {
	ty, err := w.GetType(t)
	if err != nil {
		return dat, err
	}

	kid, err := ksuid.NewRandom()
	if err != nil {
		return dat, err
	}
	id := kid.String()

	key := fmt.Sprintf("%s:%s", t, id)
	data["id"] = id
	data["type"] = t

	err = w.store.SetJSON(bktData, key, &data)
	if err != nil {
		return dat, err
	}

	// TODO: make goroutinable
	w.PutIndexEntries(ty, data)

	dat = data
	return dat, err
}

func (w *WaifuDB) GetItem(t, id string) (dat map[string]interface{}, err error) {
	key := fmt.Sprintf("%s:%s", t, id)

	err = w.store.GetJSON(bktData, key, &dat)
	if err != nil {
		return dat, err
	}

	return dat, err
}

func (w *WaifuDB) GetItemByKey(t, key string, val interface{}) (map[string]interface{}, error) {
	var dat map[string]interface{}

	ty, err := w.GetType(t)
	if err != nil {
		return dat, err
	}

	if ty.HasIndex(key) {

		// get by index
		p, err := w.GetIndexPointer(ty, key, val)
		if err != nil && err != datastore.ErrNotFound {
			return dat, err
		}

		if err == datastore.ErrNotFound {
			return w.search(t, key, val)
		}

		return w.GetItem(t, p)
	}

	return w.search(t, key, val)
}

func (w *WaifuDB) search(t, key string, val interface{}) (dat map[string]interface{}, err error) {
	prefix := []byte(t)

	err = w.store.Bolt.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bktData).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var tmp map[string]interface{}
			err := json.Unmarshal(v, &tmp)
			if err != nil {
				return err
			}

			if tmp[key] == val {
				dat = tmp
				break
			}
		}

		return nil
	})

	return dat, err
}
