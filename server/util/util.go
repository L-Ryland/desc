package util

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type Error struct {
	err    string
	code   *codes.Code
	subErr *Error
}

func (e *Error) Error() string {
	err := e.err
	if e.code != nil {
		err = fmt.Sprintf("%s:%s.", err, e.code.String())
	}
	if e.subErr != nil {
		err = fmt.Sprintf("%sWith cause:%s", err, e.subErr.Error())
	}
	return err
}

func (e *Error) WithCause(err error) *Error {
	if err == nil {
		return e
	}
	eStruct, ok := err.(*Error)
	if !ok {
		eStruct = &Error{err: err.Error()}
	}
	e.subErr = eStruct
	return e
}

func (e *Error) WithCode(code codes.Code) *Error {
	e.code = &code
	return e
}

func (e *Error) Log() *Error {
	logrus.Error(e.Error())
	return e
}

func Errorf(format string, a ...interface{}) *Error {
	return &Error{
		err: fmt.Sprintf(format, a...),
	}
}

func HaveErrorCode(e error, code codes.Code) bool {
	err, ok := e.(*Error)
	if !ok {
		return false
	}
	if err == nil {
		return false
	}
	if err.code != nil && err.code.String() == code.String() {
		return true
	} else {
		return HaveErrorCode(err.subErr, code)
	}
}
func DecodeJson(j string, item interface{}) {
	if err := json.Unmarshal([]byte(j), item); err != nil {
		panic(Errorf("decode json %s failed", j).WithCause(err))
	}
}
func EncodeJson(item interface{}) string {
	j, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		panic(Errorf("encode json %s failed", j).WithCause(err))
	}
	return string(j)
}
