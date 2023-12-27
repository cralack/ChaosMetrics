package master

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func workNodeDiff(old map[string]*NodeSpec, new map[string]*NodeSpec) (add []string, del []string, chg []string) {
	for key, newNode := range new {
		if oldNode, has := old[key]; has {
			if !reflect.DeepEqual(oldNode.Node, newNode.Node) {
				chg = append(chg, key)
			}
		} else { // !has
			add = append(add, key)
		}
	}
	for key := range old {
		if _, has := new[key]; !has {
			del = append(del, key)
		}
	}
	return
}

func getNodeID(assigned string) (string, error) {
	nodeID := strings.Split(assigned, "|")
	if len(nodeID) < 2 {
		return "", errors.New("get node id failed")
	}
	id := nodeID[0]
	return id, nil
}

func encode(s *ResourceSpec) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func genMasterID(id string, ipv4 string, GRPCAddress string) string {
	return "master" + id + "-" + ipv4 + GRPCAddress
}

func decode(ds []byte) (*ResourceSpec, error) {
	var s *ResourceSpec
	err := json.Unmarshal(ds, &s)
	return s, err
}
