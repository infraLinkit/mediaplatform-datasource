package cmd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Smoke test: migrateEntities slice harus terisi dengan pointer ke struct.
// Jaga supaya tidak accidentally kosong / di-truncate.
func TestMigrateEntities_NotEmpty(t *testing.T) {
	assert.NotEmpty(t, migrateEntities, "migrateEntities tidak boleh kosong")
	assert.GreaterOrEqual(t, len(migrateEntities), 50, "minimal 50 entity terdaftar (current: %d)", len(migrateEntities))
}

func TestMigrateEntities_AllPointers(t *testing.T) {
	for i, e := range migrateEntities {
		v := reflect.ValueOf(e)
		assert.Equal(t, reflect.Ptr, v.Kind(), "entity[%d] harus pointer (GORM convention)", i)
		assert.Equal(t, reflect.Struct, v.Elem().Kind(), "entity[%d] harus pointer ke struct", i)
	}
}

func TestMigrateEntities_NoNil(t *testing.T) {
	for i, e := range migrateEntities {
		assert.NotNil(t, e, "entity[%d] tidak boleh nil", i)
	}
}

func TestMigrateEntities_NoDuplicate(t *testing.T) {
	seen := make(map[string]int)
	for i, e := range migrateEntities {
		typeName := reflect.TypeOf(e).String()
		if prev, ok := seen[typeName]; ok {
			t.Errorf("duplicate entity %s di index %d (sebelumnya di %d)", typeName, i, prev)
		}
		seen[typeName] = i
	}
}

// === migrate sub-cmd registered ===

func TestMigrateCmd_Registered(t *testing.T) {
	assert.NotNil(t, migrateCmd)
	assert.Equal(t, "migrate", migrateCmd.Use)
}

func TestServerCmd_Registered(t *testing.T) {
	assert.NotNil(t, serverCmd)
	assert.Equal(t, "server", serverCmd.Use)
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	subs := rootCmd.Commands()
	names := make(map[string]bool)
	for _, c := range subs {
		names[c.Use] = true
	}
	assert.True(t, names["server"], "rootCmd harus punya 'server' sub-cmd")
	assert.True(t, names["migrate"], "rootCmd harus punya 'migrate' sub-cmd")
}
