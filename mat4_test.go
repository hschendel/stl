package stl

import (
	"testing"
)

func BenchmarkMultMat4(b *testing.B) {
	r := Mat4{
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
	}
	s := Mat4{
		Vec4{9, 2, 1, -6},
		Vec4{8, 3, 0, -5},
		Vec4{7, 4, -1, -4},
		Vec4{6, 5, -2, -3},
	}
	var t Mat4
	for i := 0; i < b.N; i++ {
		r.MultMat4(&s, &t)
	}
}

func TestMultMat4(t *testing.T) {
	m := Mat4{
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
		Vec4{1, 2, 3, 4},
	}
	var r Mat4
	m.MultMat4(&Mat4Identity, &r)
	if m != r {
		t.Errorf("Result: %v, Expected: %v", r, m)
	}
}

func BenchmarkMultVec3(b *testing.B) {
	m := Mat4{
		Vec4{9, 2, 1, -6},
		Vec4{8, 3, 0, -5},
		Vec4{7, 4, -1, -4},
		Vec4{6, 5, -2, -3},
	}
	v := Vec3{-1000, 234, 1000}
	for i := 0; i < b.N; i++ {
		_ = m.MultVec3(v)
	}
}

func TestMultVec3(t *testing.T) {
	m := Mat4{
		Vec4{1, 0, 0, 1000},
		Vec4{0, 2, 0, 500},
		Vec4{0, 0, 1, 250},
		Vec4{0, 0, 0, 1},
	}
	v := Vec3{1, 1, 1}
	r := m.MultVec3(v)
	expected := Vec3{1001, 502, 251}
	if r != expected {
		t.Errorf("MultVec3 result: %v, expected: %v", v, expected)
	}
}
