package sqlbuilder

import (
	"fmt"
	"io"
	"strings"
)

const (
	EOL   = "\n"
	Space = " "
)

func EndLine(sqlWriter io.StringWriter, withSpace bool) error {
	if withSpace {
		if n, err := sqlWriter.WriteString(Space); err != nil {
			return err
		} else if n != 1 {
			return fmt.Errorf("write endline failed")
		}
		return nil
	}
	if n, err := sqlWriter.WriteString(EOL); err != nil {
		return err
	} else if n != 1 {
		return fmt.Errorf("write endline failed")
	}
	return nil
}

func WriteSpace(sqlWriter io.StringWriter, level int) error {
	if level == -1 {
		return nil
	}
	space := GetSpace(level)
	if n, err := sqlWriter.WriteString(space); err != nil {
		return err
	} else if n != len(space) {
		return fmt.Errorf("write space(level:%d) failed", level)
	}
	return nil
}

func GetSpace(level int) string {
	if level < 0 {
		return ""
	}
	if level >= 0 && level < 32 {
		return spaces[level]
	}
	return strings.Repeat("  ", level)
}

var spaces = [32]string{
	"",
	"  ",
	"    ",
	"      ",
	"        ",
	"          ",
	"            ",
	"              ",
	"                ",
	"                  ",
	"                    ",
	"                      ",
	"                        ",
	"                          ",
	"                            ",
	"                              ",
	"                                ",
	"                                  ",
	"                                    ",
	"                                      ",
	"                                        ",
	"                                          ",
	"                                            ",
	"                                              ",
	"                                                ",
	"                                                  ",
	"                                                    ",
	"                                                      ",
	"                                                        ",
	"                                                          ",
	"                                                            ",
	"                                                              ",
}
