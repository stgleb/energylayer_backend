package factory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageFactory(t *testing.T) {
	_, err := StorageFactory("wrong")
	assert.NotNil(t, err)

	db, err := StorageFactory("memory")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, err)
}
