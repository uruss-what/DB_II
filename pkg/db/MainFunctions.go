package db

import (
	"DB_II/pkg/interfaces"
	"fmt"
	"unicode"
)

var (
	ErrEmptyName        = fmt.Errorf("name cannot be empty")
	ErrInvalidName      = fmt.Errorf("name contains invalid characters")
	ErrInvalidOperation = fmt.Errorf("invalid operation")
)

func isValidName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
			return ErrInvalidName
		}
	}
	return nil
}

func (db *Database) getPool(poolName string) (*DataPool, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	pool, ok := db.Pools[poolName]
	if !ok {
		return nil, ErrPoolNotFound
	}
	return pool, nil
}

func (db *Database) getSchema(poolName, schemaName string) (*DataSchema, error) {
	pool, err := db.getPool(poolName)
	if err != nil {
		return nil, err
	}

	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	schema, ok := pool.Schemas[schemaName]
	if !ok {
		return nil, ErrSchemaNotFound
	}
	return schema, nil
}

func (db *Database) getCollection(poolName, schemaName, collectionName string) (interfaces.CollectionInterface, error) {
	schema, err := db.getSchema(poolName, schemaName)
	if err != nil {
		return nil, err
	}

	schema.mutex.RLock()
	defer schema.mutex.RUnlock()

	collection, ok := schema.Collections[collectionName]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}

func (db *Database) ListPools() []string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	pools := make([]string, 0, len(db.Pools))
	for _, pool := range db.Pools {
		pools = append(pools, pool.Name)
	}
	return pools
}

func (db *Database) ListSchemas(poolName string) ([]string, error) {
	pool, err := db.getPool(poolName)
	if err != nil {
		return nil, err
	}

	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	schemas := make([]string, 0, len(pool.Schemas))
	for _, schema := range pool.Schemas {
		schemas = append(schemas, schema.Name)
	}
	return schemas, nil
}

func (db *Database) ListCollections(poolName, schemaName string) ([]string, error) {
	schema, err := db.getSchema(poolName, schemaName)
	if err != nil {
		return nil, err
	}

	schema.mutex.RLock()
	defer schema.mutex.RUnlock()

	collections := make([]string, 0, len(schema.Collections))
	for name := range schema.Collections {
		collections = append(collections, name)
	}
	return collections, nil
}
