package main

import "fmt"
// Simple factory
// 구조체를 초기화 할 때, Age 속성을 초기화해야하는 것을 잊을 수도 있다.
// NewPerson 팩토리 메서드는 Person의 인스턴스를 만들 때,
// 이름과 나이가 모두 제공되어야 생성되고, nil 값이 없음을 보장하도록 해준다.
type SimplePerson struct {
	Name string
	Age int
}

func (p *SimplePerson) Greet() {
	fmt.Printf("Hi! My Name is %s", p.Name)
}

func NewSPerson(name string, age int) *SimplePerson {
	return &SimplePerson {
		Name : name,
		Age : age,
	}
}

// Interface Factory
// 인터페이스로 구조체 내부 구현 감추기
type Person interface {
	Greet()
}

type person struct {
	name string
	age int
}

func (p *person) Greet() {
	fmt.Printf("Hi! My name is %s", p.name)
}

func NewPerson(name string, age int) Person {
	return &person {
		name: name,
		age: age,
	}
}

// Factory generators
// Factory method

type Animal struct {
	species string
	age		int
}

type AnimalHouse struct {
	name 			string
	sizeInMeters	int
}

type AnimalFactory struct {
	species 	string
	houseName	string
}

func (af *AnimalFactory) NewAnimal(age int) *Animal {
	return &Animal {
		species: af.species,
		age:	age,
	}
}

func (af *AnimalFactory) NewHouse(sizeInMeters int) *AnimalHouse{
	return &AnimalHouse{
		name:	af.houseName,
		sizeInMeters: sizeInMeters,
	}
}

func (a *Animal) String() string {
	return fmt.Sprintf("Species : %s , Age :%d", a.species, a.age)
}

func (a *AnimalHouse) String() string {
	return fmt.Sprintf("Name : %s, HouseSize : %d",a.name,a.sizeInMeters)
}

func testFactoryMethod() {
	dogFactory := AnimalFactory{
		species: "dog",
		houseName: "kennel",
	}

	dog := dogFactory.NewAnimal(2)
	kennel := dogFactory.NewHouse(3)

	fmt.Println(dog.String())
	fmt.Println(kennel.String())

}
// Factory func
// Closure 이용.
type Toy struct {
	species string
	size 	string
}

func NewToyFactory(species string) func(size string) *Toy {
	return func(size string) *Toy {
		return &Toy{
			species : species,
			size  : size,
		}
	}
}

func (t *Toy) String() string {
	return fmt.Sprintf("Species : %s , Size :%s", t.species, t.size)
}

func testFactoryFunc(){
	newRobotToy := NewToyFactory("Robot")
	robotToy := newRobotToy("Small")
	newCarToy := NewToyFactory("Car")
	carToy := newCarToy("Big")
	fmt.Println(robotToy)
	fmt.Println(carToy)

}

func main() {
	testFactoryFunc()
}