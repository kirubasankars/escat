package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	HEALTH       = "health"
	SNAPSHOTS    = "snapshots"
	ALLOCATION   = "allocation"
	NODES        = "nodes"
	INDICES      = "indices"
	SEGMENTS     = "segments"
	MASTER       = "master"
	ALIAIS       = "aliases"
	REPOSITORIES = "repositories"
	COUNT        = "count"
	PLUGINS      = "plugins"
	TEMPLATES    = "templates"
	INFO         = "info"
	ROLE         = "role"
	USER         = "user"
)

func main() {

	var host string
	flag.StringVar(&host, "host", "", "set the elasticsearch host url or use environment variable ES_HOST")
	var user string
	flag.StringVar(&user, "user", "elastic", "set the elasticsearch user or use environment variable ES_USER")
	var password string
	flag.StringVar(&password, "password", "", "set the elasticsearch password or or use environment variable ES_PASS")

	var formatJSON bool
	flag.BoolVar(&formatJSON, "json", false, "set the output format to json")
	var pretty bool
	flag.BoolVar(&pretty, "pretty", true, "pretty print (true|false)")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "show header (true|false)")
	var fields string
	flag.StringVar(&fields, "h", "", "set the fields")
	var sortFields string
	flag.StringVar(&sortFields, "s", "", "set the fields to sort")

	var debug bool
	flag.BoolVar(&debug, "d", false, "set to debug (true|false)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "\nUsage : %s [OPTIONS] COMMAND:\n\n", "escat")

		fmt.Println("Options:")

		flag.PrintDefaults()

		fmt.Fprintf(os.Stdout, "\n")

		fmt.Println("Commands:")

		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+HEALTH, "Print Cluster health")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+SNAPSHOTS, "Print Snapshots")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+ALLOCATION, "Print Allocation")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+NODES, "Print Nodes")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+INDICES, "Print Indices")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+SEGMENTS, "Print Segments")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+MASTER, "Print Master")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+ALIAIS, "Print Alias")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+REPOSITORIES, "Print Repositories")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+COUNT, "Print Count")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+PLUGINS, "Print Plugins")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+TEMPLATES, "Print Templates")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+INFO, "Print Info")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+USER, "Print User")
		fmt.Fprintf(os.Stdout, "%-20s %-20s\n", "   "+ROLE, "Print Role")
		fmt.Fprintf(os.Stdout, "\n")
	}

	flag.Parse()

	if host == "" {
		host = os.Getenv("ES_HOST")
	}
	if user == "" {
		user = os.Getenv("ES_USER")
	}
	if password == "" {
		password = os.Getenv("ES_PASS")
	}

	if host == "" {
		log.Fatal("host URL should be specified")
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	esClient := NewElasticSearchClient(host, user, password)

	command := strings.Trim(strings.ToLower(args[0]), "")

	format := "text"
	if formatJSON {
		format = "json"
	}

	commands := []string{USER, ROLE, HEALTH, SNAPSHOTS, ALLOCATION, NODES, INDICES, SEGMENTS, MASTER, ALIAIS, REPOSITORIES, COUNT, PLUGINS, TEMPLATES, INFO}
	sort.Strings(commands)

	pickedCommand := ""
	for _, x := range commands {
		if strings.HasPrefix(x, command) {
			pickedCommand = x
			break
		}
	}

	if pickedCommand == "" {
		pickedCommand = command
	}

	if esClient.isCommandHelp(pickedCommand, args) {
		return
	}

	catr := catRequest{}
	catr.fields = fields
	catr.format = format
	catr.sortFields = sortFields
	catr.verbose = verbose
	catr.debug = debug
	if len(args) >= 2 {
		catr.arg1 = args[1]
	}
	if len(args) == 3 {
		if args[2] == "_" {
			catr.def = true
		}
	}
	if catr.arg1 == "_" {
		catr.arg1 = ""
		catr.def = true
	}
	if catr.def {
		catr.format = "json"
	}

	var (
		output     []byte
		jsonOUTPUT = false
	)
	switch {

	case pickedCommand == HEALTH:
		{
			output, jsonOUTPUT = esClient.CatHealth(catr)
		}
	case pickedCommand == SNAPSHOTS:
		{
			output, jsonOUTPUT = esClient.CatSnapshot(catr)
		}
	case pickedCommand == ALLOCATION:
		{
			output, jsonOUTPUT = esClient.CatAllocation(catr)
		}
	case pickedCommand == NODES:
		{
			output, jsonOUTPUT = esClient.CatNodes(catr)
		}
	case pickedCommand == PLUGINS:
		{
			output, jsonOUTPUT = esClient.CatPlugins(catr)
		}
	case pickedCommand == TEMPLATES:
		{
			output, jsonOUTPUT = esClient.CatTemplates(catr)
		}
	case pickedCommand == INDICES:
		{
			output, jsonOUTPUT = esClient.CatIndices(catr)
		}
	case pickedCommand == SEGMENTS:
		{
			output, jsonOUTPUT = esClient.CatSegments(catr)
		}
	case pickedCommand == MASTER:
		{
			output, jsonOUTPUT = esClient.CatMaster(catr)
		}
	case pickedCommand == ALIAIS:
		{
			output, jsonOUTPUT = esClient.CatAliases(catr)
		}
	case pickedCommand == REPOSITORIES:
		{
			output, jsonOUTPUT = esClient.CatRepositories(catr)
		}
	case pickedCommand == COUNT:
		{
			output, jsonOUTPUT = esClient.CatCount(catr)
		}
	case pickedCommand == INFO:
		{
			output, jsonOUTPUT = esClient.CatInfo(catr)
		}
	case pickedCommand == ROLE:
		{
			output, jsonOUTPUT = esClient.CatRoles(catr)
		}
	case pickedCommand == USER:
		{
			output, jsonOUTPUT = esClient.CatUsers(catr)
		}
	}

	if jsonOUTPUT {
		if pretty {
			var data interface{}
			json.Unmarshal(output, &data)
			j, _ := json.MarshalIndent(data, " ", "   ")
			fmt.Println(string(j))
		} else {
			fmt.Println(string(output))
		}
	}

	if !jsonOUTPUT {
		fmt.Print(string(output))
	}
}

