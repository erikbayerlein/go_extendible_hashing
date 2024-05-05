package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"extendible-hashing.com/src"
)

func main() {
  file, err := os.Open("io/in.txt")

  if err != nil {
    fmt.Println(err)
    return
  }

  defer file.Close()

  output, err := os.Create("io/out.txt")

  if err != nil {
    fmt.Println(err)
    return
  }

  scanner := bufio.NewScanner(file)

  is_first_line := true
  var directory *models.Directory

  for scanner.Scan() {
    line := scanner.Text()

    if is_first_line {
      line_parsed := strings.Split(line, "/")

      global_depth, _ := strconv.Atoi(line_parsed[1]) 
      directory = models.CreateDirectory(global_depth) 
      is_first_line = false
    } else {
      line_parsed := strings.Split(line, ":") 

      command := line_parsed[0]
      key, _ := strconv.Atoi(line_parsed[1])

      switch command {
      case "INC":
        answer := directory.Insert(key)
        fmt.Fprintln(output, "INC=" + strconv.Itoa(key) + "," + strconv.Itoa(answer.GlobalDepth) + "," + strconv.Itoa(answer.LocalDepth))
        if answer.Duplicated {
          fmt.Fprintln(output, "DUP_DIR:" + strconv.Itoa(answer.GlobalDepth) + "," + strconv.Itoa(answer.LocalDepth))
        }
      case "BUS":
        answer := directory.Search(key)
        fmt.Fprintln(output, "BUS=" + answer)
      case "REM":
        answer := directory.Remove(key)
        fmt.Fprintln(output, "REM=" + answer[0] + "," + answer[1] + "," + answer[2])
      default:
        fmt.Println("Invalid command")
        output.Close()
      }
    }
  }

  output.Close()
}

