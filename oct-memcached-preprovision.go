package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
)

func provision(plan string) string {

	cacheparametergroupname := "memcached-14-small"
	cachenodetype := os.Getenv("SMALL_INSTANCE_TYPE")
	numcachenodes := int64(1)
	billingcode := "pre-provisioned"
	u, err := uuid.NewV4()
	name := os.Getenv("NAME_PREFIX") + "-" + strings.Split(u.String(), "-")[0]
	fmt.Println(name)

	if plan == "small" {
		cacheparametergroupname = "memcached-14-small"
		cachenodetype = os.Getenv("SMALL_INSTANCE_TYPE")
		numcachenodes = int64(1)
	}
	if plan == "medium" {
		cacheparametergroupname = "memcached-14-medium"
		cachenodetype = os.Getenv("MEDIUM_INSTANCE_TYPE")
		numcachenodes = int64(1)
	}
	if plan == "large" {
		cacheparametergroupname = "memcached-14-large"
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

	fmt.Println(resp)
	return name
}

func insertnew(name string, plan string, claimed string) {
	uri := os.Getenv("BROKER_DB")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var newname string
	err = db.QueryRow("INSERT INTO provision(name,plan,claimed) VALUES($1,$2,$3) returning name;", name, plan, claimed).Scan(&newname)

	if err != nil {
		defer db.Close()
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(newname)
}

func main() {

	uri := os.Getenv("BROKER_DB")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		defer db.Close()
		fmt.Println(err)
		os.Exit(2)
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
	fmt.Println(smallcount)

	if smallcount < provisionsmall {
		newname = provision("small")
		fmt.Println(newname)
		insertnew(newname, "small", "no")
	}

	var mediumcount int
	err = db.QueryRow("SELECT count(*) as mediumcount from provision where plan='medium' and claimed='no'").Scan(&mediumcount)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(mediumcount)

	if mediumcount < provisionmedium {
		newname = provision("medium")
		fmt.Println(newname)
		insertnew(newname, "medium", "no")
	}

	var largecount int
	err = db.QueryRow("SELECT count(*) as largecount from provision where plan='large' and claimed='no'").Scan(&largecount)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(largecount)

	if largecount < provisionlarge {
		newname = provision("large")
		fmt.Println(newname)
		insertnew(newname, "large", "no")
	}

}
