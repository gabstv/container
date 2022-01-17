package linq_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/gabstv/container/linq"
	"github.com/stretchr/testify/assert"
)

type Car struct {
	Name  string
	Year  int
	Price float64
}

func TestLinq(t *testing.T) {
	cars := []Car{
		{Name: "Ferrari F40", Year: 2009, Price: 190000},
		{Name: "Ferrari S20", Year: 2004, Price: 200000},
		{Name: "Ferrari F50", Year: 2008, Price: 100000},
		{Name: "Aladdin UX", Year: 2005, Price: 80000},
		{Name: "Tucson Plus", Year: 2006, Price: 50000},
		{Name: "Peugeot 706", Year: 2007, Price: 133212.90},
	}

	names := linq.From[Car, string](cars).Where(func(c Car) bool {
		return c.Price >= 100000
	}).Select(func(c Car) string {
		return c.Name
	}).Sort(func(name1, name2 string) bool {
		return name1 < name2
	}).All()

	assert.Equal(t, 4, len(names))
	assert.Equal(t, "Ferrari F40", names[0])
	assert.Equal(t, "Ferrari F50", names[1])
	assert.Equal(t, "Ferrari S20", names[2])
	assert.Equal(t, "Peugeot 706", names[3])

	// order by year
	// declare a struct on the fly as the result type
	result := linq.From[Car, struct {
		Name string
		Year int
	}](cars).Sort(func(c1, c2 Car) bool {
		return c1.Year < c2.Year
	}).Select(func(c Car) struct {
		Name string
		Year int
	} {
		return struct {
			Name string
			Year int
		}{
			Name: c.Name,
			Year: c.Year,
		}
	}).All()
	assert.Equal(t, 6, len(result))
	assert.Equal(t, 2004, result[0].Year)
	assert.Equal(t, 2009, result[5].Year)
	assert.Equal(t, "Aladdin UX", result[1].Name)
}

type Person struct {
	Name string
	Age  int
}

func TestLinqMap(t *testing.T) {
	people := make(map[string]Person)
	people["John"] = Person{Name: "John", Age: 30}
	people["Mary"] = Person{Name: "Mary", Age: 25}
	people["Peter"] = Person{Name: "Peter", Age: 20}
	people["Paul"] = Person{Name: "Paul", Age: 40}
	pplSlice := container.MapToSlice(people)

	q := linq.From[container.MI[string, Person], container.MI[string, Person]](pplSlice)
	mres := q.Where(func(item container.MI[string, Person]) bool {
		return item.Value.Age >= 30
	}).All()
	peopleRes := container.SliceToMap(mres)
	assert.Equal(t, 2, len(peopleRes))
	assert.Equal(t, "John", peopleRes["John"].Name)
	assert.Equal(t, 40, peopleRes["Paul"].Age)
	assert.Equal(t, "", peopleRes["Mary"].Name)
}
