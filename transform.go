package stl

// This file contains generic 3D transformation stuff

import (
	"math"
)

// RotationMatrix calculates a 4x4 rotation matrix for a rotation of Angle in radians
// around a rotation axis defined by a point on it (pos) and its direction (dir).
// The result is written into *rotationMatrix.
func RotationMatrix(pos Vec3, dir Vec3, angle float64, rotationMatrix *Mat4) {
	dirN := dir.UnitVec3()

	s := math.Sin(angle)
	c := math.Cos(angle)

	u := float64(dirN[0])
	v := float64(dirN[1])
	w := float64(dirN[2])

	su := s * u
	sv := s * v
	sw := s * w

	ic := 1 - c

	icu := ic * u

	icuu := icu * u
	icuv := icu * v
	icuw := icu * w

	icv := ic * v

	icvv := icv * v
	icvw := icv * w

	icw := ic * w
	icww := icw * w

	mRotate := Mat4{
		Vec4{icuu + c, icuv - sw, icuw + sv, 0},
		Vec4{icuv + sw, icvv + c, icvw - su, 0},
		Vec4{icuw - sv, icvw + su, icww + c, 0},
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
