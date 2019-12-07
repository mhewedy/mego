package rooms

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	buildRoomTree()
}

func buildRoomTree() {

	root := Root{Root: []Node{}}

	for i, row := range roomList {
		if i == 0 {
			continue // skip header
		}
		for j, field := range row {
			//field := strings.TrimSpace(field)
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
				b := find(root.Root, buildingKey)
				if !contains(b.Children, field) {
					b.Children = append(b.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == size {
				buildingKey := row[j-2]
				b := find(root.Root, buildingKey)
				zoneKey := row[j-1]
				z := find(b.Children, zoneKey)
				if !contains(z.Children, field) {
					z.Children = append(z.Children, Node{
						Key:   field,
						Label: field,
					})
				}
			}
			if j == displayName {
				buildingKey := row[j-3]
				b := find(root.Root, buildingKey)
				zoneKey := row[j-2]
				z := find(b.Children, zoneKey)
				sizeKey := row[j-1]
				s := find(z.Children, sizeKey)
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
	return find(nodes, key) != nil
}

func find(nodes []Node, key string) *Node {
	for i := range nodes {
		if nodes[i].Key == key {
			return &nodes[i]
		}
	}
	return nil
}
