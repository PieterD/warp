// Code generated by "stringer -type=Type"; DO NOT EDIT.

package gl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Float-0]
	_ = x[Vec2-1]
	_ = x[Vec3-2]
	_ = x[Vec4-3]
	_ = x[Mat2-4]
	_ = x[Mat3-5]
	_ = x[Mat4-6]
	_ = x[Byte-7]
	_ = x[UnsignedByte-8]
	_ = x[Short-9]
	_ = x[UnsignedShort-10]
	_ = x[Int-11]
	_ = x[UnsignedInt-12]
}

const _Type_name = "FloatVec2Vec3Vec4Mat2Mat3Mat4ByteUnsignedByteShortUnsignedShortIntUnsignedInt"

var _Type_index = [...]uint8{0, 5, 9, 13, 17, 21, 25, 29, 33, 45, 50, 63, 66, 77}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
