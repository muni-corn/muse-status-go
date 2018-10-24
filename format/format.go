package format

// Left aligns the original string to the left
func Left(original string) string {
	return "%{l}" + original;
}

// Center aligns the original string to the center
func Center(original string) string {
	return "%{c}" + original;
}

// Right aligns the original string to the right
func Right(original string) string {
	return "%{r}" + original;
}

// Dim apples a half-transparent white color to the original string
func Dim(original string) string {
	return "%{F#80FFFFFF}" + original + "%{F-}"
}
