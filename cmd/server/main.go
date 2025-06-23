package main

import (
	"DB_II/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var database *db.Database

func init() {
	postgresDB, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	authManager := db.NewAuthManager(postgresDB)

	database = db.NewDatabase(authManager)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var cmd db.Command
		err := decoder.Decode(&cmd)
		if err != nil {
			log.Printf("Error decoding command: %v", err)
			return
		}

		var response interface{}
		var responseErr error

		if cmd.Operation != "register" {
			role, err := database.AuthManager.ValidateUser(cmd.Username, cmd.Password)
			if err != nil {
				responseErr = fmt.Errorf("authentication failed: %v", err)
				encoder.Encode(map[string]string{
					"status": "error",
					"error":  responseErr.Error(),
				})
				continue
			}
			cmd.Role = role
		}

		switch cmd.Operation {
		case "register":
			responseErr = database.AuthManager.RegisterUser(cmd.Username, cmd.Password, cmd.Role)

		case "create_pool":
			responseErr = database.CreatePool(cmd.Username, cmd.Pool)

		case "create_schema":
			responseErr = database.CreateSchema(cmd.Username, cmd.Pool, cmd.Schema)

		case "create_collection":
			responseErr = database.CreateCollection(cmd.Username, cmd.Pool, cmd.Schema, cmd.Collection, cmd.TreeType)

		case "set":
			if !database.AuthManager.HasPermission(cmd.Username, db.PermWrite) {
				responseErr = db.ErrPermissionDenied
				break
			}

			pool, exists := database.Pools[cmd.Pool]
			if !exists {
				responseErr = db.ErrPoolNotFound
				break
			}

			schema, exists := pool.Schemas[cmd.Schema]
			if !exists {
				responseErr = db.ErrSchemaNotFound
				break
			}

			collection, exists := schema.Collections[cmd.Collection]
			if !exists {
				responseErr = db.ErrCollectionNotFound
				break
			}

			response = collection.Set(cmd.Key, cmd.SecondaryKey, cmd.Value)

		case "get":
			if !database.AuthManager.HasPermission(cmd.Username, db.PermRead) {
				responseErr = db.ErrPermissionDenied
				break
			}

			pool, exists := database.Pools[cmd.Pool]
			if !exists {
				responseErr = db.ErrPoolNotFound
				break
			}

			schema, exists := pool.Schemas[cmd.Schema]
			if !exists {
				responseErr = db.ErrSchemaNotFound
				break
			}

			collection, exists := schema.Collections[cmd.Collection]
			if !exists {
				responseErr = db.ErrCollectionNotFound
				break
			}

			value, status := collection.Get(cmd.Key)
			response = map[string]string{
				"value":  value,
				"status": status,
			}

		default:
			responseErr = fmt.Errorf("unknown operation: %s", cmd.Operation)
		}

		if responseErr != nil {
			encoder.Encode(map[string]string{
				"status": "error",
				"error":  responseErr.Error(),
			})
		} else {
			encoder.Encode(map[string]interface{}{
				"status":   "ok",
				"response": response,
			})
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer listener.Close()

	log.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
