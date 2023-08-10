package sqlbuilder

import (
	"fmt"
	"io"
)

const (
	EOL = "\n"
)

func EndLine(sqlWriter io.StringWriter) error {
	if n, err := sqlWriter.WriteString(EOL); err != nil {
		return err
	} else if n != 1 {
		return fmt.Errorf("write endline failed")
	}
	return nil
}

func WriteSpace(sqlWriter io.StringWriter, level int) error {
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
	return GetSpace(level-31) + spaces[31]
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
