package db

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/Sirupsen/logrus"
)

func hashIndex(ty *Type, index string, val interface{}) (string, error) {
	// return fmt.Sprintf("%s:%s:%s", ty.Name, index, val), nil

	hash := fnv.New64()
	_, err := fmt.Fprintf(hash, "%s:%s:%s", ty.Name, index, val)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(hash.Sum64(), 16), nil
}

func (w *WaifuDB) GetIndexPointer(ty *Type, index string, val interface{}) (string, error) {
	key, err := hashIndex(ty, index, val)
	if err != nil {
		return "", err
	}

	b, err := w.store.Get(bktIndexes, key)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (w *WaifuDB) PutIndexEntries(ty *Type, data map[string]interface{}) {
	for k, v := range data {
		if ty.HasIndex(k) {
			key, err := hashIndex(ty, k, v)
			if err != nil {
				w.logger.WithError(err).WithFields(logrus.Fields{
					"data": data,
					"type": ty,
				}).Error("hashing index failed")
				continue
			}
			err = w.store.Set(bktIndexes, key, []byte(data["id"].(string)))
			if err != nil {
				w.logger.WithError(err).WithFields(logrus.Fields{
					"data": data,
					"type": ty,
				}).Error("setting index failed")
				continue
			}
		}
	}
}
