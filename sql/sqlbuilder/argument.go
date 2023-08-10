package sqlbuilder

type Arg any

type ArgWriter interface {
	WriteArg(arg Arg) error
}

func WriteArgs(writer ArgWriter, args ...Arg) error {
	for _, arg := range args {
		if err := writer.WriteArg(arg); err != nil {
			return err
		}
	}
	return nil
}
