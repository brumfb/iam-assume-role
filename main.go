package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var role string
var account string
var sessionName string
var errLog *log.Logger

func init() {
	flag.StringVar(&role, "role", "", "Role to assume")
	flag.StringVar(&account, "account", "", "AccountID containing role to assume")
	flag.StringVar(&sessionName, "session", "", "Session name")
}
func main() {
	flag.Parse()
	errLog = log.New(os.Stderr, "", 0)

	if role == "" {
		errLog.Panicf("Missing --role\n")
	}

	if account == "" {
		errLog.Panicf("missing --account\n")
	}

	result := assumeRole()
	exportCredentials(result)
}

func getSessionName() string {
	if sessionName != "" {
		return sessionName
	}

	host, ok := os.Hostname()
	if ok == nil {
		return fmt.Sprintf("%s-%s-%s-%d", host, account, role, time.Now().Unix())
	}
	return fmt.Sprintf("%s-%s-%d", account, role, time.Now().Unix())
}

func assumeRole() *sts.Credentials {
	service := sts.New(session.New())
	roleInfo := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(3600),
		ExternalId:      aws.String(account),
		RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", account, role)),
		RoleSessionName: aws.String(getSessionName()),
	}

	result, err := service.AssumeRole(roleInfo)
	if err != nil {
		errLog.Panic(err)
	}

	return result.Credentials
}

func export(name string, value string) {
	fmt.Printf("export %s='%s'\n", name, value)
}
func exportCredentials(credentials *sts.Credentials) {
	export("AWS_ACCESS_KEY_ID", *credentials.AccessKeyId)
	export("AWS_SECRET_ACCESS_KEY", *credentials.SecretAccessKey)
	export("AWS_SESSION_TOKEN", *credentials.SessionToken)
}
