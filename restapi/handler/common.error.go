package handler

import (
	"bytes"
	"fmt"
)

type RespErrors interface {
	AddError(int, ...interface{})
	AddErrorf(int, string, ...interface{})
}

type ErrorRecord struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

type ErrorRecordList []*ErrorRecord

func (o *ErrorRecordList) addError(id int, desc string) {
	*o = append(*o, &ErrorRecord{Id: id, Desc: desc})
}

func (o *ErrorRecordList) AddError(id int, v ...interface{}) {
	o.addError(id, fmt.Sprint(v...))
}

func (o *ErrorRecordList) AddErrorf(id int, format string, v ...interface{}) {
	o.addError(id, fmt.Sprintf(format, v...))
}

func (o *ErrorRecordList) Append(errList ErrorRecordList) {
	*o = append(*o, errList...)
}

func (o *ErrorRecordList) AppendError(id int, errList []error) {
	for _, err := range errList {
		o.AddError(id, err.Error())
	}
}

func (o ErrorRecordList) Error() string {
	if len(o) == 0 {
		return ""
	}

	buf := bytes.NewBufferString("")

	for i, err := range o {
		if i > 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(fmt.Sprintf("[%d] %s", err.Id, err.Desc))
	}

	return buf.String()
}

type RespError struct {
	Errors ErrorRecordList `json:"error"`
}

func (o *RespError) AddError(id int, v ...interface{}) {
	o.Errors.AddError(id, v...)
}

func (o *RespError) AddErrorf(id int, format string, v ...interface{}) {
	o.Errors.AddErrorf(id, format, v...)
}

func (o *RespError) AppendError(errList ErrorRecordList) {
	o.Errors.Append(errList)
}

type ListOfErr struct {
	Id int

	err ErrorRecordList
}

func (o *ListOfErr) Check(errCondition bool, errMessage string) bool {
	if errCondition {
		o.err.AddError(o.Id, errMessage)
	}

	return errCondition
}

func (o *ListOfErr) Checkf(errCondition bool, format string, v ...interface{}) bool {
	if errCondition {
		o.err.AddErrorf(o.Id, format, v...)
	}

	return errCondition
}

func (o *ListOfErr) HasError() bool {
	return o != nil && len(o.err) > 0
}

func (o *ListOfErr) Error() string {
	return o.err.Error()
}
