package models

type Document struct {
	Id   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}
