package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	awscdk.StackProps
}

func NewGoCdkStack(scope constructs.Construct, id string, props *GoCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("userTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	lambdaFunc := awslambda.NewFunction(stack, jsii.String("lambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code: awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"),nil),
		Handler: jsii.String("main"),
	})

	table.GrantReadWriteData(lambdaFunc)
	apiGateway := awsapigateway.NewRestApi(stack, jsii.String("apiGateway"), &awsapigateway.RestApiProps{
		CloudWatchRole: jsii.Bool(true),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})
	integration := awsapigateway.NewLambdaIntegration(lambdaFunc, nil)

	authRoute := apiGateway.Root().AddResource(jsii.String("auth"),nil)

	registerRoute := authRoute.AddResource(jsii.String("register"),nil)
	registerRoute.AddMethod(jsii.String("POST"),integration, nil)

	loginRoute := authRoute.AddResource(jsii.String("login"),nil)
	loginRoute.AddMethod(jsii.String("POST"),integration, nil)

	protectedRoute := apiGateway.Root().AddResource(jsii.String("protected"),nil)
	protectedRoute.AddMethod(jsii.String("GET"),integration,nil)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoCdkStack(app, "GoCdkStack", &GoCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
