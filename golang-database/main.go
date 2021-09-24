package main

import (
   "os"
   "io/ioutil"
   "fmt"
   "encoding/json"
   "sync"
   "path/filepath"
   "github.com/jcelliott/lumber"
)

const DATABASE_VERSION = "1.0.1"

type Logger interface {
   Fatal(string, ...interface{})
   Error(string, ...interface{})
   Warn(string, ...interface{})
   Info(string, ...interface{})
   Debug(string, ...interface{})
   Trace(string, ...interface{})
}

type Driver struct {
   mutex    sync.Mutex
   mutexes  map[string]*sync.Mutex
   dir      string
   log      Logger
}

type Options struct {
   Logger
}

type Address struct {
   City     string
   State    string
   Country  string
   Pincode  json.Number
}

type User struct {
   Name     string
   Age      json.Number
   Contact  string
   Company  string
   Address  Address
}

func New(dir string, options *Options) (*Driver, error) {
   dir = filepath.Clean(dir)
   
   opts := Options{}
   if options != nil {
      opts = *options
   }
   
   if opts.Logger == nil {
      opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
   }
   
   driver := Driver{
      dir: dir,
      mutexes: make(map[string]*sync.Mutex),
      log: opts.Logger,
   }
   
   if _, err := os.Stat(dir); err == nil {
      opts.Logger.Debug("Using '%s' database already exists!\n", dir)
      return &driver, nil
   }
   
   opts.Logger.Debug("Creating the database at '%s'\n", dir)
   return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
   if collection == "" {
      return fmt.Errorf("Collection empty! No place to save record.")
   }
   
   if resource == "" {
      return fmt.Errorf("Missing resource! unable to save record. (no name)")
   }
   
   record := filepath.Join(d.dir, collection, resource)
   if _, err := stat(record); err != nil {
      return err
   }
   
   b, err := ioutil.ReadFile(record + ".json")
   if err != nil {
      return err
   }
   
   return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
   if collection == "" {
      return nil, fmt.Errorf("No collection - unable to read!")
   }
   
   dir := filepath.Join(d.dir, collection)
   if _, err := stat(dir); err != nil {
      return nil, err
   }
   
   files, _ := ioutil.ReadDir(dir)
   var records []string
   for _, file := range files {
      b, err := ioutil.ReadFile(filepath.Join(d.dir, file.Name()))
      if err != nil {
         return nil, err
      }
      
      records = append(records, string(b))
   }
   
   return records, nil
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
   if collection == "" {
      return fmt.Errorf("Collection empty! No place to save record.")
   }
   
   if resource == "" {
      return fmt.Errorf("Missing resource! unable to save record. (no name)")
   }
   
   mutex := d.GetOrCreateMutex(collection)
   mutex.Lock()
   defer mutex.Unlock()
   
   dir := filepath.Join(d.dir, collection)
   finalPath := filepath.Join(dir, resource + ".json")
   tmpPath := finalPath + ".tmp"
   
   if err := os.MkdirAll(dir, 0755); err != nil {
      return err
   }
   
   b, err := json.MarshalIndent(v, "", "\t")
   if err != nil {
      return err
   }
   b = append(b, byte('\n'))
   
   if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
      return err
   }
   
   return os.Rename(tmpPath, finalPath)
}

func (d *Driver) Delete(collection, resource string) error {
   path := filepath.Join(collection, resource)
   
   mutex := d.GetOrCreateMutex(collection)
   mutex.Lock()
   defer mutex.Unlock()
   
   dir := filepath.Join(d.dir, path)
   
   switch fi, err := stat(dir); {
      case fi == nil, err != nil:
         return fmt.Errorf("Unable to find file or directory named %v\n", path)
      case fi.Mode().IsDir():
         return os.RemoveAll(dir)
      case fi.Mode().IsRegular():
         return os.RemoveAll(dir + ".json")
   }
   
   return nil
}

func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {
   d.mutex.Lock()
   defer d.mutex.Unlock()
   
   mutex, ok := d.mutexes[collection]
   if !ok {
      mutex = &sync.Mutex{}
      d.mutexes[collection] = mutex
   }
   
   return mutex
}

func stat(path string) (fi os.FileInfo, err error) {
   if fi, err = os.Stat(path); os.IsNotExist(err) {
      fi, err = os.Stat(path + ".json")
   }
   return
}

func main() {
   dir := "./"
   db, err := New(dir, nil)
   if err != nil {
      fmt.Println("Error:", err)
   }
   
   employees := []User{
      {
         Name: "Ayush",
         Age: "19",
         Contact: "9004903289",
         Company: "Expansion Tricks",
         Address: Address{
            City: "Mumbai",
            State: "Maharashtra",
            Country: "India",
            Pincode: "400033",
         },
      },
      
      {
         Name: "Alkesh",
         Age: "20",
         Contact: "9004267596",
         Company: "Dashboard.io",
         Address: Address{
            City: "Mumbai",
            State: "Maharashtra",
            Country: "India",
            Pincode: "400033",
         },
      },
      
      {
         Name: "John",
         Age: "19",
         Contact: "1234567890",
         Company: "Expansion Tricks",
         Address: Address{
            City: "Kochi",
            State: "Kerala",
            Country: "India",
            Pincode: "40123",
         },
      },
   }
   
   for _, user := range employees {
      db.Write("users", user.Name, User{
         Name: user.Name,
         Age: user.Age,
         Contact: user.Contact,
         Company: user.Company,
         Address: user.Address,
      })
   }
   
   records, err := db.ReadAll("users")
   if err != nil {
      fmt.Println("Error:", err)
   }
   fmt.Println(records)
   
   allUsers := []User{}
   for _, found := range records {
      employee := User{}
      if err := json.Unmarshal([]byte(found), &employee); err != nil {
         fmt.Println("Error:", err)
      }
      allUsers = append(allUsers, employee)
   }
   fmt.Println(allUsers)
   
   /*
   if err := db.Delete("users", "John"); err != nil {
      fmt.Println("Error:", err)
   }
   */
   
   /*
   if err := db.Delete("users", ""); err != nil {
      fmt.Println("Error:", err)
   }
   */
}