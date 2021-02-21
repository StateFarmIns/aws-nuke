package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

// if there are no tags what should properties be?

func init() {
	register("Budget", ListBudgets)
	//register("Budgets", ListBudgets)
}

type Budget struct {
	svc       *budgets.Budgets
	name      *string
	accountid *string
}

// when include a region in the config, only works with global
// error: service 'budgets' is global, but the session is not
// session.go line #235
func ListBudgets(sess *session.Session) ([]Resource, error) {
	svc := budgets.New(sess)

	fmt.Println("Entered Budgets Resource")

	//, &aws.Config{Region: aws.String(awsutil.DefaultRegionID)}).GetCallerIdentity(nil)
	// Lookup current account ID
	identityOutput, err := sts.New(sess).GetCallerIdentity(nil)
	if err != nil {
		fmt.Printf("sts error: %s \n", err)
		return nil, err
	}
	accountID := identityOutput.Account

	resources := []Resource{}

	fmt.Printf("account_num: %s \n", *accountID)

	params := &budgets.DescribeBudgetsInput{
		AccountId:  aws.String(*accountID),
		MaxResults: aws.Int64(100),
	}

	output, err := svc.DescribeBudgets(params)
	if err != nil {
		fmt.Printf("budgets error: %s \n", err)
		return nil, err
	}

	fmt.Printf("len=%d\n", len(output.Budgets))

	for _, bud := range output.Budgets {
		fmt.Printf("%s\n", *bud.BudgetName)
		resources = append(resources, &Budget{
			svc:  svc,
			name: bud.BudgetName,
		})
	}

	// return resources, nil
	return nil, nil
}

// fully implement remove later
func (b *Budget) Remove() error {

	_, err := b.svc.DeleteBudget(&budgets.DeleteBudgetInput{
		AccountId:  b.accountid,
		BudgetName: b.name,
	})

	return err
}

func (b *Budget) Properties() types.Properties {
	return types.NewProperties().
		Set("name", *b.name).
		Set("accoundid", *b.accountid)
}

func (b *Budget) String() string {
	return *b.name
}
