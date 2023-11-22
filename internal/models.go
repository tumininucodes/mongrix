package internal

type Todo struct {
	Name string `bson:"name"`
	Description string `bson:"description"`
}