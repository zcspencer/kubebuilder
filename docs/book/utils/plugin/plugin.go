package plugin

import (
	"fmt"
	"io"
	"encoding/json"
	"os"
	"io/ioutil"
)

// Plugin represents a mdBook plugin.
type Plugin interface {
	// SupportsOutput checks if the given plugin supports the given output format.
	SupportsOutput(string) bool
	// Process modifies the book in the input, which gets returned as the result of the plugin.
	Process(*Input) error
}

// Run runs the given plugin on the given input stream, outputting its result to the given
// result, assuming the given command-line args (without program name).
func Run(plug Plugin, inputRaw io.Reader, outputRaw io.Writer, args ...string) error {
	if len(args) > 1 && args[0] == "supports" {	
		// we support any renderer, no need to check (name is in Args[1])
		if plug.SupportsOutput(args[1]) {
			return nil
		}
		return fmt.Errorf("output format %q not supported", args[1])
	}

	var input Input
	dec := json.NewDecoder(inputRaw)
	if err := dec.Decode(&input); err != nil {
		return fmt.Errorf("unable to decode preprocessor input: %v", err)
	}

	if err := plug.Process(&input); err != nil {
		return err
	}

	out, err := json.Marshal(&input.Book)
	if err != nil {
		return fmt.Errorf("unable to encode output book object: %v", err)
	}

	if n, err := outputRaw.Write(out); err != nil || n < len(out) {
		if err == nil && n < len(out) {
			err = io.ErrShortWrite
		}
		return fmt.Errorf("unable to write output book object: %v", err)
	}

	ioutil.WriteFile("/tmp/litout.json", out, os.ModePerm)

	return nil
}
