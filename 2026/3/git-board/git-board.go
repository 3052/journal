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

type GitBoard struct {
   Add          int
   AddStatus    string
   Delete       int
   DeleteStatus string
   Change       int
   ChangeStatus string
   Target       int
   Then         string
   Now          string
   DateStatus   string
}

func GenerateGitBoard() (*GitBoard, error) {
   g := &GitBoard{}

   command := exec.Command("git", "add", ".")
   fmt.Println(command.Args)
   if err := command.Run(); err != nil {
      return nil, err
   }

   command = exec.Command("git", "diff", "--cached", "--numstat")
   fmt.Println(command.Args)
   data, err := command.Output()
   if err != nil {
      return nil, err
   }

   lines := strings.FieldsFunc(string(data), func(r rune) bool {
      return r == '\n'
   })
   for _, line := range lines {
      var add, del int
      fmt.Sscan(line, &add, &del)
      g.Add += add
      g.Delete += del
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

   // Renamed get_then to getLastCommitDate
   then, err := getLastCommitDate()
   if err != nil {
      return nil, err
   }
   g.Then = then

   g.Now = time.Now().AddDate(0, 0, -1).String()[:10]
   if g.Then <= g.Now {
      g.DateStatus = pass
   } else {
      g.DateStatus = fail
   }

   return g, nil
}

func main() {
   // Updated function call
   board, err := GenerateGitBoard()
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

func getLastCommitDate() (string, error) {
   command := exec.Command("git", "log", "-1", "--format=%cI")
   fmt.Println(command.Args)
   data, err := command.Output()
   if err != nil {
      return "", err
   }
   if len(data) >= 11 {
      data = data[:10]
   }
   return string(data), nil
}

const (
   fail = "\x1b[30;101m Fail \x1b[m"
   pass = "\x1b[30;102m Pass \x1b[m"
)

const format = "{{ .AddStatus }} additions\ttarget:{{ .Target }}\tactual:{{ .Add }}\n" +
   "{{ .DeleteStatus }} deletions\ttarget:{{ .Target }}\tactual:{{ .Delete }}\n" +
   "{{ .ChangeStatus }} changed files\ttarget:{{ .Target}}\tactual:{{ .Change }}\n" +
   "{{ .DateStatus }} last commit\ttarget:{{ .Then }}\tactual:{{ .Now }}\n"
