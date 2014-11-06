package main

import "fmt"

type User struct {
	Name string
	Age  int
	Tags []string
}

func (u *User) String() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

func main() {
	users := map[string]*User{
		"rob":   &User{Name: "Rob Pike", Age: 58, Tags: []string{"plan9"}},
		"linus": {Name: "Linus Torvalds", Age: 44, Tags: []string{"linux"}},
	}

	for _, u := range users {
		fmt.Println(u)
	}
}
