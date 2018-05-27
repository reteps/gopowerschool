package gopowerschool

import (
	"fmt"
)

func Client(url string) *PublicPortalServiceJSONPortType {
	auth := DigestAuth{Login: "pearson", Password: "m0bApP5"}
	if url[len(url)-1] != '/' {
		url += "/"
	}
	wsdl_url := fmt.Sprintf("%s/pearson-rest/services/PublicPortalServiceJSON?wsdl", url)
	return NewPublicPortalServiceJSONPortType(wsdl_url, true, &auth)
}
func (client *PublicPortalServiceJSONPortType) CreateUserSessionAndStudent(username, password string) (*UserSessionVO, int64, error) {

	PublicPortalLogin := LoginToPublicPortal{Username: username, Password: password}
	response, err := client.LoginToPublicPortal(&PublicPortalLogin)
	if err != nil {
		return nil, 0, err
	}
	if response.Return_.MessageVOs != nil {
		return nil, 0, fmt.Errorf("error: %s - %s", response.Return_.MessageVOs[0].Title, response.Return_.MessageVOs[0].Description)
	}
	newSession := UserSessionVO{
		UserId:            response.Return_.UserSessionVO.UserId,
		ServiceTicket:     response.Return_.UserSessionVO.ServiceTicket,
		ServerInfo:        &ServerInfo{ApiVersion: response.Return_.UserSessionVO.ServerInfo.ApiVersion},
		ServerCurrentTime: response.Return_.UserSessionVO.ServerCurrentTime,
		UserType:          response.Return_.UserSessionVO.UserType}
	return &newSession, int64(response.Return_.UserSessionVO.StudentIDs[0]), nil
}
func (client *PublicPortalServiceJSONPortType) GetStudent(username, password string) (*StudentDataVO, error) {
	session, userID, err := client.CreateUserSessionAndStudent(username, password)
	if err != nil {
		return nil, err
	}
	studentDataArguments := GetStudentData{UserSessionVO: session, StudentIDs: []int64{userID}, Qil: &QueryIncludeListVO{Includes: []int32{1}}}
	student, err := client.GetStudentData(&studentDataArguments)
	if err != nil {
		panic(err)
	}
	return student.Return_.StudentDataVOs[0], nil
}
