package main 
 
import ( 
    "fmt" 
    "log" 
    "log/syslog" 
    "os" 
    "path/filepath" 
) 
 
func main() { 
 
    programName := filepath.Base(os.Args[0]) 
    sysLog, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, programName)

    if err != nil { 
        log.Fatal(err) 
    } else { 
        log.SetOutput(sysLog) 
    } 
    log.Println("LOG_INFO + LOG_LOCAL6: Logging in Go!")
    
    sysLog, err = syslog.New(syslog.LOG_ERR|syslog.LOG_LOCAL6, "Some program!") 
    if err != nil { 
        log.Fatal(err) 
    } else { 
        log.SetOutput(sysLog) 
    } 
 
    log.Println("LOG_ERR + LOG_LOCAL6: Logging in Go!") 
    fmt.Println("Will you see this?") 
} 
