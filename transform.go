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

	ic := 1 - c

	uv_ic := uv * ic
	uw_ic := uw * ic
	vw_ic := vw * ic

	mRotate := Mat4{
		Vec4{uu + (1-uu)*c, uv_ic - ws, uw_ic + vs, 0},
		Vec4{uv_ic + ws, vv + (1-vv)*c, vw_ic - us, 0},
		Vec4{uw_ic - vs, vw_ic + us, ww + (1-ww)*c, 0},
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
