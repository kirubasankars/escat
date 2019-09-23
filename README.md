# escat
command line tool for managing elasticsearch


Usage :

  escat -help


Setup Environment Variables
  export ES_HOST=http://locahost:9200
  export ES_USER=elastic
  export ES_PASS=changeme

To Print Indices

./escat i 
./escat i a*

To Print Index Mapping, add "_" in the end.

./escat i _
./eacat i a* _

Note: "_" only works indices,templates and snapshots.
