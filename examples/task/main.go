package main

import (
	"github.com/BellerophonMobile/logberry"
	"io/ioutil"
	"net/http"
)

func geticon(proc logberry.Context) error {

	url := "https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/logberry.png"

	get := proc.Task("Download strawberry icon").Resource(url)
	res, err := http.Get(url)
	if err != nil {
		return get.Error(err)
	} else if res.StatusCode != http.StatusOK {
		return get.Failure(http.StatusText(res.StatusCode))
	}
	get.Success()

	return nil

}

func main() {

	var value = struct {
		StringField string
		IntField    int
	}{
		StringField: "Banana",
		IntField:    24,
	}

	processor := logberry.Main.Task("Some data task", value).Begin()

	processor.Info("Prepare some data")

	myfilename := "/home/nouser/doesnotexist"
	read := processor.Task("Read app data").Resource(myfilename)
	if _, err := ioutil.ReadFile(myfilename); err != nil {
		read.Error(err)
	} else {
		read.Success()
	}

	if e := geticon(processor); e != nil {
		logberry.Main.Error("Could not get icon", e)
	}

	compute := processor.Task("Compute results")
	compute.Success()

	processor.AddData("Throughput", 23.0/100.0)

	processor.Success()

}
