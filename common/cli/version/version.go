package version

import (
	"flag"
	"fmt"
	"os"
)

const (
	undefined = "undefined"
	name      = "version"
)

// Version is a object that is used to process the version flag.
type Version struct {
	run func(flag *bool, commit, version string) error
	val *bool
	ver string
	com string
}

// NewVersionFlag initializes the object to process the version flag.
func NewVersionFlag(commit, version string) *Version {
	return &Version{
		run: func(flag *bool, commit, version string) error {
			Print(*flag, commit, version)
			return nil
		},
		val: flag.Bool("version", false, "Prints current installer version"),
		ver: version,
		com: commit,
	}
}

// Name returns flag name.
func (v *Version) Name() string {
	return name
}

// Value returns value of version flag.
func (v *Version) Value() interface{} {
	return v.val
}

// Process performs processing of the version flag.
func (v *Version) Process() error {
	return v.run(v.val, v.com, v.ver)
}

// Print prints version and completes the program.
func Print(run bool, commit, version string) {
	if run {
		fmt.Println(message(commit, version))
		os.Exit(0)
	}
}

func message(commit, version string) string {
	var c string
	var v string

	if commit == "" {
		c = undefined
	} else if len(commit) > 7 {
		c = commit[:7]
	} else {
		c = commit
	}

	if version == "" {
		v = undefined
	} else {
		v = version
	}

	return fmt.Sprintf("%s (%s)", v, c)
}
