package flags

import (
	"fmt"
	"strings"
)

type StringArray []string

func (i *StringArray) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*i, ", "))
}

func (i *StringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
