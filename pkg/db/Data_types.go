package db

import (
	"DB_II/pkg/interfaces"
	"errors"
	"sync"
)

var (
	ErrPoolExists         = errors.New("pool already exists")
	ErrPoolNotFound       = errors.New("pool not found")
	ErrSchemaExists       = errors.New("schema already exists")
	ErrSchemaNotFound     = errors.New("schema not found")
	ErrCollectionExists   = errors.New("collection already exists")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrPermissionDenied   = errors.New("permission denied")
)

type Database struct {
	Pools       map[string]*DataPool
	AuthManager *AuthManager
	mutex       *sync.RWMutex
}

func NewDatabase(authManager *AuthManager) *Database {
	return &Database{
		Pools:       make(map[string]*DataPool),
		AuthManager: authManager,
		mutex:       &sync.RWMutex{},
	}
}

type DataPool struct {
	Name    string
	Schemas map[string]*DataSchema
	mutex   *sync.RWMutex
}

func NewDataPool(name string) *DataPool {
	return &DataPool{
		Name:    name,
		Schemas: make(map[string]*DataSchema),
		mutex:   &sync.RWMutex{},
	}
}

type DataSchema struct {
	Name        string
	Collections map[string]interfaces.CollectionInterface
	mutex       *sync.RWMutex
}

func NewDataSchema(name string) *DataSchema {
	return &DataSchema{
		Name:        name,
		Collections: make(map[string]interfaces.CollectionInterface),
		mutex:       &sync.RWMutex{},
	}
}

func (db *Database) CreatePool(username string, poolName string) error {
	if !db.AuthManager.HasPermission(username, PermCreatePool) {
		return ErrPermissionDenied
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.Pools[poolName]; exists {
		return ErrPoolExists
	}

	db.Pools[poolName] = NewDataPool(poolName)
	return nil
}

func (db *Database) CreateSchema(username string, poolName, schemaName string) error {
	if !db.AuthManager.HasPermission(username, PermCreateSchema) {
		return ErrPermissionDenied
	}

	db.mutex.RLock()
	pool, exists := db.Pools[poolName]
	db.mutex.RUnlock()

	if !exists {
		return ErrPoolNotFound
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if _, exists := pool.Schemas[schemaName]; exists {
		return ErrSchemaExists
	}

	pool.Schemas[schemaName] = NewDataSchema(schemaName)
	return nil
}

func (db *Database) CreateCollection(username string, poolName, schemaName, collectionName string, treeType TreeType) error {
	if !db.AuthManager.HasPermission(username, PermCreateCollection) {
		return ErrPermissionDenied
	}

	db.mutex.RLock()
	pool, exists := db.Pools[poolName]
	db.mutex.RUnlock()

	if !exists {
		return ErrPoolNotFound
	}

	pool.mutex.RLock()
	schema, exists := pool.Schemas[schemaName]
	pool.mutex.RUnlock()

	if !exists {
		return ErrSchemaNotFound
	}

	schema.mutex.Lock()
	defer schema.mutex.Unlock()

	if _, exists := schema.Collections[collectionName]; exists {
		return ErrCollectionExists
	}

	schema.Collections[collectionName] = NewTreeCollection(treeType)
	return nil
}

type Command struct {
	Username     string
	Password     string
	Role         Role
	Operation    string
	Pool         string
	Schema       string
	Collection   string
	TreeType     TreeType
	Key          string
	Value        string
	SecondaryKey string
}