//ElasticSearchClient is client object
type ElasticSearchClient struct {
	Host     string
	User     string
	Password string
	http     *http.Client
}

//NewElasticSearchClient to build the new elasticsearch client
func NewElasticSearchClient(host, user, password string) *ElasticSearchClient {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	client := &ElasticSearchClient{
		Host:     host,
		User:     user,
		Password: password,
		http: &http.Client{
			Timeout: time.Duration(5 * time.Second),
		},
	}

	return client
}

//CatHealth do /_cluster/health or /_cat/health depends upo then r.format parameter
func (client *ElasticSearchClient) CatHealth(r catRequest) ([]byte, bool) {
	url := "/_cluster/health"
	var query []string
	if r.format == "text" {
		url = "/_cat/health"
		if r.verbose {
			query = append(query, "v")
		}
	}
	return client.do(url, query, r)
}

//CatSnapshot do _cat/snapshots/{r.arg1}
func (client *ElasticSearchClient) CatSnapshot(r catRequest) ([]byte, bool) {
	url := "/_cat/snapshots/" + r.arg1
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}

	output, jsonOUTPUT := client.do(url, query, r)

	if r.def {
		snaps := []snapshots{}
		json.Unmarshal(output, &snaps)
		output = []byte("[")
		for _, i := range snaps {
			snapDef, _ := client.do("/_snapshot/"+r.arg1+"/"+i.ID, nil, r)
			output = append(output, snapDef...)
			output = append(output, ',')
		}
		output[len(output)-1] = ']'
	}

	return output, jsonOUTPUT
}

