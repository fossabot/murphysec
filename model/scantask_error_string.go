// Code generated by "stringer -type scanTaskError -linecomment -output scantask_error_string.go"; DO NOT EDIT.

package model

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrPathIsNotAbsolute-1]
	_ = x[ErrMustBeDirectory-2]
	_ = x[ErrMustBeFile-3]
	_ = x[ErrSubtaskNameIsEmpty-4]
}

const _scanTaskError_name = "task: path must be absolutetask: path must be a directorytask: path must be a filetask: subtask name is empty"

var _scanTaskError_index = [...]uint8{0, 27, 57, 82, 109}

func (i scanTaskError) String() string {
	i -= 1
	if i < 0 || i >= scanTaskError(len(_scanTaskError_index)-1) {
		return "scanTaskError(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _scanTaskError_name[_scanTaskError_index[i]:_scanTaskError_index[i+1]]
}