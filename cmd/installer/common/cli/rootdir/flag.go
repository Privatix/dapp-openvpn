package rootdir

import (
	"flag"
	"fmt"

	"gopkg.in/reform.v1"
)

const (
	name = "rootdir"
)

// FlagRootDir is a object that is used to process the rootdir flag.
type FlagRootDir struct {
	run func(flag *string, customization *bool,
		db *reform.DB, agent bool) error
	val   *string
	val2  *bool
	db    *reform.DB
	agent bool
}

// NewFlagRootDir initializes the object to process the rootdir flag.
func NewFlagRootDir() *FlagRootDir {
	return &FlagRootDir{
		run: processor,
		val: flag.String("rootdir", "",
			"Full path to root directory of service adapter"),
		val2: flag.Bool("setauth", false,
			"Generate authentication credentials "+
				"for service adapter"),
	}
}

// Name returns rootdir flag name.
func (v *FlagRootDir) Name() string {
	return name
}

// Value returns value of rootdir flag.
func (v *FlagRootDir) Value() interface{} {
	return v.val
}

// Process performs processing of the rootdir flag.
func (v *FlagRootDir) Process() error {
	if v.db == nil {
		return fmt.Errorf("database is null")
	}
	return v.run(v.val, v.val2, v.db, v.agent)
}

// AdditionalParams adds new parameters for processing rootdir flag.
func (v *FlagRootDir) AdditionalParams(db *reform.DB, agent *bool) error {
	v.db = db
	v.agent = *agent
	return nil
}
