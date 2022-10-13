package aws_go_cost

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/zhangtaomox/tablib"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func GenerateCost(region string) {

	//Must be in YYYY-MM-DD Format
	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()
	metrics := []string{
		"BlendedCost",
		"UnblendedCost",
		"UsageQuantity",
	}
	// Initialize a session in us-east-1 that the SDK will use to load credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		exitErrorf("Unable to generate report, %v", err)
	}

	acctMap := make(map[string]string)

	osvc := organizations.New(sess)
	oresult, _ := osvc.ListAccounts(&organizations.ListAccountsInput{})
	// We try and keep track of friendly names for accounts
	for _, acct := range oresult.Accounts {
		acctMap[*acct.Id] = *acct.Name
	}

	// Create Cost Explorer Service Client
	svc := costexplorer.New(sess)
	ctx := context.Background()
	var results []*costexplorer.ResultByTime

	// pagination handling
	var paginationToken string = ""

	for {
		params := &costexplorer.GetCostAndUsageInput{
			TimePeriod: &costexplorer.DateInterval{
				Start: aws.String(start.Format("2006-01") + "-01"),
				End:   aws.String(end.Format("2006-01") + "-01"),
			},
			Granularity: aws.String("MONTHLY"),
			GroupBy: []*costexplorer.GroupDefinition{
				{
					Type: aws.String("DIMENSION"),
					Key:  aws.String("SERVICE"),
				},
				{
					Type: aws.String("DIMENSION"),
					Key:  aws.String("LINKED_ACCOUNT"),
				},
			},
			Metrics: aws.StringSlice(metrics),
		}
		if paginationToken != "" {
			params.NextPageToken = aws.String(paginationToken)
		}

		result, err := svc.GetCostAndUsageWithContext(
			ctx,
			params,
		)
		if err != nil {
			fmt.Println("Error happened.", err)
			os.Exit(1)
		}
		results = append(results, result.ResultsByTime...)
		if result.NextPageToken == nil {
			break
		}
		paginationToken = *result.NextPageToken
	}

	dataset := tablib.NewDataSet().SetTitle("tablib").SetHeaders([]string{
		"Start", "End", "Account ID", "Account Name", "Service", "Unit", "cost",
	})
	for _, p := range results {
		for _, g := range p.Groups {
			acctID := *g.Keys[1]
			serviceName := *g.Keys[0]
			fname := acctMap[acctID]
			if fname == "" {
				fname = acctID
			}
			_ = dataset.Append([]string{
				*p.TimePeriod.Start,
				*p.TimePeriod.End,
				acctID,
				fname, serviceName,
				*g.Metrics["BlendedCost"].Unit,
				formatNumber(*g.Metrics["BlendedCost"].Amount),
			})
			//cost := formatNumber(*g.Metrics["BlendedCost"].Amount)
			//fmt.Printf("%s | %s | %s | %s | %s | %s %s\n", *p.TimePeriod.Start, *p.TimePeriod.End, acctID, fname, serviceName, *g.Metrics["BlendedCost"].Unit, cost)
		}
	}

	fCsv, err := os.OpenFile("/tmp/report.csv", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		exitErrorf("Unable to generate report, %v", err)
	}
	_ = dataset.Export(fCsv, tablib.CSV)

	defer fCsv.Close()
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func formatNumber(s string) string {
	f, _ := strconv.ParseFloat(s, 64)
	return fmt.Sprintf("%.2f", f)
}
