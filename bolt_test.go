package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetupDB (t *testing.T) {
	var fn = "test.db"
	t.Run("NormalCase", func(t *testing.T) {
		db, err := setupDB(fn)
		assert.NoError(t, err)
		assert.Equal(t, db.Path(), "test.db")
	})

	t.Run("ErrorOpening", func(t *testing.T) {
		f := os.NewFile(0200, fn)
		f.Close()

		//assert.Nil(t, f)
		db, err := setupDB(f.Name())
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, db.Path(), "test.db")
	})

	err := os.Remove(fn)
	assert.NoError(t, err)

}
