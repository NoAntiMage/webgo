package stringx

import (
	"testing"
)

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person
	Grade int
	Class int
}

func TestStructFieldNames(t *testing.T) {
	expected := []string{"Name", "Age"}
	t.Run("1", func(t *testing.T) {
		res := StructFieldNames(&Person{})
		if len(res) != len(expected) {
			t.Fatal("fail")
		}
		for i := range res {
			if res[i] != expected[i] {
				t.Fatal("fail")
			}
		}
	})

}

func TestStructFieldNamesRecursion(t *testing.T) {
	expected := []string{"Name", "Age", "Grade", "Class"}
	t.Run("2", func(t *testing.T) {
		res := StructFieldNames(&Student{})
		if len(res) != len(expected) {
			t.Fatal("fail")
		}
		for i := range res {
			if res[i] != expected[i] {
				t.Fatal("fail")
			}
		}
	})
}
