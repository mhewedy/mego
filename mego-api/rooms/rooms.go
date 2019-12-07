package rooms

import (
	"encoding/csv"
	"encoding/json"
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
var roomTree []Node

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

	file, err := os.Open("rooms.csv")
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

	var tree []Node

	for i, row := range roomList {
		if i == 0 {
			continue // skip header
		}
		for j, field := range row {
			if j == building {
				if !contains(tree, field) {
					tree = append(tree, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == zone {
				buildingKey := row[j-1]
				b := get(tree, buildingKey)
				if !contains(b.Children, field) {
					b.Children = append(b.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == size {
				buildingKey := row[j-2]
				b := get(tree, buildingKey)
				zoneKey := row[j-1]
				z := get(b.Children, zoneKey)
				if !contains(z.Children, field) {
					z.Children = append(z.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == displayName {
				buildingKey := row[j-3]
				b := get(tree, buildingKey)
				zoneKey := row[j-2]
				z := get(b.Children, zoneKey)
				sizeKey := row[j-1]
				s := get(z.Children, sizeKey)
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
	roomTree = tree
}

func contains(nodes []Node, key string) bool {
	for i := range nodes {
		if nodes[i].Key == key {
			return true
		}
	}
	return false
}

func get(nodes []Node, key string) *Node {
	for i := range nodes {
		if nodes[i].Key == key {
			return &nodes[i]
		}
	}
	log.Fatalf("should never happen: key \"%s\" not found in nodes \"%#v\"\n", key, nodes)
	return nil
}