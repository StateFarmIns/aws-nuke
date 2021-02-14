package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Budget struct {
	svc  *budgets.Budgets
	name *string
	// tags []*budgets.Tag
}

// if there are no tags what should properties be?

func init() {
	register("Budgets", ListBudgets)
}

func ListBudgets(sess *session.Session) ([]Resource, error) {
	svc := budgets.New(sess)

	// resources := []Resource{}

	fmt.Println("Entered Budgets Resource")

	// Lookup current account ID
	callerID, err := sts.New(sess).GetCallerIdentity(nil) //&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	// accountID := callerID.Account
	fmt.Printf("account_num: %s \n", *callerID.Account)

	params := &budgets.DescribeBudgetsInput{
		AccountId:  aws.String(*callerID.Account),
		MaxResults: aws.Int64(100),
	}

	output, err := svc.DescribeBudgets(params)
	if err != nil {
		fmt.Printf("error: %s \n", err)
		return nil, err
	}

	fmt.Printf("len=%d\n", len(output.Budgets))

	for _, bud := range output.Budgets {
		fmt.Println(bud.BudgetName)
		// resources = append(resources, Budget{
		// 	svc:  svc,
		// 	name: bud.BudgetName,
		// })
	}

	// return resources, nil
	return nil, nil
}
