package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
)

func provision(db *sql.DB, plan string) string {

	cacheparametergroupname := "memcached-14-small"
	cachenodetype := os.Getenv("SMALL_INSTANCE_TYPE")
	numcachenodes := int64(1)
	billingcode := "pre-provisioned"
	u, err := uuid.NewV4()
	name := os.Getenv("NAME_PREFIX") + "-" + strings.Split(u.String(), "-")[0]

	if plan == "small" {
		cacheparametergroupname = os.Getenv("SMALL_PARAMETER_GROUP")
		cachenodetype = os.Getenv("SMALL_INSTANCE_TYPE")
		numcachenodes = int64(1)
	}
	if plan == "medium" {
		cacheparametergroupname = os.Getenv("MEDIUM_PARAMETER_GROUP")
		cachenodetype = os.Getenv("MEDIUM_INSTANCE_TYPE")
		numcachenodes = int64(1)
	}
	if plan == "large" {
		cacheparametergroupname = os.Getenv("LARGE_PARAMETER_GROUP")
		cachenodetype = os.Getenv("LARGE_INSTANCE_TYPE")
		numcachenodes = int64(1)
	}

	svc := elasticache.New(session.New(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}))

	params := &elasticache.CreateCacheClusterInput{
		CacheClusterId:          aws.String(name),
		AutoMinorVersionUpgrade: aws.Bool(true),
		CacheNodeType:           aws.String(cachenodetype),
		CacheParameterGroupName: aws.String(cacheparametergroupname),
		CacheSubnetGroupName:    aws.String("memcached-subnet-group"),
		Engine:                  aws.String("Memcached"),
		NumCacheNodes:           aws.Int64(numcachenodes),
		Port:                    aws.Int64(6379),
		SecurityGroupIds: []*string{
			aws.String(os.Getenv("ELASTICACHE_SECURITY_GROUP")),
		},
		Tags: []*elasticache.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
			{
				Key:   aws.String("billingcode"),
				Value: aws.String(billingcode),
			},
		},
	}
	resp, err := svc.CreateCacheCluster(params)

	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	return name
}

func insertnew(db *sql.DB, name string, plan string, claimed string) {
	var newname string
	err :=db.QueryRow("INSERT INTO provision(name,plan,claimed) VALUES($1,$2,$3) returning name;", name, plan, claimed).Scan(&newname)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func main() {

	uri := os.Getenv("BROKER_DB")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer db.Close()

	// setup the database (or modify it as necessary)
	buf, err := ioutil.ReadFile("create.sql")
	if err != nil {
		log.Fatalf("Unable to read create.sql: %s\n", err)
	}
	_, err = db.Exec(string(buf))
	if err != nil {
		log.Fatal("Unable to create database: %s\n", err)
	}

	newname := "new"

	provisionsmall, _ := strconv.Atoi(os.Getenv("PROVISION_SMALL"))
	provisionmedium, _ := strconv.Atoi(os.Getenv("PROVISION_MEDIUM"))
	provisionlarge, _ := strconv.Atoi(os.Getenv("PROVISION_LARGE"))

	var smallcount int
	err = db.QueryRow("SELECT count(*) as smallcount from provision where plan='small' and claimed='no'").Scan(&smallcount)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if smallcount < provisionsmall {
		newname = provision(db, "small")
		insertnew(db,newname, "small", "no")
	}

	var mediumcount int
	err = db.QueryRow("SELECT count(*) as mediumcount from provision where plan='medium' and claimed='no'").Scan(&mediumcount)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if mediumcount < provisionmedium {
		newname = provision(db, "medium")
		insertnew(db, newname, "medium", "no")
	}

	var largecount int
	err = db.QueryRow("SELECT count(*) as largecount from provision where plan='large' and claimed='no'").Scan(&largecount)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if largecount < provisionlarge {
		newname = provision(db, "large")
		insertnew(db,newname, "large", "no")
	}

}
