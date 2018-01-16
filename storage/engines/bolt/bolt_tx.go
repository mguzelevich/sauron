package bolt

import (
	"github.com/boltdb/bolt"
)

type TxContext struct {
	tx *bolt.Tx
}

func NewTx(tx *bolt.Tx) *TxContext {
	return &TxContext{
		tx: tx,
	}
}

func (t *TxContext) name() {

}

func (t *TxContext) getBucket(names []string) (*bolt.Bucket, error) {
	var bucket *bolt.Bucket
	for _, bn := range names {
		if bucket == nil {
			bucket = t.tx.Bucket([]byte(bn))
		} else {
			bucket = bucket.Bucket([]byte(bn))
		}
	}
	return bucket, nil
}

func (t *TxContext) getOrCreateBucket(names []string) (*bolt.Bucket, error) {
	var bucket *bolt.Bucket
	var err error
	for _, bn := range names {
		if err != nil {
			break
		}
		if bn == "" {
			bn = "<?>"
		}
		if bucket == nil {
			bucket, err = t.tx.CreateBucketIfNotExists([]byte(bn))
		} else {
			bucket, err = bucket.CreateBucketIfNotExists([]byte(bn))
		}
	}
	return bucket, err
}

func (t *TxContext) get(names []string, key string) (string, error) {
	bucket, _ := t.getBucket(names)
	response := string(bucket.Get([]byte(key)))
	return response, nil
}

func (t *TxContext) put(names []string, key string, value string) error {
	bucket, _ := t.getBucket(names)
	err := bucket.Put([]byte(key), []byte(value))
	return err
}
