package main

import (
	// go 
	"fmt"
	"os"
	// own
	cfd "cloudformationdeploy"
	// utils
	"github.com/thatisuday/clapper"
)

func main() {

	const cmdDestroyString = "destroy"
	const cmdDeployString = "deploy"
	const cmdStatusString = "status"
	const cmdShowString = "show"
	const cmdHelpString = "help"

	// Look for commands
	registry := clapper.NewRegistry()
	registry.Register(cmdDeployString)
	registry.Register(cmdDestroyString)
	registry.Register(cmdStatusString)
	registry.Register(cmdShowString)
	registry.Register(cmdHelpString)
	
	// parse command-line arguments
	command, err := registry.Parse(os.Args[1:])

	// check for command line error
	if err != nil {
		fmt.Printf("error => %#v\n", err)
		help()
		return
	}

	// no command
	if( len(command.Name) == 0 ) {
		help()
		os.Exit(1)
	}
	cmd := command.Name;
	
	// create aws config & CloudFormation client
	client := cfd.Client()
	
	const stackname = "demotemplate"
	template := cfd.CreateTemplate(stackname);
	

	if cmd == cmdDeployString {
		// Create a new CloudFormation template	
		cfd.CreateStack(client,stackname, template)
		cfd.ShowStatus(client,stackname,template,cfd.StatusCreateComplete);
	}
	
	if cmd == cmdDestroyString {
		cfd.DeleteStack(client,stackname)
		cfd.ShowStatus(client,stackname,template,cfd.StatusDeleteComplete);
	}
	
	if cmd == cmdStatusString {
		cfd.ShowStatus(client,stackname,template,cfd.StatusCreateComplete);
	}

	if cmd == cmdShowString {
		y,_ := template.YAML()
		fmt.Println(string(y));
	}


	if cmd == cmdHelpString {
		help();
		os.Exit(0);
	}

}

func help(){
	fmt.Println("CloudFormation deploy app. ")
	fmt.Println("CloudFormation Template is generated automatically.")
	fmt.Println("Please call with [deploy|destroy|status|show|help] . ")
}