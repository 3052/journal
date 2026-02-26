package main

import (
   "bytes"
   "io/fs"
   "log"
   "os"
   "os/exec"
   "path/filepath"
)

func walk_dir(path string, _ fs.DirEntry, err error) error {
   if err != nil {
      return err
   }
   if filepath.Ext(path) != ".go" {
      return nil
   }
   log.Print(path)
   data, err := os.ReadFile(path)
   if err != nil {
      return err
   }
   data = bytes.ReplaceAll(
      data, []byte{'\t'}, []byte("   "),
   )
   return os.WriteFile(path, data, os.ModePerm)
}

func run(name string, arg ...string) error {
   command := exec.Command(name, arg...)
   log.Println(command.Args)
   return command.Run()
}

func main() {
   log.SetFlags(log.Ltime)
   err := run("gofmt", "-w", ".")
   if err != nil {
      log.Fatal(err)
   }
   err = filepath.WalkDir(".", walk_dir)
   if err != nil {
      log.Fatal(err)
   }
}
