package testutil

func ToPointerInt32(value int) *int32 {
	i := int32(value)

	return &i
}
