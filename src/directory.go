package models

import (
  "strconv"
  "fmt"
)

type Directory struct {
  global_depth    int
  directory_lines []DirectoryLine
}

type DirectoryLine struct {
  index       string
  local_depth int
  bucket      *Bucket
}

type Bucket struct {
  name string
  data []int
}

type InsertAnswer struct {
  Duplicated  bool
  GlobalDepth int
  LocalDepth  int
}


func CreateDirectory(global_depth int) *Directory {
  directory_lines := make([]DirectoryLine, 0)

  directory := &Directory {
    global_depth: global_depth,
    directory_lines: directory_lines,
  }

  binary_numbers := GenerateBinaryNumbers(global_depth)

  for _, binary := range binary_numbers {
    bucket := &Bucket {name: binary}
    line := &DirectoryLine {index: binary, local_depth: global_depth, bucket: bucket}
    directory.directory_lines = append(directory.directory_lines, *line)
  }

  return directory
}

func (d *Directory) SearchByIndex(index string) *DirectoryLine {
  for _, line := range d.directory_lines {
    if line.index == index {
      return &line
    }
  }

  return nil
}

func (d* Directory) Search(key int) string {
  bucket_index := Hasher(key, d.global_depth)

  lines := make([]DirectoryLine, 0)

  for _, line := range d.directory_lines {
    if line.index == bucket_index {
      lines = append(lines, line)
    }
  }

  num_of_tuples_found := 0

  for _, line := range lines {
    bucket := line.bucket

    if bucket != nil {
      for _, year := range bucket.data {
        if year == key {
          num_of_tuples_found++
          fmt.Printf("Key found in bucket %s\n", bucket.name)
        } else {
          fmt.Printf("Key not found in bucket %s\n", bucket.name)
        }
      } 
    } else {
      fmt.Println("Bucket does not exist")
    }
  }

  return strconv.Itoa(num_of_tuples_found)
} 


func (d *Directory) Insert(key int) *InsertAnswer {
  bucketIndex := Hasher(key, d.global_depth)

  answer := &InsertAnswer{false, 0, 0}

  var lines []DirectoryLine

  for _, directoryLine := range d.directory_lines {
    if directoryLine.index == bucketIndex {
        lines = append(lines, directoryLine)
      }
  }

  bucket := lines[0].bucket

  if len(bucket.data) <= 2 { // bucket isn't full
    bucket.data = append(bucket.data, key)
    fmt.Println("Key inserted in bucket " + bucket.name)
  } else { // bucket is full
    var newBucket Bucket

    if len(lines) > 1 { // local depth < global depth
      if len(bucket.data) == 2 { // bucket is full
        lines[1].bucket = &newBucket
        d.distributeBucket(lines[0], lines[1], d.global_depth, key)
        fmt.Println("Key inserted in bucket " + lines[1].bucket.name)
      } else { // bucket isn't full
        bucket.data = append(bucket.data, key)
        fmt.Println("Key inserted in bucket " + bucket.name)
      }
    } else { // local depth == global depth
      if len(bucket.data) < 2 { // bucket isn't full
        bucket.data = append(bucket.data, key)
        fmt.Println("Key inserted in bucket " + bucket.name)
      } else { // bucket is full
        newIndex := "1" + lines[0].bucket.name
        d.duplicateDirectory()
        line := d.SearchByIndex(newIndex)
        d.distributeBucket(lines[0], *line, d.global_depth, key)
        answer.LocalDepth = line.local_depth + 1
        answer.Duplicated = true
      }
    }
  }

  answer.GlobalDepth = d.global_depth
  if answer.LocalDepth == 0 {
    answer.LocalDepth = lines[0].local_depth
  }

  return answer
}

func (d *Directory) duplicateDirectory() {
  d.global_depth++
  directorySize := len(d.directory_lines)
  for i:=0; i<directorySize; i++ {
    line := d.directory_lines[i]
    var newLine DirectoryLine

    newLine.index = "1" + line.index
    line.index = "0" + line.index
    newLine.bucket = line.bucket

    newLine.local_depth = line.local_depth
    d.directory_lines = append(d.directory_lines, newLine)
  }

  fmt.Println("Directory duplicated")
}

func (d *Directory) distributeBucket(oldLine DirectoryLine, newLine DirectoryLine, depth int, newKey int) {
  auxData := oldLine.bucket.data
  auxData = append(auxData, newKey)

  var bucket Bucket

  bucket.name = oldLine.bucket.name + "k"
  newLine.bucket = &bucket

  oldLine.bucket.data = oldLine.bucket.data[:0]

  var keys []int
  
  for _, key := range auxData {
    keys = append(keys, key)
  }

  yearsWithoutDuplicates := removeDuplicates(keys)

  for _, year := range yearsWithoutDuplicates {
    d.Insert(year)
  }

  oldLine.local_depth = depth
  newLine.local_depth = depth
}

func (d *Directory) Remove(key int) []string {
  bucketIndex := Hasher(key, d.global_depth)

  var directoryLine *DirectoryLine
  
  for _, line := range d.directory_lines {
    if bucketIndex == line.index {
      directoryLine = &line
      break
    }
  }

  var tuplesRemoved int
  answer := make([]string, 3)

  if directoryLine != nil {
    bucket := directoryLine.bucket

    for i, year := range bucket.data {
      if year == key {
        bucket.data = RemoveIndex(bucket.data, i)
        tuplesRemoved++

        fmt.Println("Key removed from bucket " + bucket.name)
      }
    }
  }

  answer[0] = strconv.Itoa(tuplesRemoved)
  answer[1] = strconv.Itoa(d.global_depth)
  answer[2] = strconv.Itoa(directoryLine.local_depth)

  return answer
}

func RemoveIndex(s []int, index int) []int {
    return append(s[:index], s[index+1:]...)
}

func removeDuplicates(arr []int) []int {
    encountered := map[int]bool{}
    result := []int{}

    for _, v := range arr {
        if encountered[v] == false {
            encountered[v] = true
            result = append(result, v)
        }
    }

    return result
}

func Hasher(key int, depth int) string {
  binary_string := strconv.FormatInt(int64(key), 2)

  if depth >= len(binary_string) {
    return binary_string
  }

  return binary_string[len(binary_string)-depth:]  
}

func GenerateBinaryNumbers(global_depth int) []string {
    binary_numbers := make([]string, 0)

    for i := 0; i < 1<<global_depth; i++ {
        binary := fmt.Sprintf("%0*b", global_depth, i)
        binary_numbers = append(binary_numbers, binary)
    }

    return binary_numbers
}

