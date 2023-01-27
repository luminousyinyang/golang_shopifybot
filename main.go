package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type ProfileCommand struct {
	fs *flag.FlagSet
	name bool
}

func (g *ProfileCommand) Name() string {
	return g.fs.Name()
}


func (g *ProfileCommand) Run(subcommand string) error {
	scanner := bufio.NewScanner(os.Stdin)
	switch subcommand {
	case "profile": 
	fmt.Println("Select 1 of the following options:")
	fmt.Println("1. create")
	fmt.Println("2. edit")
	fmt.Println("3. delete")


	scanner.Scan()
		switch strings.ToLower(scanner.Text()) {
		case "create":
			ProfileCreate()
			return nil
		case "1":
			ProfileCreate()
			return nil
		case "edit":
			ProfileEdit()
			return nil
		case "2":
			ProfileEdit()
			return nil
		case "delete":
			ProfileDelete()
			return nil
		case "3":
			ProfileDelete()
			return nil
		default:
			return errors.New("sub command missing or invalid")
		}	


	case "task":
	fmt.Println("Select 1 of the following options:")
	fmt.Println("1. create")
	fmt.Println("2. edit")
	fmt.Println("3. start")

	scanner.Scan()
		switch strings.ToLower(scanner.Text()) {
		case "create":
			TaskCreate()
			return nil
		case "1":
			TaskCreate()
			return nil
		case "start":
			TaskEdit()
			return nil
		case "2":
			TaskEdit()
			return nil
		case "edit":
			TaskStart()
			return nil
		case "3":
			TaskStart()
			return nil
		default:
			return errors.New("sub command missing or invalid")
		}	
	default:
		return fmt.Errorf("unknown subcommand: [%s]", os.Args[1])
	}
	

}

type Runner interface {
	Run(subcommand string) error
	Name() string
}


func root(args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var subcommand string
	if len(args) < 1 {
		fmt.Println("Please enter a subcommand. [PROFILE or TASK]")
		scanner.Scan()
		subcommand = strings.ToLower(scanner.Text())
	} else if len(args) >= 2 {
		fmt.Println("only pass 1 subcommand, remove any extra spaces or characters")
		fmt.Println("Please enter a subcommand. [PROFILE or TASK]")
		scanner.Scan()
		subcommand = strings.ToLower(scanner.Text())
	} else if len(args) == 1 {
		subcommand = os.Args[1]
	}

	cmds := []Runner{
		ProfileCommands(),
		TaskCommands(),
	}

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			return cmd.Run(subcommand)
		}
	}

	return fmt.Errorf("unknown subcommand: [%s]", subcommand)
}