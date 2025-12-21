package printobject

import (
	"encoding/json"
	"fmt"
)

// Utils =)
func PrintObject(message string, object any) {
	p, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message, string(p))
}
