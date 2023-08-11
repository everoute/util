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

// End the line with EOL or Space
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

// Write indentation with space.
// If you want to add indentation when writing a string, then WriteStringWithSpace is a better choice.
func WriteSpace(sqlWriter io.StringWriter, level int) error {
	// Redundant judgments will be optimized in GetSpace.
	if CompactLevel(level) {
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

// Get the Spaces corresponding to the level.
// And custom the indentation strategy is not supported currently.
// When you want to write indentation, you should always call WriteSpace instead of GetSpace.
func GetSpace(level int) string {
	if CompactLevel(level) {
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
