package main

import (
	"DB_II/pkg/db"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn     net.Conn
	username string
	password string
}

func NewClient() (*Client, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) register() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	c.username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	c.password = strings.TrimSpace(password)

	fmt.Println("\nSelect role:")
	fmt.Println("1. Superuser")
	fmt.Println("2. Admin")
	fmt.Println("3. Editor")
	fmt.Println("4. User")

	var role db.Role
	for {
		fmt.Print("Enter choice (1-4): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			role = db.RoleSuperUser
		case "2":
			role = db.RoleAdmin
		case "3":
			role = db.RoleEditor
		case "4":
			role = db.RoleUser
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue
		}
		break
	}

	cmd := db.Command{
		Operation: "register",
		Username:  c.username,
		Password:  c.password,
		Role:      role,
	}

	response, err := c.sendCommand(cmd)
	if err != nil {
		return err
	}

	if status, ok := response["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v", response["error"])
	}

	fmt.Printf("User '%s' registered successfully with role %s\n", c.username, role)
	return nil
}

func (c *Client) login() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	c.username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	c.password = strings.TrimSpace(password)

	cmd := db.Command{
		Operation: "create_pool",
		Username:  c.username,
		Password:  c.password,
		Pool:      "",
	}

	response, err := c.sendCommand(cmd)
	if err != nil {
		return err
	}

	if status, ok := response["status"].(string); ok && status == "error" {
		if strings.Contains(response["error"].(string), "authentication failed") {
			return fmt.Errorf("invalid credentials")
		}
	}

	return nil
}

func (c *Client) sendCommand(cmd db.Command) (map[string]interface{}, error) {
	encoder := json.NewEncoder(c.conn)
	decoder := json.NewDecoder(c.conn)

	if cmd.Operation != "register" {
		cmd.Password = c.password
	}

	if err := encoder.Encode(cmd); err != nil {
		return nil, fmt.Errorf("failed to send command: %v", err)
	}

	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return response, nil
}

func (c *Client) createPool() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter pool name: ")
	poolName, _ := reader.ReadString('\n')
	poolName = strings.TrimSpace(poolName)

	cmd := db.Command{
		Operation: "create_pool",
		Username:  c.username,
		Pool:      poolName,
	}

	response, err := c.sendCommand(cmd)
	if err != nil {
		return err
	}

	if status, ok := response["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v", response["error"])
	}

	fmt.Printf("Pool '%s' created successfully\n", poolName)
	return nil
}

func (c *Client) createSchema() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter pool name: ")
	poolName, _ := reader.ReadString('\n')
	poolName = strings.TrimSpace(poolName)

	fmt.Print("Enter schema name: ")
	schemaName, _ := reader.ReadString('\n')
	schemaName = strings.TrimSpace(schemaName)

	cmd := db.Command{
		Operation: "create_schema",
		Username:  c.username,
		Pool:      poolName,
		Schema:    schemaName,
	}

	response, err := c.sendCommand(cmd)
	if err != nil {
		return err
	}

	if status, ok := response["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v", response["error"])
	}

	fmt.Printf("Schema '%s' created successfully in pool '%s'\n", schemaName, poolName)
	return nil
}

func (c *Client) createCollection() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter pool name: ")
	poolName, _ := reader.ReadString('\n')
	poolName = strings.TrimSpace(poolName)

	fmt.Print("Enter schema name: ")
	schemaName, _ := reader.ReadString('\n')
	schemaName = strings.TrimSpace(schemaName)

	fmt.Print("Enter collection name: ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)

	fmt.Println("Select tree type:")
	fmt.Println("1. AVL Tree")
	fmt.Println("2. Red-Black Tree")
	fmt.Println("3. B-Tree")

	var treeType db.TreeType
	for {
		fmt.Print("Enter choice (1-3): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			treeType = db.TreeTypeAVL
		case "2":
			treeType = db.TreeTypeRedBlack
		case "3":
			treeType = db.TreeTypeBTree
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue
		}
		break
	}

	cmd := db.Command{
		Operation:  "create_collection",
		Username:   c.username,
		Pool:       poolName,
		Schema:     schemaName,
		Collection: collectionName,
		TreeType:   treeType,
	}

	response, err := c.sendCommand(cmd)
	if err != nil {
		return err
	}

	if status, ok := response["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v", response["error"])
	}

	fmt.Printf("Collection '%s' created successfully with %s tree type\n", collectionName, treeType)
	return nil
}

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Welcome to DB II Client")
	fmt.Println("------------------------")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nDo you want to:")
	fmt.Println("1. Register new user")
	fmt.Println("2. Login")

	for {
		fmt.Print("\nEnter choice (1-2): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		var err error
		switch choice {
		case "1":
			err = client.register()
		case "2":
			err = client.login()
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		break
	}

	fmt.Printf("\nWelcome, %s!\n", client.username)

	for {
		fmt.Println("\nAvailable commands:")
		fmt.Println("1. Create Pool")
		fmt.Println("2. Create Schema")
		fmt.Println("3. Create Collection")
		fmt.Println("4. Exit")

		fmt.Print("\nEnter command (1-4): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		var err error

		switch choice {
		case "1":
			err = client.createPool()
		case "2":
			err = client.createSchema()
		case "3":
			err = client.createCollection()
		case "4":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command")
			continue
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
