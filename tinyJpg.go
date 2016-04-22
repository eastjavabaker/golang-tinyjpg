// Image TinyJpg auto compression v 1.0
// Author : Wahyu S - 21 April 2016 

package main

import (
     "fmt"
     "os"
     "io"
     "io/ioutil"
     "log"
     "path"
     "os/exec"
     "strings"
     "net/http"
     "path/filepath"
     "encoding/json"
)


var Config MyConfig

type MyConfig struct {
      
      TinyJpgUrl string  `json:"tinyjpg_url"`
      TinyJpgUserApi string  `json:"tinyjpg_user"`
      StartPath string `json:"start_path"` 
      TargetPath string `json:"target_path"`
      Option string  `json:"option"`   // Option 1 : Optimize , 2 : Resize, 3 : Crop           
      OptionParams Options `json:"option_params"`
      
}


type Options struct {
    
      TargetDimensions Dimensions  `json:"dimensions"`

}

type Dimensions struct {
      
      Width string  `json:"width"`
      Height string `json:"height"`      

}



func getImg(url string, filename string, copyfolder string){

    response, e := http.Get(url)
    if e != nil {
        log.Fatal(e)
    }

    defer response.Body.Close()
    
    str := filename
    i, j := strings.LastIndex(str, "/"), strings.LastIndex(str, path.Ext(str))
    name := str[i:j]+path.Ext(str)     
    //open a file for writing
    //file, err := os.Create(Config.TargetPath+"/"+name+path.Ext(str))
    file, err := os.Create(Config.TargetPath+copyfolder+name)
    if err != nil {
        log.Fatal(err)
    }
    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()

}


func main() {
	var (
		cmdOut []byte
		err    error
	)

        // pull from json data .
        file, e := ioutil.ReadFile("configs/config.json")
        if e != nil {
           fmt.Printf("File error: %v\n", e)
           os.Exit(1)
        }    
    
        json.Unmarshal(file, &Config)


	cmdName := "curl"

    searchDir := Config.StartPath

    fileList := []string{}
    err2 := filepath.Walk(searchDir, func(path string, f os.FileInfo, err2 error) error {
        fileList = append(fileList, path)
        return nil
    })

    if err2 != nil {

           fmt.Println("fail")

    }

    startIndex := strings.LastIndex(searchDir,"/")    
    fmt.Println("TinyJPG Start !!!") 
    for _, file := range fileList {
       if path.Ext(file) == ".jpg" {
        
	cmdArgs := []string{Config.TinyJpgUrl, "--user", "api:"+Config.TinyJpgUserApi, "--data-binary", "@"+file, "--dump-header", "/dev/stdout"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error command: ", err)
		os.Exit(1)
	}
	
	sha := string(cmdOut)

        lastIndex := strings.LastIndex(file,"/")
        substring := file[startIndex:lastIndex]
         
        lines := strings.Split(sha, "\n")
        image := strings.Split(lines[8], " ")
        
        os.MkdirAll(Config.TargetPath+substring, 0777)         

        getImg(strings.TrimSpace(image[1]), file, substring)  
        fmt.Println(file)
         
        }     
    }


}
