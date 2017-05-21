package stl

// This file contains generic 3D transformation stuff

import (
	"math"
)

// Calculates a 4x4 rotation matrix for a rotation of angle in radians
// around a axis defined by dir.
// The returns a rotation matrix.
func RotationMatrix(pos Vec3, dir Vec3, angle float64, rotationMatrix *Mat4) {
	dirN := dir.unitVec()

	s := math.Sin(angle)
	c := math.Cos(angle)

	u := float64(dirN[0])
	v := float64(dirN[1])
	w := float64(dirN[2])

	uu := u * u
	uv := u * v
	uw := u * w
	vv := v * v
	vw := v * w
	ww := w * w
	us := u * s
	vs := v * s
	ws := w * s
	
	iuu := 1 - uu
	ivv := 1 - vv
	iww := 1 - ww
	ic := 1 - c
	
	mRotate := Mat4{
		Vec4{uu + iuu * c, uv * ic - ws, uw * ic + vs, 0},
		Vec4{uv * ic + ws, vv + ivv * c, vw * ic - us, 0},
		Vec4{uw * ic - vs, vw * ic + us, ww + iww * c, 0},
		Vec4{0, 0, 0, 1},
	}
	
	mTranslateForward := Mat4{
		Vec4{1, 0, 0, float64(pos[0])},
		Vec4{0, 1, 0, float64(pos[1])},
		Vec4{0, 0, 1, float64(pos[2])},
		Vec4{0, 0, 0, 1},
	}

	mTranslateBackward := Mat4{
		Vec4{1, 0, 0, float64(-pos[0])},
		Vec4{0, 1, 0, float64(-pos[1])},
		Vec4{0, 0, 1, float64(-pos[2])},
		Vec4{0, 0, 0, 1},
	}

	var mTmp Mat4
	mRotate.MultMat4(&mTranslateForward, &mTmp)
	mTranslateBackward.MultMat4(&mTmp, rotationMatrix)
}