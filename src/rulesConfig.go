package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"omr/api/conf"
	"os"
)

type RuleScenarioCondition struct {
	ScenarioName  string `json:"scenarioName"`
	ChannelName   string `json:"channelName"`
	HasAttachment bool   `json:"hasAttachment"`
	HasSignature  bool   `json:"hasSignature"`
	IsSmartCard   bool   `json:"isSmartCard"`
}

type RuleScenarioResult struct {
	ConditionName        string `json:"conditionName"`
	IsSharepointUpload   bool   `json:"isSharepointUpload"`
	IsEFormGenerate      bool   `json:"isEFormGenerate"`
	IsShortFormGenerate  bool   `json:"isShortFormGenerate"`
	IsAttachmentGenerate bool   `json:"isAttachmentGenerate"`
	IsWarningMessage     bool   `json:"isWarningMessage"`
	WarningMessage       string `json:"warningMessage"`
	IsContactItemEForm   bool   `json:"isContactItemEForm"`
}

type RuleScenario struct {
	RuleCondition RuleScenarioCondition `json:"ruleScenarioCondition"`
	RuleResult    RuleScenarioResult    `json:"ruleScenarioResult"`
}

type ListRuleScenario struct {
	Scenario []RuleScenario `json:"genPDFRules"`
}

var ruleBase map[RuleScenarioCondition]RuleScenario
var (
	DEFAULT_CONDITION_NAME        = Get("form.default.conditionName")
	DEFAULT_IS_SHAREPOINT_UPLOAD  = GetBool("form.default.isSharepointUpload")
	DEFAULT_IS_EFORM_GEN          = GetBool("form.default.IsEFormGenerate")
	DEFAULT_IS_SHORT_FORM_GEN     = GetBool("form.default.isShortFormGenerate")
	DEFAULT_IS_ATTACHMENT_GEN     = GetBool("form.default.isAttachmentGenerate")
	DEFAULT_IS_WARNING_MESSAGE    = GetBool("form.default.isWarningMessage")
	DEFAULT_WARNING_MESSAGE       = Get("form.default.warningMessage")
	DEFAULT_IS_CONTACT_ITEM_EFORM = GetBool("form.default.isContactItemEForm")
)

var url = conf.Get("ad.url")

func InitialRule() error {
	fileLocation := getFileLocation("form.rules.json")

	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		fmt.Printf("RULE_001:  %v\n", err)
	}

	listRuleScenario := new(ListRuleScenario)
	err = json.Unmarshal(file, listRuleScenario)
	if err != nil {
		fmt.Printf("RULE_002:  %v\n", err)
	}

	timeout := conf.Get("ad.timeout")
	ruleBase = make(map[RuleScenarioCondition]RuleScenario)

	for _, item := range listRuleScenario.Scenario {
		ruleBase[item.RuleCondition] = item
	}

	err = checkDuplicateRule(listRuleScenario.Scenario, ruleBase)
	return err
}

func checkDuplicateRule(jsonFile []RuleScenario, jsonMap map[RuleScenarioCondition]RuleScenario) (err error) {
	if len(jsonFile) != len(jsonMap) {
		err = errors.New("RULE_003: found duplicate rule, please check file formrules.json")
	}
	return
}

func GetRules(condition RuleScenarioCondition) (RuleScenarioResult, error) {
	var err error
	if ruleBase == nil {
		err = InitialRule()
		if err != nil {
			fmt.Printf("InitialRule err=%v", err)
		}
	}

	ruleScenarioResult := RuleScenarioResult{}
	if ruleScenario, ok := ruleBase[condition]; ok {
		ruleScenarioResult = ruleScenario.RuleResult
	} else {
		// Default rule when not match
		ruleScenarioResult.ConditionName = DEFAULT_CONDITION_NAME
		ruleScenarioResult.IsSharepointUpload = DEFAULT_IS_SHAREPOINT_UPLOAD
		ruleScenarioResult.IsEFormGenerate = DEFAULT_IS_EFORM_GEN
		ruleScenarioResult.IsShortFormGenerate = DEFAULT_IS_SHORT_FORM_GEN
		ruleScenarioResult.IsAttachmentGenerate = DEFAULT_IS_ATTACHMENT_GEN
		ruleScenarioResult.IsWarningMessage = DEFAULT_IS_WARNING_MESSAGE
		ruleScenarioResult.WarningMessage = DEFAULT_WARNING_MESSAGE
		ruleScenarioResult.IsContactItemEForm = DEFAULT_IS_CONTACT_ITEM_EFORM
	}

	return ruleScenarioResult, err
}

func getFileLocation(location string) string {
	return os.Getenv("GOPATH") + Get(location)
}
