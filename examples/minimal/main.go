package main

import (
	"github.com/BellerophonMobile/logberry"
	"os"
	"errors"
)


func somecomputation(data interface{}) (int,error) {
	return 7,nil
}


func main() {

	// Output autogenerated build information and a hello
	logberry.Main.Build(buildmeta)

	logberry.Main.Info("Start program")


	// Construct a new component of our program
	cmplog := logberry.Main.Component("MyComponent")


// Create some structured application data
	var data = struct {
		MyString string
		MyInt int
	}{ "alpha", 9 }

	// Do some activity on that data, which may fail, within the component
	tlog := cmplog.Task("Some computation", &data)
	res,err := somecomputation(data)
	if err != nil {
		tlog.Error(err)
		os.Exit(1)
	}
	tlog.Complete(&logberry.D{"Result": res})


	// Shut down the component
	cmplog.Finalize()


	// An error has occurred out of nowhere!
	logberry.Main.Fatal("Unrecoverable error", errors.New("Arbitrary fault"))

}
