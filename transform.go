package stl

// This file contains generic 3D transformation stuff

import (
	"math"
)

// Calculates a 4x4 rotation matrix for a rotation of angle in radians
// around a rotation axis defined by a point on it (pos) and its direction (dir).
// The result is written into *rotationMatrix.
func RotationMatrix(pos Vec3, dir Vec3, angle float64, rotationMatrix *Mat4) {
	dirN := dir.unitVec()

	s := math.Sin(angle)
	c := math.Cos(angle)

	u := float64(dirN[0])
	v := float64(dirN[1])
	w := float64(dirN[2])

	su := s * u
	sv := s * v
	sw := s * w

	ic := 1 - c

	ic_u := ic * u

	ic_uu := ic_u * u
	ic_uv := ic_u * v
	ic_uw := ic_u * w

	ic_v := ic * v

	ic_vv := ic_v * v
	ic_vw := ic_v * w

	ic_w := ic * w
	ic_ww := ic_w * w

	mRotate := Mat4{
		Vec4{ic_uu + c, ic_uv - sw, ic_uw + sv, 0},
		Vec4{ic_uv + sw, ic_vv + c, ic_vw - su, 0},
		Vec4{ic_uw - sv, ic_vw + su, ic_ww + c, 0},
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
