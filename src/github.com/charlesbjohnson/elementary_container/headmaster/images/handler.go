package images

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/charlesbjohnson/elementary_container/communication"
	"github.com/charlesbjohnson/elementary_container/headmaster"
	"github.com/gorilla/context"
)

func startHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	clientUrl := (&url.URL{
		Scheme: os.Getenv("CLIENT_URL_PROTOCOL"),
		Host:   os.Getenv("CLIENT_URL_HOST"),
		Path:   os.Getenv("CLIENT_URL_PATH"),
	}).String()

	serverUrl := (&url.URL{
		Scheme: os.Getenv("SERVER_URL_PROTOCOL"),
		Host:   os.Getenv("SERVER_URL_HOST"),
	}).String()

	context := struct {
		ClientUrl string
		ServerUrl string
	}{clientUrl, serverUrl}

	application.View.Plain(response, http.StatusOK, "execute.sh.tmpl", context)
}

func createHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	decoder := json.NewDecoder(request.Body)
	input := new(communication.ImageCreateInput)

	if err := decoder.Decode(input); err != nil {
		application.View.JSON(response, http.StatusBadRequest, communication.Error{Message: err.Error()})
		return
	}

	if ok, err := regexp.MatchString("[[:alnum:]]{64}", input.Checksum); err != nil || !ok {
		application.View.JSON(response, http.StatusBadRequest, communication.Error{Message: "Invalid checksum"})
		return
	}

	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))
	stsService := sts.New(awsConfig)

	assumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/os/os-client-role", os.Getenv("AWS_ACCOUNT_ID"))),
		RoleSessionName: aws.String("os-client-role"),
	}
	assumeRoleOutput, err := stsService.AssumeRole(assumeRoleInput)
	if err != nil {
		application.View.JSON(response, http.StatusInternalServerError, communication.Error{Message: err.Error()})
		return
	}

	awsCredentials := credentials.NewStaticCredentials(
		*assumeRoleOutput.Credentials.AccessKeyId,
		*assumeRoleOutput.Credentials.SecretAccessKey,
		*assumeRoleOutput.Credentials.SessionToken,
	)
	awsConfig = awsConfig.WithCredentials(awsCredentials)
	s3Service := s3.New(awsConfig)

	listObjectsInput := &s3.ListObjectsInput{Bucket: aws.String(os.Getenv("AWS_OS_IMAGES_BUCKET"))}
	listObjectsOutput, err := s3Service.ListObjects(listObjectsInput)
	if err != nil {
		application.View.JSON(response, http.StatusInternalServerError, communication.Error{Message: fmt.Sprintf("aws listobjects error: %s", err.Error())})
		return
	}

	exists, key := false, ""
	for _, object := range listObjectsOutput.Contents {
		if strings.Contains(*object.Key, input.Checksum) {
			exists = true
			key = *object.Key
			break
		}
	}

	application.View.JSON(response, http.StatusOK, communication.ImageCreateOutput{
		AWSAccessKeyId:     *assumeRoleOutput.Credentials.AccessKeyId,
		AWSSecretAccessKey: *assumeRoleOutput.Credentials.SecretAccessKey,
		AWSSessionToken:    *assumeRoleOutput.Credentials.SessionToken,
		AWSRegion:          *s3Service.Config.Region,
		AWSS3Bucket:        *listObjectsInput.Bucket,
		AWSS3ObjectKey:     key,
		Exists:             exists,
	})
}

func commitHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	decoder := json.NewDecoder(request.Body)
	input := new(communication.ImageCommitInput)

	if err := decoder.Decode(input); err != nil {
		application.View.JSON(response, http.StatusBadRequest, communication.Error{Message: err.Error()})
		return
	}

	awsCredentials := credentials.NewStaticCredentials(
		input.AWSAccessKeyId,
		input.AWSSecretAccessKey,
		input.AWSSessionToken,
	)
	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")).WithCredentials(awsCredentials)
	s3Service := s3.New(awsConfig)

	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_OS_IMAGES_BUCKET")),
		Key:    aws.String(input.AWSS3ObjectKey),
	}
	if _, err := s3Service.HeadObject(headObjectInput); err != nil {
		application.View.JSON(response, http.StatusInternalServerError, communication.Error{Message: fmt.Sprintf("aws headobject error: %s", err.Error())})
		return
	}

	response.WriteHeader(http.StatusOK)
	// TODO
	// send a message to an image importer channel
	// return success message
}

func addHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	application.Log.Info("add")

	vars := mux.Vars(request)
	application.Log.Info(vars["checksum"])
}

func pollHandler(response http.ResponseWriter, request *http.Request) {
	server := context.Get(request, "server").(headmaster.Server)
	application := server.Application()

	application.Log.Info("poll")
	// TODO
	// check the map for their checksum
	// attempt to pull checksum from the queue
	// race against a time out
	// if time out occurs, put their checksum in the map (if its not already there)
}
