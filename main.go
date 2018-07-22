package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"syscall"
)

/*PGCon postgres connection paramters */
type PGCon struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

// PGAliasFile Parsed config file
type PGAliasFile struct {
	Filepath string
	Conns    *map[string]PGCon
}

func readConfig(filepath string) *PGAliasFile {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalln("Fail to open the config file")
	}
	contents := string(raw)

	conns := make(map[string]PGCon)
	for _, line := range strings.Split(contents, "\n") {
		elems := strings.Split(line, "\t")
		conns[elems[0]] = PGCon{
			Host:     elems[1],
			Port:     elems[2],
			DB:       elems[3],
			User:     elems[4],
			Password: elems[5],
		}

	}

	return &PGAliasFile{Filepath: filepath, Conns: &conns}
}

func listAllConn(config *PGAliasFile) {
	if config.Conns == nil {
		panic("Unparsed config file")
	}

	conn := *config.Conns
	keys := make([]string, 0, len(conn))
	for key := range conn {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Printf("Alias in config file: %s\n\n", config.Filepath)
	for _, alias := range keys {
		c := conn[alias]
		fmt.Printf("  %-8s ", alias)
		if c.DB != "*" {
			fmt.Printf("%s@%s:%s/%s\n", c.User, c.Host, c.Port, c.DB)
		} else {
			fmt.Printf("%s@%s:%s\n", c.User, c.Host, c.Port)
		}
	}
}

func main() {

	configPath := path.Join(os.Getenv("HOME"), ".pgapass")

	configFile := flag.String("config", configPath, "PG connection alias")
	isPgdump := flag.Bool("pgdump", false, "Use pgdump instead of psql")
	ls := flag.Bool("l", false, "List all avialable connection and alias")

	flag.Parse()

	c := readConfig(*configFile)

	if *ls {
		listAllConn(c)
		os.Exit(0)
	}

	prog := "psql"

	if *isPgdump {
		prog = "pg_dump"
	}

	binary, err := exec.LookPath(prog)

	if err != nil {
		log.Fatalln("Please install psql first")
	}

	con, exist := (*c.Conns)[flag.Arg(0)]

	if !exist {
		log.Fatalf("No connection found for alias: %s", os.Args[1])
	}

	pgargs := []string{prog, "-h", con.Host, "-p", con.Port, "-U", con.User}
	if con.DB != "*" {
		if *isPgdump {
			pgargs = append(pgargs, "-d", con.DB)
		} else {
			pgargs = append(pgargs, con.DB)
		}
	}
	args := append(pgargs, flag.Args()[1:]...)
	env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", con.Password))

	syscall.Exec(binary, args, env)
}
