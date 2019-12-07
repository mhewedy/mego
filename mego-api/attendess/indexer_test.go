package attendess

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

const emptyRespHead = `<?xml version="1.0" encoding="utf-8"?><s:Envelope 
   xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
  <s:Header>
    <h:ServerVersionInfo MajorVersion="15" 
                         MinorVersion="0" 
                         MajorBuildNumber="349" 
                         MinorBuildNumber="0" 
                         Version="Exchange2013" 
                         xmlns="http://schemas.microsoft.com/exchange/services/2006/types" 
                         xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
                         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
                         xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" />
  </s:Header>
  <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
          xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <FindPeopleResponse ResponseClass="Success" 
                        xmlns="http://schemas.microsoft.com/exchange/services/2006/messages">
      <ResponseCode>NoError</ResponseCode>
      <People>`

const emptyRespTail = `</People>
      <TotalNumberOfPeopleInView>1</TotalNumberOfPeopleInView>
    </FindPeopleResponse>
  </s:Body>
</s:Envelope>`

const terry = `<Persona xmlns="http://schemas.microsoft.com/exchange/services/2006/types">
 <PersonaId Id="some id" />
 <CreationTime>2012-01-11T22:25:37Z</CreationTime>
 <DisplayName>Terry Adams</DisplayName>
 <EmailAddress>
	 <EmailAddress>terry@litwareinc.com</EmailAddress>
 </EmailAddress>
</Persona>`

const abbas = `<Persona xmlns="http://schemas.microsoft.com/exchange/services/2006/types">
 <PersonaId Id="some id 2" />
 <CreationTime>2012-01-11T22:25:37Z</CreationTime>
 <DisplayName>Abbas Adams</DisplayName>
 <EmailAddress>
	 <EmailAddress>abbas@litwareinc.com</EmailAddress>
 </EmailAddress>
</Persona>`

type mockEWSClient struct {
}

func (m mockEWSClient) SendAndReceive(body []byte) ([]byte, error) {
	if strings.Contains(string(body), "<m:QueryString>a</m:QueryString>") {
		return []byte(emptyRespHead + terry + abbas + emptyRespTail), nil
	}

	if strings.Contains(string(body), "<m:QueryString>b</m:QueryString>") {
		return []byte(emptyRespHead + terry + emptyRespTail), nil
	}

	return []byte(emptyRespHead + emptyRespTail), nil
}

func (m mockEWSClient) GetEWSAddr() string {
	return ""
}

func (m mockEWSClient) GetUsername() string {
	return ""
}

func Test_indexAttendees(t *testing.T) {

	EWSClient = &mockEWSClient{}
	indexAttendees()

	assert.Equal(t, 2, len(attendeesIndex))
	assert.ElementsMatch(t, attendeesIndex, []Attendee{
		{
			DisplayName:  "Terry Adams",
			Title:        "",
			EmailAddress: "terry@litwareinc.com",
			Image:        "",
		},
		{
			DisplayName:  "Abbas Adams",
			Title:        "",
			EmailAddress: "abbas@litwareinc.com",
			Image:        "",
		},
	})
}

func Test_searchAttendees(t *testing.T) {

	attendeesIndex = []Attendee{
		{
			DisplayName:  "Terry Adams",
			Title:        "",
			EmailAddress: "terry@litwareinc.com",
			Image:        "",
		},
		{
			DisplayName:  "Abbas Fernas",
			Title:        "",
			EmailAddress: "abbas@litwareinc.com",
			Image:        "",
		},
	}

	type args struct {
		q string
	}
	tests := []struct {
		name string
		args args
		want []Attendee
	}{
		{
			name: "test start with",
			args: args{q: "terry"},
			want: []Attendee{
				{
					DisplayName:  "Terry Adams",
					Title:        "",
					EmailAddress: "terry@litwareinc.com",
					Image:        "",
				},
			},
		}, {
			name: "test contains",
			args: args{q: "erry"},
			want: []Attendee{
				{
					DisplayName:  "Terry Adams",
					Title:        "",
					EmailAddress: "terry@litwareinc.com",
					Image:        "",
				},
			},
		}, {
			name: "test contains of last name",
			args: args{q: "ams"},
			want: []Attendee{
				{
					DisplayName:  "Terry Adams",
					Title:        "",
					EmailAddress: "terry@litwareinc.com",
					Image:        "",
				},
			},
		}, {
			name: "test non existence",
			args: args{q: "xyz"},
			want: []Attendee{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchAttendees(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchAttendees() = %v, want %v", got, tt.want)
			}
		})
	}
}