package master

import (
	"encoding/json"
	"reflect"
	"strings"

	"go-micro.dev/v5/registry"
)

func workNodeDiff(old map[string]*registry.Node, new map[string]*registry.Node) (add []string, del []string, chg []string) {
	for key, newNode := range new {
		if oldNode, has := old[key]; has {
			if !reflect.DeepEqual(oldNode, newNode) {
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

// func GetNodeID(name string) (string, error) {
// 	nodeID := strings.Split(name, "|")
// 	if len(nodeID) < 2 {
// 		return "", errors.New("get node id failed")
// 	}
// 	id := nodeID[0]
// 	return id, nil
// }

func Encode(s *TaskSpec) string {
	buff, _ := json.Marshal(s)
	return string(buff)
}

func Decode(ds []byte) (*TaskSpec, error) {
	var tmp *TaskSpec
	if err := json.Unmarshal(ds, &tmp); err != nil {
		return nil, err
	} else {
		return tmp, nil
	}
}

// func getTaskPath(name string) string {
// 	return fmt.Sprintf("%s/%s", global.TaskPath, name)
// }

func getLeaderAddr(str string) string {
	s := strings.Split(str, "@")
	if len(s) < 2 {
		return ""
	}
	return s[1]
}
