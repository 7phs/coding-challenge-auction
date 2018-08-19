package helpers

import "bytes"

type ErrList []error

func (o *ErrList) Add(err ...error) {
	if len(err) > 0 {
		*o = append(*o, err...)
	}
}

func (o ErrList) Error() string {
	if len(o) == 0 {
		return ""
	}

	buf := bytes.NewBufferString("")
	for i, err := range o {
		if i > 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(err.Error())
	}

	return buf.String()
}
