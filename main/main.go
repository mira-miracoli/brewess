package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mira-miracoli/brewess/internal/model"
	"github.com/objectbox/objectbox-go/objectbox"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func main() {
	// load objectbox
	ob := initObjectBox()
	defer ob.Close()

	box := model.BoxForRecipe(ob)

	runInteractiveShell(box)
}

func runInteractiveShell(box *model.RecipeBox) {
	// our simple interactive shell
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to brewess, the smart brewing recipe management")
	printHelp()

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		//input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSpace(input)
		args := strings.Fields(input)

		switch strings.ToLower(args[0]) {
		case "new":
			if len(args) != 1 {
				fmt.Fprintf(os.Stderr, "wrong number of arguments, expecting exactly one\n")
			} else {
				createRecipe(box)
			}
		case "delete":
			if len(args) != 2 {
				fmt.Fprintf(os.Stderr, "wrong number of arguments, expecting exactly one\n")
			} else if id, err := strconv.ParseUint(args[1], 10, 64); err != nil {
				fmt.Fprintf(os.Stderr, "could not parse ID: %s\n", err)
			} else {
				delRecipe(box, id)
			}
		case "ls":
			if len(args) < 2 {
				printList(box, false)
			} else if args[1] == "-a" {
				printList(box, true)
			} else {
				fmt.Fprintf(os.Stderr, "unknown argument %s\n", args[1])
				fmt.Println()
			}
		case "exit":
			return
		case "help":
			printHelp()
		default:
			fmt.Fprintf(os.Stderr, "unknown command %s\n", input)
			printHelp()
		}
	}
}

func initObjectBox() *objectbox.ObjectBox {
	objectBox, err := objectbox.NewBuilder().Model(model.ObjectBoxModel()).Build()
	if err != nil {
		panic(err)
	}
	return objectBox
}

func ScanOrErrorString() (text string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := strconv.Atoi(scanner.Text()); err == nil {
			fmt.Printf("This looks like a number, please try again\n")
		} else {
			break
		}
	}
	return scanner.Text()
}

func ScanOrErrorNumber() (f32 float32) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if f64, err := strconv.ParseFloat(scanner.Text(), 32); err == nil {
			var f32 = float32(f64)
			return f32
		} else {
			fmt.Printf("This looks like a number, please try again\n")

		}
	}
	return f32
}

func FindMalt(box *model.RecipeBox)

// create Recipe how to parse/ insert the fields and check for errors?
func createRecipe(box *model.RecipeBox) {
	fmt.Printf("Let's get started with your recipe.\nFirst, enter a name for it.\n")
	name := ScanOrErrorString()

	fmt.Printf("Well done! Now enter a short description\n")
	description := ScanOrErrorString()

	fmt.Printf("Now enter the malt(s) you want to use and it's proportion in the bulk.\n No Worries - if not found in the library, we'll create a new entry together\n")
	malts = FindMalts()
	recipe := &model.Recipe{
		Name:        name,
		Description: description,
		Malts:       malts,
		DateCreated: obNow(),
	}

	if id, err := box.Put(recipe); err != nil {
		fmt.Fprintf(os.Stderr, "could not create task: %s\n", err)
	} else {
		task.Id = id
		fmt.Printf("task ID %d successfully created\n", task.Id)
	}
}
