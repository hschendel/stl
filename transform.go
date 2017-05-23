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
	sinA := math.Sin(angle)
	sinAd0 := sinA * float64(dirN[0])
	sinAd1 := sinA * float64(dirN[1])
	sinAd2 := sinA * float64(dirN[2])
	cosA := math.Cos(angle)
	icosA := 1 - cosA
	icosAd0 := icosA * float64(dirN[0])
	icosAd00 := icosAd0 * float64(dirN[0])
	icosAd01 := icosAd0 * float64(dirN[1])
	icosAd02 := icosAd0 * float64(dirN[2])
	icosAd1 := icosA * float64(dirN[1])
	icosAd11 := icosAd1 * float64(dirN[1])
	icosAd12 := icosAd1 * float64(dirN[2])
	icosAd2 := icosA * float64(dirN[2])
	icosAd22 := icosAd2 * float64(dirN[2])

	mRotate := Mat4{
		Vec4{icosAd00 + cosA, icosAd01 - sinAd2, icosAd02 + sinAd1, 0},
		Vec4{icosAd01 + sinAd2, icosAd11 + cosA, icosAd12 - sinAd0, 0},
		Vec4{icosAd02 - sinAd1, icosAd12 + sinAd0, icosAd22 + cosA, 0},
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
