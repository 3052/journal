package main

import (
   "fmt"
   "log"
   "os"
   "os/exec"
   "strings"
   "text/template"
   "time"
)

func (g *git_board) New() error {
   cmd := exec.Command("git", "add", ".")
   fmt.Println(cmd.Args)
   err := cmd.Run()
   if err != nil {
      return err
   }
   cmd = exec.Command("git", "diff", "--cached", "--numstat")
   fmt.Println(cmd.Args)
   text, err := cmd.Output()
   if err != nil {
      return err
   }
   // split fails on empty string
   for _, line := range strings.FieldsFunc(string(text), lines) {
      var add, del int
      // binary files will be "- - hello.txt", so ignore error
      fmt.Sscan(line, &add, &del)
      // Add
      g.Add += add
      // Delete
      g.Delete += del
      // Change
      g.Change++
   }
   g.Target = 100
   if g.Add >= g.Target {
      g.AddStatus = pass
   } else {
      g.AddStatus = fail
   }
   if g.Delete >= g.Target {
      g.DeleteStatus = pass
   } else {
      g.DeleteStatus = fail
   }
   if g.Change >= g.Target {
      g.ChangeStatus = pass
   } else {
      g.ChangeStatus = fail
   }
   // Then
   then, err := get_then()
   if err != nil {
      return err
   }
   g.Then = then
   // Now
   g.Now = time.Now().AddDate(0, 0, -1).String()[:10]
   if g.Then <= g.Now {
      g.DateStatus = pass
   } else {
      g.DateStatus = fail
   }
   return nil
}

type git_board struct {
   Add int
   AddStatus string
   Delete int
   DeleteStatus string
   Change int
   ChangeStatus string
   Target int
   Then string
   Now string
   DateStatus string
}

func main() {
   var board git_board
   err := board.New()
   if err != nil {
      log.Fatal(err)
   }
   temp, err := new(template.Template).Parse(format)
   if err != nil {
      log.Fatal(err)
   }
   if err := temp.Execute(os.Stdout, board); err != nil {
      log.Fatal(err)
   }
}

func get_then() (string, error) {
   cmd := exec.Command("git", "log", "-1", "--format=%cI")
   fmt.Println(cmd.Args)
   text, err := cmd.Output()
   if err != nil {
      return "", err
   }
   if len(text) >= 11 {
      text = text[:10]
   }
   return string(text), nil
}

const (
   fail = "\x1b[30;101m Fail \x1b[m"
   pass = "\x1b[30;102m Pass \x1b[m"
)

const format =
   "{{ .AddStatus }} additions\ttarget:{{ .Target }}\tactual:{{ .Add }}\n" +
   "{{ .DeleteStatus }} deletions\ttarget:{{ .Target }}\tactual:{{ .Delete }}\n" +
   "{{ .ChangeStatus }} changed files\ttarget:{{ .Target}}\tactual:{{ .Change }}\n" +
   "{{ .DateStatus }} last commit\ttarget:{{ .Then }}\tactual:{{ .Now }}\n"

func lines(r rune) bool {
   return r == '\n'
}
