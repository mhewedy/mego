package rooms

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	email = iota
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

func ListRooms(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	once.Do(loadRoomList)

	return ListRoomEmails(), nil
}

func ListRoomsTree(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	once.Do(loadRoomList)

	return roomTree, nil
}

func ListRoomEmails() []string {
	once.Do(loadRoomList)

	roomCodes := make([]string, 0)
	for i, rr := range roomList {
		if i == 0 {
			continue // skip header
		}
		roomCodes = append(roomCodes, rr[email])
	}

	return roomCodes
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

func FindByEmail(toFound string) (string, error) {
	once.Do(loadRoomList)

	for i, rr := range roomList {
		if i == 0 {
			continue // skip header
		}
		if rr[email] == toFound {
			return rr[displayName], nil
		}
	}
	return "", fmt.Errorf("no such room found: %s", toFound)
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
				key := fmt.Sprintf("%s-%s", buildingKey, field)
				if !contains(b.Children, key) {
					b.Children = append(b.Children, Node{
						Key:   key,
						Label: field,
					})
				}
			}
			if j == size {
				buildingKey := row[j-2]
				b := get(tree, buildingKey)
				zoneKey := fmt.Sprintf("%s-%s", buildingKey, row[j-1])
				z := get(b.Children, zoneKey)
				key := fmt.Sprintf("%s-%s", buildingKey, field)
				if !contains(z.Children, key) {
					z.Children = append(z.Children, Node{
						Key:   key,
						Label: fmt.Sprintf("%s Person", field),
					})
				}
			}
			if j == displayName {
				buildingKey := row[j-3]
				b := get(tree, buildingKey)
				zoneKey := fmt.Sprintf("%s-%s", buildingKey, row[j-2])
				z := get(b.Children, zoneKey)
				sizeKey := fmt.Sprintf("%s-%s", buildingKey, row[j-1])
				s := get(z.Children, sizeKey)
				key := row[email]
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