//CatAllocation do _cat/allocation
func (client *ElasticSearchClient) CatAllocation(r catRequest) ([]byte, bool) {
	url := "/_cat/allocation"
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	return client.do(url, query, r)
}

//CatNodes do _cat/nodes
func (client *ElasticSearchClient) CatNodes(r catRequest) ([]byte, bool) {
	url := "/_cat/nodes"
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	return client.do(url, query, r)
}

//CatPlugins do _cat/plugins
func (client *ElasticSearchClient) CatPlugins(r catRequest) ([]byte, bool) {
	url := "/_cat/plugins"
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	return client.do(url, query, r)
}

//CatTemplates do _cat/templates, with "_" do _template/{templateName}
func (client *ElasticSearchClient) CatTemplates(r catRequest) ([]byte, bool) {
	url := "/_cat/templates"

	if r.arg1 != "" {
		url += "/" + r.arg1
	}

	var query []string

	if r.format == "json" {
		query = append(query, "format=json")
	}

	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}

	if r.sortFields != "" {
		query = append(query, "s="+r.sortFields)
	} else {
		query = append(query, "s=n")
	}

	output, jsonOUTPUT := client.do(url, query, r)

	if r.def {
		templates := []templates{}
		json.Unmarshal(output, &templates)
		output = []byte("[")
		for _, i := range templates {
			templateDef, _ := client.do("/_template/"+i.Name, nil, r)
			output = append(output, templateDef...)
			output = append(output, ',')
		}
		output[len(output)-1] = ']'
	}

	return output, jsonOUTPUT
}

//CatMaster do _cat/master
func (client *ElasticSearchClient) CatMaster(r catRequest) ([]byte, bool) {
	url := "/_cat/master"
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	return client.do(url, query, r)
}

//CatIndices do _cat/indices, with "_" do _template/{indexName}
func (client *ElasticSearchClient) CatIndices(r catRequest) ([]byte, bool) {
	url := "/_cat/indices"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	if r.sortFields != "" {
		query = append(query, "s="+r.sortFields)
	} else {
		query = append(query, "s=i")
	}

	output, jsonOUTPUT := client.do(url, query, r)

	if r.def {
		indices := []indices{}
		json.Unmarshal(output, &indices)
		output = []byte("[")
		for _, i := range indices {
			indexDef, _ := client.do("/"+i.Index, nil, r)
			output = append(output, indexDef...)
			output = append(output, ',')
		}
		output[len(output)-1] = ']'
	}

	return output, jsonOUTPUT
}

//CatSegments do _cat/segments
func (client *ElasticSearchClient) CatSegments(r catRequest) ([]byte, bool) {
	url := "/_cat/segments"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	if r.sortFields != "" {
		query = append(query, "s="+r.sortFields)
	} else {
		query = append(query, "s=i")
	}
	return client.do(url, query, r)
}

//CatAliases do _cat/aliases
func (client *ElasticSearchClient) CatAliases(r catRequest) ([]byte, bool) {
	url := "/_cat/aliases"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	if r.sortFields != "" {
		query = append(query, "s="+r.sortFields)
	} else {
		query = append(query, "s=a")
	}
	return client.do(url, query, r)
}

//CatRepositories do _cat/repositories, with "_" do _snapshot/{repoName}
func (client *ElasticSearchClient) CatRepositories(r catRequest) ([]byte, bool) {
	url := "/_cat/repositories"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}

	output, jsonOUTPUT := client.do(url, query, r)

	if r.def {
		snaps := []snapshots{}
		json.Unmarshal(output, &snaps)
		output = []byte("[")
		for _, i := range snaps {
			snapDef, _ := client.do("/_snapshot/"+i.ID, nil, r)
			output = append(output, snapDef...)
			output = append(output, ',')
		}
		output[len(output)-1] = ']'
	}

	return output, jsonOUTPUT
}

