package person

// Person ...
type Person struct {
	name string
	age  int
}

// NewPerson ...
func NewPerson(name string, age int) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

// Set ...
func (p *Person) Set(name string, age int) {
	p.name = name
	p.age = age
}

// Get ...
func (p *Person) Get() (name string, age int) {
	return p.name, p.age
}
