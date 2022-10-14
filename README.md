# CronJob AWS Cost Reporting
The application sends an email with a CSV file where you can review the expenses by Service and Accounts of the organization.


## Required parameters
- Region: AWS Region
- From: From email address
- To: To email address
- AWS Credentials

### Local implementation
```
export From=sender@email.com
export To=recipient@email.com
export Region=us-east-1

go run main.go
```
This generates a CSV file in the path **/tmp/report.csv** and then sends an email with the following information 

## kubernetes

In [k8s](./k8s) we can find the maifest to deploy the CronJob, for the correct operation we must create a secret called **cost-export** with the following data:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_DEFAULT_REGION
- FROM
- TO