//CatCount do _cat/count
func (client *ElasticSearchClient) CatCount(r catRequest) ([]byte, bool) {
	url := "/_cat/count"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	return client.do(url, query, r)
}

//CatInfo do /
func (client *ElasticSearchClient) CatInfo(r catRequest) ([]byte, bool) {
	url := "/"
	var query []string
	return client.do(url, query, r)
}

//CatRoles do _xpack/security/role
func (client *ElasticSearchClient) CatRoles(r catRequest) ([]byte, bool) {
	url := "/_xpack/security/role"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	return client.do(url, query, r)
}

//CatUsers do _xpack/security/user
func (client *ElasticSearchClient) CatUsers(r catRequest) ([]byte, bool) {
	url := "/_xpack/security/user"
	if r.arg1 != "" {
		url += "/" + r.arg1
	}
	var query []string
	if r.format == "json" {
		query = append(query, "format=json")
	}
	if r.format == "text" {
		if r.verbose {
			query = append(query, "v")
		}
	}
	if r.fields != "" {
		query = append(query, "h="+r.fields)
	}
	return client.do(url, query, r)
}

func (client *ElasticSearchClient) isCommandHelp(command string, args []string) bool {

	firstArg := ""
	if len(args) >= 2 {
		firstArg = strings.ToLower(args[1])
	}

	if firstArg == "help" {
		text := ""
		if command == HEALTH {
			text = "escat [OPTIONS] health"
		}

		if command == SNAPSHOTS {
			text = "escat [OPTIONS] snapshots respository-name [_]"
		}

		if command == ALLOCATION {
			text = "escat [OPTIONS] allocation"
		}

		if command == NODES {
			text = "escat [OPTIONS] nodes"
		}

		if command == PLUGINS {
			text = "escat [OPTIONS] plugins"
		}

		if command == MASTER {
			text = "escat [OPTIONS] master"
		}

		if command == INDICES {
			text = "escat [OPTIONS] indices [PREFIX] [_]"
		}

		if command == SEGMENTS {
			text = "escat [OPTIONS] segments [PREFIX]"
		}

		if command == ALIAIS {
			text = "escat [OPTIONS] aliases [PREFIX]"
		}

		if command == REPOSITORIES {
			text = "escat [OPTIONS] repositories [PREFIX] [_]"
		}

		if command == INFO {
			text = "escat [OPTIONS] info"
		}

		if command == COUNT {
			text = "escat [OPTIONS] count [PREFIX]"
		}

		if command == USER {
			text = "escat [OPTIONS] user [PREFIX]"
		}

		if command == ROLE {
			text = "escat [OPTIONS] role [PREFIX]"
		}

		if text != "" {
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintf(os.Stdout, "Usage : %s\n", text)
			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintf(os.Stdout, "Options:\n")
			flag.PrintDefaults()
			fmt.Fprintf(os.Stdout, "\n")

			return true
		}
	}

	return false
}

func (client *ElasticSearchClient) do(url string, query []string, r catRequest) ([]byte, bool) {
	if query != nil && len(query) > 0 {
		url += "?" + strings.Join(query, "&")
	}
	req, _ := http.NewRequest("GET", client.Host+url, nil)
	req.SetBasicAuth(client.User, client.Password)
	res, err := client.http.Do(req)
	if err != nil {
		panic(err)
	}

	responseJSON := false
	if res.Header.Get("Content-Type") == "application/json; charset=UTF-8" {
		responseJSON = true
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if body[len(body)-1] != '\n' {
		body = append(body, '\n')
	}

	return body, responseJSON
}

type indices struct {
	Index string `json:"index"`
}

type templates struct {
	Name string `json:"name"`
}

type snapshots struct {
	ID string `json:"id"`
}

type catRequest struct {
	format     string
	verbose    bool
	arg1       string
	fields     string
	sortFields string
	def        bool
	debug      bool
}
