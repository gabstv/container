package container

import (
	"constraints"
	"encoding/json"
)

type List2D[T any] struct {
	data   []T
	width  int
	height int
}

func min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func (l *List2D[T]) Resize(width, height int) {
	prevl := l.Copy()
	l.data = make([]T, width*height)
	l.width = width
	l.height = height
	for x := 0; x < min(width, prevl.Width()); x++ {
		for y := 0; y < min(height, prevl.Height()); y++ {
			l.Set(x, y, prevl.Get(x, y))
		}
	}
}

func (l *List2D[T]) Set(x, y int, v T) {
	l.data[x+y*l.width] = v
}

func (l *List2D[T]) Get(x, y int) T {
	return l.data[x+y*l.width]
}

// GetW is the same as Get, but wraps around if is out of bounds. It returns the
// zero value if the width or height is zero
func (l *List2D[T]) GetW(x, y int) T {
	if l.width < 1 || l.height < 1 {
		var zv T
		return zv
	}
	for x < 0 {
		x += l.width
	}
	for x >= l.width {
		x -= l.width
	}
	for y < 0 {
		y += l.height
	}
	for y >= l.height {
		y -= l.height
	}
	return l.data[x+y*l.width]
}

// Copy returns a copy of the list
func (l *List2D[T]) Copy() *List2D[T] {
	l2 := &List2D[T]{
		data:   make([]T, len(l.data)),
		width:  l.width,
		height: l.height,
	}
	copy(l2.data, l.data)
	return l2
}

// Copy returns a copy of the list determined by the given parameters.
// It panics if the given parameters are out of bounds.
func (l *List2D[T]) CopyRect(x, y, width, height int) *List2D[T] {
	l2 := &List2D[T]{
		data:   make([]T, width*height),
		width:  width,
		height: height,
	}
	for lx := 0; lx < width; lx++ {
		for ly := 0; ly < height; ly++ {
			l2.Set(lx, ly, l.Get(x+lx, y+ly))
		}
	}
	return l2
}

// ShiftX shifts all the elements in the x and y dimensions of the list by
// the offsets.
func (l *List2D[T]) Shift(xoffset, yoffset int) {
	if l.width < 1 || l.height < 1 {
		return
	}
	l2 := l.Copy()
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			l.Set(x, y, l2.GetW(x-xoffset, y-yoffset))
		}
	}
}

func (l *List2D[T]) Width() int {
	return l.width
}

func (l *List2D[T]) Height() int {
	return l.height
}

func NewList2D[T any](width, height int) *List2D[T] {
	l := new(List2D[T])
	l.width = width
	l.height = height
	l.data = make([]T, width*height)
	return l
}

type list2dj[T any] struct {
	W int `json:"w"`
	H int `json:"h"`
	D []T `json:"d"`
}

func (l List2D[T]) MarshalText() (text []byte, err error) {
	j := list2dj[T]{
		W: l.width,
		H: l.height,
		D: l.data,
	}
	return json.Marshal(j)
}

func (l *List2D[T]) UnmarshalText(text []byte) error {
	j := new(list2dj[T])
	err := json.Unmarshal(text, j)
	if err != nil {
		return err
	}
	l.width = j.W
	l.height = j.H
	l.data = j.D
	return nil
}

func (l List2D[T]) MarshalJSON() ([]byte, error) {
	j := list2dj[T]{
		W: l.width,
		H: l.height,
		D: l.data,
	}
	return json.Marshal(j)
}

func (l *List2D[T]) UnmarshalJSON(text []byte) error {
	j := new(list2dj[T])
	err := json.Unmarshal(text, j)
	if err != nil {
		return err
	}
	l.width = j.W
	l.height = j.H
	l.data = j.D
	return nil
}
