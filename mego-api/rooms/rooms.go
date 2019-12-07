package rooms

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	code = iota
	building
	zone
	size
	displayName
)

var roomList [][]string
var roomTree *Root

type Root struct {
	Root []Node `json:"root"`
}
type Node struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Children []Node `json:"children,omitempty"`
}

var once sync.Once

func ListRoomsTree(w http.ResponseWriter, r *http.Request) {
	once.Do(loadRoomList)

	json.NewEncoder(w).Encode(roomTree)
}

func loadRoomList() {

	file, err := os.OpenFile("rooms.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	roomList, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// trim spaces
	for i := range roomList {
		for j := range roomList[i] {
			roomList[i][j] = strings.TrimSpace(roomList[i][j])
		}
	}

	buildRoomTree()
}

func buildRoomTree() {

	root := Root{Root: []Node{}}

	for i, row := range roomList {
		if i == 0 {
			continue // skip header
		}
		for j, field := range row {
			if j == building {
				if !contains(root.Root, field) {
					root.Root = append(root.Root, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == zone {
				buildingKey := row[j-1]
				b, _ := find(root.Root, buildingKey)
				if !contains(b.Children, field) {
					b.Children = append(b.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == size {
				buildingKey := row[j-2]
				b, _ := find(root.Root, buildingKey)
				zoneKey := row[j-1]
				z, _ := find(b.Children, zoneKey)
				if !contains(z.Children, field) {
					z.Children = append(z.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == displayName {
				buildingKey := row[j-3]
				b, _ := find(root.Root, buildingKey)
				zoneKey := row[j-2]
				z, _ := find(b.Children, zoneKey)
				sizeKey := row[j-1]
				s, _ := find(z.Children, sizeKey)
				key := row[code]
				if !contains(s.Children, key) {
					s.Children = append(s.Children, Node{
						Key:   key,
						Label: field,
					})
				}
			}
		}
	}
	roomTree = &root
}

func contains(nodes []Node, key string) bool {
	for i := range nodes {
		if nodes[i].Key == key {
			return true
		}
	}
	return false
}

func find(nodes []Node, key string) (*Node, error) {
	for i := range nodes {
		if nodes[i].Key == key {
			return &nodes[i], nil
		}
	}
	return nil, fmt.Errorf("should never happen: key %s not found in nodes %s", key, nodes)
}
