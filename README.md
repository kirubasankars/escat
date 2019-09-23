# escat
command line tool for managing elasticsearch


Usage :

  escat -help
  
  ```

Usage : escat [OPTIONS] COMMAND:

Options:
  -d	set to debug (true|false)
  -f string
    	set the fields
  -host string
    	set the elasticsearch host url or use environment variable ES_HOST
  -json
    	set the output format to json
  -password string
    	set the elasticsearch password or or use environment variable ES_PASS
  -pretty
    	pretty print (true|false) (default true)
  -s string
    	set the fields to sort
  -user string
    	set the elasticsearch user or use environment variable ES_USER (default "elastic")
  -v	set header (true|false)

Commands:
   health            Print Cluster health
   snapshots         Print Snapshots     
   allocation        Print Allocation    
   nodes             Print Nodes         
   indices           Print Indices       
   segments          Print Segments      
   master            Print Master        
   aliases           Print Alias         
   repositories      Print Repositories  
   count             Print Count         
   plugins           Print Plugins       
   templates         Print Templates     
   info              Print Info          
   user              Print User          
   role              Print Role        

```


Setup Environment Variables
```
  export ES_HOST=http://locahost:9200
  export ES_USER=elastic
  export ES_PASS=changeme
```

To Print Indices, You can just type copule of letters in the begining of the command

```
./escat indices 
or 
./escat i 
./escat i a*
```

To Print Index Mapping, add "_" in the end.

```
./escat i _
./escat i a* _
```

Note: "_" only works indices,templates and snapshots.
