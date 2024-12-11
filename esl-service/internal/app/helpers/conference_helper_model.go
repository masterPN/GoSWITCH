package helpers

import (
	"esl-service/internal/data"
	"sync"
	"time"

	"github.com/0x19/goesl"
)

type ConferenceHelperInterface interface {
	LoadConfiguration() (int, string)
	InitiateConferenceCalls(client *goesl.Client, conferenceData []string, externalDomain string, sipPort int, message map[string]string) error
	ValidateRadiusAndHandleConference(client *goesl.Client, conferenceInitData []string, msg map[string]string) bool
	FetchOptimalRouteData(destination string) (data.OptimalRouteData, error)
	FetchInternalCodemapping(internalCode string) (data.InternalCodemappingData, error)
	OriginateCallToOperator(client *goesl.Client, conferenceData []string, operatorCode int, operatorPrefix, externalDomain string, sipPort int) bool
	waitForCallToBeEstablished(client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool
	isWithinTimeout(start time.Time) bool
	handleBackgroundApiCall(client *goesl.Client, startTime time.Time, once *sync.Once)
	processCalleeAndConnection(message *goesl.Message, client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool
	handleConnectedCall(client *goesl.Client, baseClass int, conferenceName string) bool
	handleCalleeIssue(client *goesl.Client, message *goesl.Message, operatorPrefix, destination, conferenceName string)
	notifyCalleeIssue(message *goesl.Message, operatorPrefix, destination string)
	handleReadError(err error)
	isConnected(msg *goesl.Message, operatorPrefix, destination string) bool
	isCalleeUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool
	isOperatorUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool
	logOperatorIssue(msg *goesl.Message, operatorPrefix string)
}

type Models struct {
	ConferenceHelper ConferenceHelperInterface
}

func NewModels(client *goesl.Client, message *goesl.Message) Models {
	return Models{
		ConferenceHelper: ConferenceHelperModel{client: client, message: message},
	}
}
