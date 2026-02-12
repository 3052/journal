package main

import (
   "bufio"
   "bytes"
   "embed"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/textproto"
   "net/url"
   "os"
   "strings"
   "text/template"
)

// this is needed because http.ReadRequest is trash
func read_request(r *bufio.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(r)
   // .Method
   raw_method_path, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   method_path := strings.Fields(raw_method_path)
   req.Method = method_path[0]
   // .URL
   ref, err := url.ParseRequestURI(method_path[1])
   if err != nil {
      return nil, err
   }
   req.URL = ref
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   if req.URL.Host == "" {
      req.URL.Host = head.Get("Host")
   }
   // .Header
   req.Header = http.Header(head)
   // .Body
   data := &bytes.Buffer{}
   length, err := text.R.WriteTo(data)
   if err != nil {
      return nil, err
   }
   if length >= 1 {
      req.Body = io.NopCloser(data)
   }
   // .ContentLength
   req.ContentLength = length
   return &req, nil
}

func main() {
   var set command
   set.New()
   if set.in.name == "" {
      flag.Usage()
   } else {
      var err error
      set.out.file, err = new_file(set.out.name)
      if err != nil {
         log.Fatal(err)
      }
      defer set.out.file.Close()
      set.in.file, err = os.Open(set.in.name)
      if err != nil {
         log.Fatal(err)
      }
      defer set.in.file.Close()
      req, err := read_request(bufio.NewReader(set.in.file))
      if err != nil {
         log.Fatal(err)
      }
      if req.URL.Scheme == "" {
         if set.https {
            req.URL.Scheme = "https"
         } else {
            req.URL.Scheme = "http"
         }
      }
      if set.golang {
         err = set.write_go(req)
      } else {
         err = set.write(req)
      }
      if err != nil {
         log.Fatal(err)
      }
   }
}

func (f *command) write(req *http.Request) error {
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   if f.out.name != "" {
      // 1. body to file
      _, err = f.out.file.ReadFrom(resp.Body)
      if err != nil {
         return err
      }
      resp.ContentLength = 0
   }
   // 2. head to stdout
   return resp.Write(os.Stdout)
}

func (f *command) write_go(req *http.Request) error {
   var value request
   value.Method = req.Method
   value.URL = req.URL
   value.Header = req.Header
   if req.Body != nil {
      data, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      if f.form {
         form, err := url.ParseQuery(string(data))
         if err != nil {
            return err
         }
         value.RawBody = fmt.Sprintf("\n%#v.Encode(),\n", form)
      } else {
         value.RawBody = fmt.Sprintf("%#q", data)
      }
      value.Body = "io.NopCloser(strings.NewReader(data))"
   } else {
      value.RawBody = `""`
      value.Body = "nil"
   }
   temp, err := template.ParseFS(content, "_template.go")
   if err != nil {
      return err
   }
   return temp.Execute(f.out.file, value)
}

type command struct {
   golang bool
   https  bool
   form   bool
   in     struct {
      name string
      file *os.File
   }
   out struct {
      name string
      file *os.File
   }
}

func (f *command) New() {
   flag.BoolVar(&f.form, "f", false, "form")
   flag.BoolVar(&f.golang, "g", false, "request as Go code")
   flag.BoolVar(&f.https, "s", false, "HTTPS")
   flag.StringVar(&f.in.name, "i", "", "in file")
   flag.StringVar(&f.out.name, "o", "", "output file")
   flag.Parse()
}

type request struct {
   Method  string
   URL     *url.URL
   Header  http.Header
   Body    string
   RawBody string
}

//go:embed _template.go
var content embed.FS

func new_file(name string) (*os.File, error) {
   if name != "" {
      return os.Create(name)
   }
   return os.Stdout, nil
}
