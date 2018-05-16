package gopowerschool

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type Locale struct {
	XMLName xml.Name `xml:"locale"`

	ISO3Country            string `xml:"ISO3Country,omitempty"`
	ISO3Language           string `xml:"ISO3Language,omitempty"`
	Country                string `xml:"country,omitempty"`
	DisplayCountry         string `xml:"displayCountry,omitempty"`
	DisplayLanguage        string `xml:"displayLanguage,omitempty"`
	DisplayName            string `xml:"displayName,omitempty"`
	DisplayScript          string `xml:"displayScript,omitempty"`
	DisplayVariant         string `xml:"displayVariant,omitempty"`
	ExtensionKeys          *ExSet `xml:"extensionKeys,omitempty"`
	Language               string `xml:"language,omitempty"`
	Script                 string `xml:"script,omitempty"`
	UnicodeLocaleAttribute *UASet `xml:"unicodeLocaleAttributes,omitempty"`
	UnicodeLocaleKeys      *ULSet `xml:"unicodeLocaleKeys,omitempty"`
	Variant                string `xml:"variant,omitempty"`
}

// Workaround
type ExSet struct {
	XMLName xml.Name `xml:"http://util.java/xsd extensionKeys"`

	Empty bool `xml:"empty,omitempty"`
}
type UASet struct {
	XMLName xml.Name `xml:"http://util.java/xsd unicodeLocaleAttributes"`

	Empty bool `xml:"empty,omitempty"`
}

type ULSet struct {
	XMLName xml.Name `xml:"http://util.java/xsd unicodeLocaleKeys"`

	Empty bool `xml:"empty,omitempty"`
}
type Login struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd login"`

	Username string `xml:"username,omitempty"`
	Password string `xml:"password,omitempty"`
	UserType int32  `xml:"userType,omitempty"`
}

type LoginResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd loginResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type Logout struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd logout"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
}

type LogoutResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd logoutResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type GetSchoolMapBySchoolNumber struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getSchoolMapBySchoolNumber"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	SchoolNumber  int64          `xml:"schoolNumber,omitempty"`
}

type GetSchoolMapBySchoolNumberResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getSchoolMapBySchoolNumberResponse"`

	Return_ []byte `xml:"return,omitempty"`
}

type GetStudentData struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStudentData"`

	UserSessionVO *UserSessionVO      `xml:"userSessionVO,omitempty"`
	StudentIDs    []int64             `xml:"studentIDs,omitempty"`
	Qil           *QueryIncludeListVO `xml:"qil,omitempty"`
}

type GetStudentDataResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStudentDataResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type GetStudentPhoto struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStudentPhoto"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	StudentID     int64          `xml:"studentID,omitempty"`
}

type GetStudentPhotoResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStudentPhotoResponse"`

	Return_ []byte `xml:"return,omitempty"`
}

type RecoverUsername struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd recoverUsername"`

	EmailAddress string `xml:"emailAddress,omitempty"`
}

type RecoverUsernameResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd recoverUsernameResponse"`

	Return_ *MessageVO `xml:"return,omitempty"`
}

type RecoverPassword struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd recoverPassword"`

	UserType      int32  `xml:"userType,omitempty"`
	UserName      string `xml:"userName,omitempty"`
	RecoveryToken string `xml:"recoveryToken,omitempty"`
	NewPassword   string `xml:"newPassword,omitempty"`
}

type RecoverPasswordResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd recoverPasswordResponse"`

	Return_ *PasswordResetVO `xml:"return,omitempty"`
}

type LinkDeviceTokenToUser struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd linkDeviceTokenToUser"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	DeviceToken   string         `xml:"deviceToken,omitempty"`
}

type LinkDeviceTokenToUserResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd linkDeviceTokenToUserResponse"`

	Return_ *MessageVO `xml:"return,omitempty"`
}

type StoreCourseRequests struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd storeCourseRequests"`

	UserSessionVO       *UserSessionVO          `xml:"userSessionVO,omitempty"`
	StudentId           int64                   `xml:"studentId,omitempty"`
	CourseRequestGroups []*CourseRequestGroupVO `xml:"courseRequestGroups,omitempty"`
}

type StoreCourseRequestsResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd storeCourseRequestsResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type GetCredentialComplexityRules struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getCredentialComplexityRules"`

	UserType int32 `xml:"userType,omitempty"`
}

type GetCredentialComplexityRulesResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getCredentialComplexityRulesResponse"`

	Return_ *CredentialComplexityRulesVO `xml:"return,omitempty"`
}

type StoreNotificationSettings struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd storeNotificationSettings"`

	UserSessionVO *UserSessionVO          `xml:"userSessionVO,omitempty"`
	Ns            *NotificationSettingsVO `xml:"ns,omitempty"`
}

type StoreNotificationSettingsResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd storeNotificationSettingsResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type LogoutAndDelinkDeviceToken struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd logoutAndDelinkDeviceToken"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	DeviceToken   string         `xml:"deviceToken,omitempty"`
}

type LogoutAndDelinkDeviceTokenResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd logoutAndDelinkDeviceTokenResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type LoginToPublicPortal struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd loginToPublicPortal"`

	Username string `xml:"username,omitempty"`
	Password string `xml:"password,omitempty"`
}

type LoginToPublicPortalResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd loginToPublicPortalResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type GetStartStopTimeForAllSections struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStartStopTimeForAllSections"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	StudentIDs    []int64        `xml:"studentIDs,omitempty"`
	Month         int32          `xml:"month,omitempty"`
	Year          int32          `xml:"year,omitempty"`
}

type GetStartStopTimeForAllSectionsResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getStartStopTimeForAllSectionsResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type GetAllCourseRequests struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getAllCourseRequests"`

	UserSessionVO *UserSessionVO `xml:"userSessionVO,omitempty"`
	StudentId     int64          `xml:"studentId,omitempty"`
}

type GetAllCourseRequestsResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd getAllCourseRequestsResponse"`

	Return_ *ResultsVO `xml:"return,omitempty"`
}

type SendPasswordRecoveryEmail struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd sendPasswordRecoveryEmail"`

	UserType     int32  `xml:"userType,omitempty"`
	UserName     string `xml:"userName,omitempty"`
	EmailAddress string `xml:"emailAddress,omitempty"`
}

type SendPasswordRecoveryEmailResponse struct {
	XMLName xml.Name `xml:"http://publicportal.rest.powerschool.pearson.com/xsd sendPasswordRecoveryEmailResponse"`

	Return_ *MessageVO `xml:"return,omitempty"`
}

type BaseResultsVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd BaseResultsVO"`

	MessageVOs []*MessageVO `xml:"messagesVO,omitempty"`
}

type MessageVO struct {
	XMLName xml.Name `xml:"messageVOs"`

	Description string `xml:"description,omitempty"`
	Id          string `xml:"id,omitempty"`
	MsgCode     int32  `xml:"msgCode,omitempty"`
	Title       string `xml:"title,omitempty"`
}

type ResultsVO struct {
	XMLName xml.Name `xml:"return"`

	//	*BaseResultsVO
	MessageVOs []*MessageVO `xml:"messageVOs,omitempty"`

	CourseRequestGroupsVOs []*CourseRequestGroupVO `xml:"courseRequestGroupsVOs,omitempty"`
	CourseRequestRulesVO   *CourseRequestRulesVO   `xml:"courseRequestRulesVO,omitempty"`
	StudentDataVOs         []*StudentDataVO        `xml:"studentDataVOs,omitempty"`
	UserSessionVO          *UserSessionVO          `xml:"userSessionVO,omitempty"`
}

type CourseRequestGroupVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd CourseRequestGroupVO"`

	Courses        []*CourseRequestVO `xml:"courses,omitempty"`
	Description    string             `xml:"description,omitempty"`
	EmptyAdvice    string             `xml:"emptyAdvice,omitempty"`
	GradeLevel     int32              `xml:"gradeLevel,omitempty"`
	Id             int64              `xml:"id,omitempty"`
	ItemType       string             `xml:"itemType,omitempty"`
	MaxCourseCount float32            `xml:"maxCourseCount,omitempty"`
	MinCourseCount float32            `xml:"minCourseCount,omitempty"`
	Name           string             `xml:"name,omitempty"`
	RequestType    string             `xml:"requestType,omitempty"`
	Requests       []*CourseRequestVO `xml:"requests,omitempty"`
	SchoolId       int64              `xml:"schoolId,omitempty"`
	SortOrder      int32              `xml:"sortOrder,omitempty"`
	YearId         int32              `xml:"yearId,omitempty"`
}

type CourseRequestVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd CourseRequestVO"`

	CourseName   string  `xml:"courseName,omitempty"`
	CourseNumber string  `xml:"courseNumber,omitempty"`
	CreditHours  float32 `xml:"creditHours,omitempty"`
}

type CourseRequestRulesVO struct {
	XMLName xml.Name `xml:"courseRequestRulesVO"`

	Description string  `xml:"description,omitempty"`
	MaxCredits  float64 `xml:"maxCredits,omitempty"`
	MinCredits  float64 `xml:"minCredits,omitempty"`
}

type StudentDataVO struct {
	Activities             []*ActivityVO           `xml:"activities,omitempty"`
	ArchivedFinalGrades    []*ArchivedFinalGradeVO `xml:"archivedFinalGrades,omitempty"`
	AssignmentCategories   []*AsmtCatVO            `xml:"AsmtCatVO,omitempty"`
	AssignmentScores       []*AssignmentScoreVO    `xml:"AssignmentScoreVO,omitempty"`
	Assignments            []*AssignmentVO         `xml:"AssignmentVO,omitempty"`
	Attendance             []*AttendanceVO         `xml:"AttendancVO,omitempty"`
	AttendanceCodes        []*AttendanceCodeVO     `xml:"AttendanceCodeVO,omitempty"`
	Bulletins              []*BulletinLite         `xml:"bulletins,omitempty"`
	CitizenCodes           []*CitizenCodeVO        `xml:"citizenCodes,omitempty"`
	CitizenGrades          []*CitizenGradeVO       `xml:"citizenGrades,omitempty"`
	CourseRequests         []*CourseRequestVO      `xml:"courseRequests,omitempty"`
	Enrollments            []*SectionEnrollmentVO  `xml:"enrollments,omitempty"`
	FeeBalance             *FeeBalanceVO           `xml:"feeBalance,omitempty"`
	FeeTransactions        []*FeeTransactionVO     `xml:"feeTransactions,omitempty"`
	FeeTypes               []*FeeTypeVO            `xml:"feeTypes,omitempty"`
	FinalGrades            []*FinalGradeVO         `xml:"finalGradeVO,omitempty"`
	GradeScales            []*GradeScaleVO         `xml:"gradeScales,omitempty"`
	LunchTransactions      []*LunchTransactionVO   `xml:"lunchTransactions,omitempty"`
	NotInSessionDays       []*NotInSessionDayVO    `xml:"notInSessionDayVO,omitempty"`
	NotificationSettingsVO *NotificationSettingsVO `xml:"notificationSettingsVO,omitempty"`
	Periods                []*PeriodVO             `xml:"PeriodVO,omitempty"`
	RemoteSchools          []*SchoolVO             `xml:"remoteSchools,omitempty"`
	ReportingTerms         []*ReportingTermVO      `xml:"ReportingTermsVO,omitempty"`
	Schools                []*SchoolVO             `xml:"schoolVO,omitempty"`
	Sections               []*SectionVO            `xml:"SectionVO,omitempty"`
	Standards              []*StandardVO           `xml:"standards,omitempty"`
	StandardsGrades        []*StandardGradeVO      `xml:"standardsGrades,omitempty"`
	Student                *StudentVO              `xml:"student,omitempty"`
	StudentDcid            int64                   `xml:"studentDcid,omitempty"`
	StudentId              int64                   `xml:"studentId,omitempty"`
	Teachers               []*TeacherVO            `xml:"TeacherVO,omitempty"`
	Terms                  []*TermVO               `xml:"TermVO,omitempty"`
	YearId                 int32                   `xml:"yearId,omitempty"`
}

type ActivityVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd ActivityVO"`

	Category string `xml:"category,omitempty"`
	Id       int64  `xml:"id,omitempty"`
	Name     string `xml:"name,omitempty"`
	Required bool   `xml:"required,omitempty"`
}

type FinalGradeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd FinalGradeVO"`

	CommentValue    string    `xml:"commentValue,omitempty"`
	DateStored      time.Time `xml:"dateStored,omitempty"`
	Grade           string    `xml:"grade,omitempty"`
	Id              int64     `xml:"id,omitempty"`
	Percent         float64   `xml:"percent,omitempty"`
	ReportingTermId int64     `xml:"reportingTermId,omitempty"`
	Sectionid       int64     `xml:"sectionid,omitempty"`
	StoreType       int32     `xml:"storeType,omitempty"`
}

type ArchivedFinalGradeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd ArchivedFinalGradeVO"`

	*FinalGradeVO

	CourseName    string    `xml:"courseName,omitempty"`
	CourseNumber  string    `xml:"courseNumber,omitempty"`
	SchoolId      int64     `xml:"schoolId,omitempty"`
	SortOrder     int32     `xml:"sortOrder,omitempty"`
	StoreCode     string    `xml:"storeCode,omitempty"`
	TeacherName   string    `xml:"teacherName,omitempty"`
	TermEndDate   time.Time `xml:"termEndDate,omitempty"`
	TermId        int64     `xml:"termId,omitempty"`
	TermStartDate time.Time `xml:"termStartDate,omitempty"`
	YearId        int64     `xml:"yearId,omitempty"`
}

type AsmtCatVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd AsmtCatVO"`

	Abbreviation  string `xml:"abbreviation,omitempty"`
	Description   string `xml:"description,omitempty"`
	GradeBookType int32  `xml:"gradeBookType,omitempty"`
	Id            int64  `xml:"id,omitempty"`
	Name          string `xml:"name,omitempty"`
}

type AssignmentScoreVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd AssignmentScoreVO"`

	AssignmentId  int64  `xml:"assignmentId,omitempty"`
	Collected     bool   `xml:"collected,omitempty"`
	Comment       string `xml:"comment,omitempty"`
	Exempt        bool   `xml:"exempt,omitempty"`
	GradeBookType int32  `xml:"gradeBookType,omitempty"`
	Id            int64  `xml:"id,omitempty"`
	Late          bool   `xml:"late,omitempty"`
	LetterGrade   string `xml:"letterGrade,omitempty"`
	Missing       bool   `xml:"missing,omitempty"`
	Percent       string `xml:"percent,omitempty"`
	Score         string `xml:"score,omitempty"`
	Scoretype     int32  `xml:"scoretype,omitempty"`
}

type AssignmentVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd AssignmentVO"`

	Abbreviation          string    `xml:"abbreviation,omitempty"`
	AdditionalCategoryIds []int32   `xml:"additionalCategoryIds,omitempty"`
	Assignmentid          int64     `xml:"assignmentid,omitempty"`
	CategoryId            int32     `xml:"categoryId,omitempty"`
	Description           string    `xml:"description,omitempty"`
	DueDate               time.Time `xml:"dueDate,omitempty"`
	GradeBookType         int32     `xml:"gradeBookType,omitempty"`
	Id                    int64     `xml:"id,omitempty"`
	Includeinfinalgrades  int32     `xml:"includeinfinalgrades,omitempty"`
	Name                  string    `xml:"name,omitempty"`
	Pointspossible        float64   `xml:"pointspossible,omitempty"`
	PublishDaysBeforeDue  int32     `xml:"publishDaysBeforeDue,omitempty"`
	PublishState          int32     `xml:"publishState,omitempty"`
	Publishonspecificdate time.Time `xml:"publishonspecificdate,omitempty"`
	Publishscores         int32     `xml:"publishscores,omitempty"`
	SectionDcid           int64     `xml:"sectionDcid,omitempty"`
	Sectionid             int64     `xml:"sectionid,omitempty"`
	Type_                 int32     `xml:"type,omitempty"`
	Weight                float64   `xml:"weight,omitempty"`
}

type AttendanceVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd AttendanceVO"`

	AdaValueCode    float64   `xml:"adaValueCode,omitempty"`
	AdaValueTime    float64   `xml:"adaValueTime,omitempty"`
	AdmValue        float64   `xml:"admValue,omitempty"`
	AttCodeid       int64     `xml:"attCodeid,omitempty"`
	AttComment      string    `xml:"attComment,omitempty"`
	AttDate         time.Time `xml:"attDate,omitempty"`
	AttFlags        int32     `xml:"attFlags,omitempty"`
	AttInterval     int32     `xml:"attInterval,omitempty"`
	AttModeCode     string    `xml:"attModeCode,omitempty"`
	Ccid            int64     `xml:"ccid,omitempty"`
	Id              int64     `xml:"id,omitempty"`
	Periodid        int64     `xml:"periodid,omitempty"`
	Schoolid        int64     `xml:"schoolid,omitempty"`
	Studentid       int64     `xml:"studentid,omitempty"`
	TotalMinutes    float64   `xml:"totalMinutes,omitempty"`
	TransactionType string    `xml:"transactionType,omitempty"`
	Yearid          int32     `xml:"yearid,omitempty"`
}

type AttendanceCodeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd AttendanceCodeVO"`

	AttCode     string `xml:"attCode,omitempty"`
	CodeType    int32  `xml:"codeType,omitempty"`
	Description string `xml:"description,omitempty"`
	Id          int64  `xml:"id,omitempty"`
	Schoolid    int64  `xml:"schoolid,omitempty"`
	Sortorder   int32  `xml:"sortorder,omitempty"`
	Yearid      int32  `xml:"yearid,omitempty"`
}

type CitizenCodeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd CitizenCodeVO"`

	CodeName    string `xml:"codeName,omitempty"`
	Description string `xml:"description,omitempty"`
	Id          int64  `xml:"id,omitempty"`
	SortOrder   int32  `xml:"sortOrder,omitempty"`
}

type CitizenGradeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd CitizenGradeVO"`

	CodeId          int64 `xml:"codeId,omitempty"`
	ReportingTermId int64 `xml:"reportingTermId,omitempty"`
	SectionId       int64 `xml:"sectionId,omitempty"`
	StoreType       int32 `xml:"storeType,omitempty"`
}

type SectionEnrollmentVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd SectionEnrollmentVO"`

	EndDate      time.Time `xml:"endDate,omitempty"`
	EnrollStatus int32     `xml:"enrollStatus,omitempty"`
	Id           int64     `xml:"id,omitempty"`
	StartDate    time.Time `xml:"startDate,omitempty"`
}

type FeeBalanceVO struct {
	XMLName xml.Name `xml:"feeBalance"`

	Balance  float64 `xml:"balance,omitempty"`
	Credit   float64 `xml:"credit,omitempty"`
	Debit    float64 `xml:"debit,omitempty"`
	Id       int64   `xml:"id,omitempty"`
	Schoolid int64   `xml:"schoolid,omitempty"`
	Yearid   int32   `xml:"yearid,omitempty"`
}

type FeeTransactionVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd FeeTransactionVO"`

	Adjustment         float64   `xml:"adjustment,omitempty"`
	CourseName         string    `xml:"courseName,omitempty"`
	CourseNumber       string    `xml:"courseNumber,omitempty"`
	Creationdate       time.Time `xml:"creationdate,omitempty"`
	DateValue          time.Time `xml:"dateValue,omitempty"`
	DepartmentName     string    `xml:"departmentName,omitempty"`
	Description        string    `xml:"description,omitempty"`
	FeeAmount          float64   `xml:"feeAmount,omitempty"`
	FeeBalance         float64   `xml:"feeBalance,omitempty"`
	FeeCategoryName    string    `xml:"feeCategoryName,omitempty"`
	FeePaid            float64   `xml:"feePaid,omitempty"`
	FeeTypeId          int64     `xml:"feeTypeId,omitempty"`
	FeeTypeName        string    `xml:"feeTypeName,omitempty"`
	Feecharged         float64   `xml:"feecharged,omitempty"`
	GroupTransactionId int64     `xml:"groupTransactionId,omitempty"`
	Id                 int64     `xml:"id,omitempty"`
	Modificationdate   time.Time `xml:"modificationdate,omitempty"`
	Originalfee        float64   `xml:"originalfee,omitempty"`
	Priority           int32     `xml:"priority,omitempty"`
	ProRated           int32     `xml:"proRated,omitempty"`
	SchoolfeeId        int64     `xml:"schoolfeeId,omitempty"`
	Schoolid           int64     `xml:"schoolid,omitempty"`
	Termid             int64     `xml:"termid,omitempty"`
}

type FeeTypeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd FeeTypeVO"`

	Descript        string `xml:"descript,omitempty"`
	FeeCategoryName string `xml:"feeCategoryName,omitempty"`
	Id              int64  `xml:"id,omitempty"`
	SchoolNumber    int32  `xml:"schoolNumber,omitempty"`
	Sort            int32  `xml:"sort,omitempty"`
	Title           string `xml:"title,omitempty"`
}

type GradeScaleVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd GradeScaleVO"`

	Description      string              `xml:"description,omitempty"`
	GradeBookType    int32               `xml:"gradeBookType,omitempty"`
	GradeScaleItems  []*GradeScaleItemVO `xml:"gradeScaleItems,omitempty"`
	Id               int64               `xml:"id,omitempty"`
	Name             string              `xml:"name,omitempty"`
	Numeric          int32               `xml:"numeric,omitempty"`
	NumericMax       int32               `xml:"numericMax,omitempty"`
	NumericMin       int32               `xml:"numericMin,omitempty"`
	NumericPrecision int32               `xml:"numericPrecision,omitempty"`
	NumericScale     int32               `xml:"numericScale,omitempty"`
}

type GradeScaleItemVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd GradeScaleItemVO"`

	CutoffPercent     float64 `xml:"cutoffPercent,omitempty"`
	DefaultZeroCutoff bool    `xml:"defaultZeroCutoff,omitempty"`
	Description       string  `xml:"description,omitempty"`
	GradeBookType     int32   `xml:"gradeBookType,omitempty"`
	GradeLabel        string  `xml:"gradeLabel,omitempty"`
	Id                int64   `xml:"id,omitempty"`
	PercentValue      float64 `xml:"percentValue,omitempty"`
	PointsValue       float64 `xml:"pointsValue,omitempty"`
	SortOrder         int32   `xml:"sortOrder,omitempty"`
}

type LunchTransactionVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd LunchTransactionVO"`

	Cash        float64   `xml:"cash,omitempty"`
	Credit      float64   `xml:"credit,omitempty"`
	DateValue   time.Time `xml:"dateValue,omitempty"`
	Debit       float64   `xml:"debit,omitempty"`
	Description string    `xml:"description,omitempty"`
	Id          int64     `xml:"id,omitempty"`
	Mealprice   float64   `xml:"mealprice,omitempty"`
	Neteffect   float64   `xml:"neteffect,omitempty"`
	Time        int32     `xml:"time,omitempty"`
}

type NotInSessionDayVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd NotInSessionDayVO"`

	CalType      string    `xml:"calType,omitempty"`
	CalendarDay  time.Time `xml:"calendarDay,omitempty"`
	Description  string    `xml:"description,omitempty"`
	Id           int64     `xml:"id,omitempty"`
	SchoolNumber int64     `xml:"schoolNumber,omitempty"`
}

type NotificationSettingsVO struct {
	XMLName xml.Name `xml:"notificationSettingsVO"`
	ApplyToAllStudents  bool     `xml:"applyToAllStudents,omitempty"`
	BalanceAlerts       bool     `xml:"balanceAlerts,omitempty"`
	DetailedAssignments bool     `xml:"detailedAssignments,omitempty"`
	DetailedAttendance  bool     `xml:"detailedAttendance,omitempty"`
	EmailAddresses      []string `xml:"emailAddresses,omitempty"`
	Frequency           int32    `xml:"frequency,omitempty"`
	GradeAndAttSummary  bool     `xml:"gradeAndAttSummary,omitempty"`
	GuardianStudentId   int64    `xml:"guardianStudentId,omitempty"`
	MainEmail           string   `xml:"mainEmail,omitempty"`
	SchoolAnnouncements bool     `xml:"schoolAnnouncements,omitempty"`
	SendNow             bool     `xml:"sendNow,omitempty"`
}

type PeriodVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd PeriodVO"`

	Abbreviation string `xml:"abbreviation,omitempty"`
	Id           int64  `xml:"id,omitempty"`
	Name         string `xml:"name,omitempty"`
	PeriodNumber int32  `xml:"periodNumber,omitempty"`
	Schoolid     int64  `xml:"schoolid,omitempty"`
	SortOrder    int32  `xml:"sortOrder,omitempty"`
	Yearid       int32  `xml:"yearid,omitempty"`
}

type SchoolVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd SchoolVO"`

	Abbreviation          string              `xml:"abbreviation,omitempty"`
	Address               string              `xml:"address,omitempty"`
	DisabledFeatures      *DisabledFeaturesVO `xml:"disabledFeatures,omitempty"`
	HighGrade             int32               `xml:"highGrade,omitempty"`
	LowGrade              int32               `xml:"lowGrade,omitempty"`
	MapMimeType           string              `xml:"mapMimeType,omitempty"`
	Name                  string              `xml:"name,omitempty"`
	SchoolDisabled        bool                `xml:"schoolDisabled,omitempty"`
	SchoolDisabledMessage string              `xml:"schoolDisabledMessage,omitempty"`
	SchoolDisabledTitle   string              `xml:"schoolDisabledTitle,omitempty"`
	SchoolId              int64               `xml:"schoolId,omitempty"`
	SchoolMapModifiedDate time.Time           `xml:"schoolMapModifiedDate,omitempty"`
	SchoolNumber          int64               `xml:"schoolNumber,omitempty"`
	Schooladdress         string              `xml:"schooladdress,omitempty"`
	Schoolcity            string              `xml:"schoolcity,omitempty"`
	Schoolcountry         string              `xml:"schoolcountry,omitempty"`
	Schoolfax             string              `xml:"schoolfax,omitempty"`
	Schoolphone           string              `xml:"schoolphone,omitempty"`
	Schoolstate           string              `xml:"schoolstate,omitempty"`
	Schoolzip             string              `xml:"schoolzip,omitempty"`
}

type DisabledFeaturesVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd disabledFeatures"`

	Activities  bool `xml:"activities,omitempty"`
	Assignments bool `xml:"assignments,omitempty"`
	Attendance  bool `xml:"attendance,omitempty"`
	Citizenship bool `xml:"citizenship,omitempty"`
	CurrentGpa  bool `xml:"currentGpa,omitempty"`
	Emailalerts bool `xml:"emailalerts,omitempty"`
	Fees        bool `xml:"fees,omitempty"`
	FinalGrades bool `xml:"finalGrades,omitempty"`
	Meals       bool `xml:"meals,omitempty"`
	Standards   bool `xml:"standards,omitempty"`
}

type ReportingTermVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd ReportingTermVO"`

	Abbreviation     string    `xml:"abbreviation,omitempty"`
	EndDate          time.Time `xml:"endDate,omitempty"`
	Id               int64     `xml:"id,omitempty"`
	Schoolid         int64     `xml:"schoolid,omitempty"`
	SendingGrades    bool      `xml:"sendingGrades,omitempty"`
	SortOrder        int32     `xml:"sortOrder,omitempty"`
	StartDate        time.Time `xml:"startDate,omitempty"`
	SuppressGrades   bool      `xml:"suppressGrades,omitempty"`
	SuppressPercents bool      `xml:"suppressPercents,omitempty"`
	Termid           int64     `xml:"termid,omitempty"`
	Title            string    `xml:"title,omitempty"`
	Yearid           int64     `xml:"yearid,omitempty"`
}

type SectionVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd SectionVO"`

	CourseCode        string                 `xml:"courseCode,omitempty"`
	Dcid              int64                  `xml:"dcid,omitempty"`
	Description       string                 `xml:"description,omitempty"`
	Enrollments       []*SectionEnrollmentVO `xml:"enrollments,omitempty"`
	Expression        string                 `xml:"expression,omitempty"`
	GradeBookType     int32                  `xml:"gradeBookType,omitempty"`
	Id                int64                  `xml:"id,omitempty"`
	PeriodSort        int32                  `xml:"periodSort,omitempty"`
	RoomName          string                 `xml:"roomName,omitempty"`
	SchoolCourseTitle string                 `xml:"schoolCourseTitle,omitempty"`
	SchoolNumber      int64                  `xml:"schoolNumber,omitempty"`
	SectionNum        string                 `xml:"sectionNum,omitempty"`
	StartStopDates    []*StartStopDateVO     `xml:"startStopDates,omitempty"`
	TeacherID         int64                  `xml:"teacherID,omitempty"`
	TermID            int64                  `xml:"termID,omitempty"`
}

type StartStopDateVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd StartStopDateVO"`

	SectionEnrollmentId int64     `xml:"sectionEnrollmentId,omitempty"`
	Start               time.Time `xml:"start,omitempty"`
	Stop                time.Time `xml:"stop,omitempty"`
}

type StandardVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd StandardVO"`

	Description      string `xml:"description,omitempty"`
	GradeBookType    int32  `xml:"gradeBookType,omitempty"`
	GradeScaleID     int64  `xml:"gradeScaleID,omitempty"`
	Id               int64  `xml:"id,omitempty"`
	Identifier       string `xml:"identifier,omitempty"`
	Name             string `xml:"name,omitempty"`
	ParentStandardID int64  `xml:"parentStandardID,omitempty"`
	SortOrder        int32  `xml:"sortOrder,omitempty"`
}

type StandardGradeVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd StandardGradeVO"`

	Comment            string    `xml:"comment,omitempty"`
	CommentLastUpdated time.Time `xml:"commentLastUpdated,omitempty"`
	Exempt             int32     `xml:"exempt,omitempty"`
	GradeBookType      int32     `xml:"gradeBookType,omitempty"`
	GradeEntered       string    `xml:"gradeEntered,omitempty"`
	GradeLastUpdated   time.Time `xml:"gradeLastUpdated,omitempty"`
	GradeType          int32     `xml:"gradeType,omitempty"`
	Id                 int64     `xml:"id,omitempty"`
	Late               int32     `xml:"late,omitempty"`
	Missing            int32     `xml:"missing,omitempty"`
	ReportingTermId    int64     `xml:"reportingTermId,omitempty"`
	SectionDcid        int64     `xml:"sectionDcid,omitempty"`
	SectionId          int64     `xml:"sectionId,omitempty"`
	StandardId         int64     `xml:"standardId,omitempty"`
}

type StudentVO struct {
	XMLName xml.Name `xml:"student"`

	CurrentGPA             string    `xml:"currentGPA,omitempty"`
	CurrentMealBalance     float64   `xml:"currentMealBalance,omitempty"`
	CurrentTerm            string    `xml:"currentTerm,omitempty"`
	Dcid                   int64     `xml:"dcid,omitempty"`
	Dob                    time.Time `xml:"dob,omitempty"`
	Ethnicity              string    `xml:"ethnicity,omitempty"`
	FirstName              string    `xml:"firstName,omitempty"`
	Gender                 string    `xml:"gender,omitempty"`
	GradeLevel             int32     `xml:"gradeLevel,omitempty"`
	GuardianAccessDisabled bool      `xml:"guardianAccessDisabled,omitempty"`
	Id                     int64     `xml:"id,omitempty"`
	LastName               string    `xml:"lastName,omitempty"`
	MiddleName             string    `xml:"middleName,omitempty"`
	PhotoDate              time.Time `xml:"photoDate,omitempty"`
	StartingMealBalance    float64   `xml:"startingMealBalance,omitempty"`
}

type TeacherVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd TeacherVO"`

	Email       string `xml:"email,omitempty"`
	FirstName   string `xml:"firstName,omitempty"`
	Id          int64  `xml:"id,omitempty"`
	LastName    string `xml:"lastName,omitempty"`
	SchoolPhone string `xml:"schoolPhone,omitempty"`
}

type TermVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd TermVO"`

	Abbrev       string    `xml:"abbrev,omitempty"`
	EndDate      time.Time `xml:"endDate,omitempty"`
	Id           int64     `xml:"id,omitempty"`
	ParentTermId int64     `xml:"parentTermId,omitempty"`
	SchoolNumber string    `xml:"schoolNumber,omitempty"`
	StartDate    time.Time `xml:"startDate,omitempty"`
	Title        string    `xml:"title,omitempty"`
}

type UserSessionVO struct {
	XMLName xml.Name `xml:"userSessionVO"`

	Locale            *Locale     `xml:"locale,omitempty"`
	ServerCurrentTime time.Time   `xml:"serverCurrentTime,omitempty"`
	ServerInfo        *ServerInfo `xml:"serverInfo,omitempty"`
	ServiceTicket     string      `xml:"serviceTicket,omitempty"`
	StudentIDs        []int32     `xml:"studentIDs,omitempty"`
	UserId            int64       `xml:"userId,omitempty"`
	UserType          int32       `xml:"userType,omitempty"`
}

type QueryIncludeListVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd qil"`

	Includes []int32 `xml:"includes,omitempty"`
}

type PasswordResetVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd PasswordResetVO"`

	*BaseResultsVO

	MinPasswordLength int32  `xml:"minPasswordLength,omitempty"`
	ServiceTicket     string `xml:"serviceTicket,omitempty"`
	Successful        bool   `xml:"successful,omitempty"`
}

type CredentialComplexityRulesVO struct {
	XMLName xml.Name `xml:"http://vo.rest.powerschool.pearson.com/xsd CredentialComplexityRulesVO"`

	*BaseResultsVO

	LettersAndNumRequired    bool  `xml:"lettersAndNumRequired,omitempty"`
	MixOfCaseRequired        bool  `xml:"mixOfCaseRequired,omitempty"`
	RequiredCharacterCount   int32 `xml:"requiredCharacterCount,omitempty"`
	SpecialCharacterRequired bool  `xml:"specialCharacterRequired,omitempty"`
	Successful               bool  `xml:"successful,omitempty"`
}

type BulletinLite struct {
	XMLName xml.Name `xml:"http://model.rest.powerschool.pearson.com/xsd BulletinLite"`

	Audience  int64     `xml:"audience,omitempty"`
	Body      string    `xml:"body,omitempty"`
	EndDate   time.Time `xml:"endDate,omitempty"`
	Id        int64     `xml:"id,omitempty"`
	Name      string    `xml:"name,omitempty"`
	SchoolId  int64     `xml:"schoolId,omitempty"`
	SortOrder int32     `xml:"sortOrder,omitempty"`
	StartDate time.Time `xml:"startDate,omitempty"`
}

type ServerInfo struct {
	XMLName xml.Name `xml:"serverInfo"`

	ApiVersion                  string    `xml:"apiVersion,omitempty"`
	DayLightSavings             int32     `xml:"dayLightSavings,omitempty"`
	ParentSAMLEndPoint          string    `xml:"parentSAMLEndPoint,omitempty"`
	PublicPortalDisabled        bool      `xml:"publicPortalDisabled,omitempty"`
	PublicPortalDisabledMessage string    `xml:"publicPortalDisabledMessage,omitempty"`
	RawOffset                   int32     `xml:"rawOffset,omitempty"`
	ServerTime                  time.Time `xml:"serverTime,omitempty"`
	StudentSAMLEndPoint         string    `xml:"studentSAMLEndPoint,omitempty"`
	TeacherSAMLEndPoint         string    `xml:"teacherSAMLEndPoint,omitempty"`
	TimeZoneName                string    `xml:"timeZoneName,omitempty"`
}

type PublicPortalServiceJSONPortType struct {
	client *SOAPClient
}

func NewPublicPortalServiceJSONPortType(url string, tls bool, auth *DigestAuth) *PublicPortalServiceJSONPortType {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &PublicPortalServiceJSONPortType{
		client: client,
	}
}

func (service *PublicPortalServiceJSONPortType) GetCredentialComplexityRules(request *GetCredentialComplexityRules) (*GetCredentialComplexityRulesResponse, error) {
	response := new(GetCredentialComplexityRulesResponse)
	err := service.client.Call("urn:getCredentialComplexityRules", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) LogoutAndDelinkDeviceToken(request *LogoutAndDelinkDeviceToken) (*LogoutAndDelinkDeviceTokenResponse, error) {
	response := new(LogoutAndDelinkDeviceTokenResponse)
	err := service.client.Call("urn:logoutAndDelinkDeviceToken", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) GetStudentData(request *GetStudentData) (*GetStudentDataResponse, error) {
	response := new(GetStudentDataResponse)
	err := service.client.Call("urn:getStudentData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) Login(request *Login) (*LoginResponse, error) {
	response := new(LoginResponse)
	err := service.client.Call("urn:login", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) SendPasswordRecoveryEmail(request *SendPasswordRecoveryEmail) (*SendPasswordRecoveryEmailResponse, error) {
	response := new(SendPasswordRecoveryEmailResponse)
	err := service.client.Call("urn:sendPasswordRecoveryEmail", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) Logout(request *Logout) (*LogoutResponse, error) {
	response := new(LogoutResponse)
	err := service.client.Call("urn:logout", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) LoginToPublicPortal(request *LoginToPublicPortal) (*LoginToPublicPortalResponse, error) {
	response := new(LoginToPublicPortalResponse)
	err := service.client.Call("urn:loginToPublicPortal", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) RecoverUsername(request *RecoverUsername) (*RecoverUsernameResponse, error) {
	response := new(RecoverUsernameResponse)
	err := service.client.Call("urn:recoverUsername", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) LinkDeviceTokenToUser(request *LinkDeviceTokenToUser) (*LinkDeviceTokenToUserResponse, error) {
	response := new(LinkDeviceTokenToUserResponse)
	err := service.client.Call("urn:linkDeviceTokenToUser", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) GetStudentPhoto(request *GetStudentPhoto) (*GetStudentPhotoResponse, error) {
	response := new(GetStudentPhotoResponse)
	err := service.client.Call("urn:getStudentPhoto", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) RecoverPassword(request *RecoverPassword) (*RecoverPasswordResponse, error) {
	response := new(RecoverPasswordResponse)
	err := service.client.Call("urn:recoverPassword", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) GetSchoolMapBySchoolNumber(request *GetSchoolMapBySchoolNumber) (*GetSchoolMapBySchoolNumberResponse, error) {
	response := new(GetSchoolMapBySchoolNumberResponse)
	err := service.client.Call("urn:getSchoolMapBySchoolNumber", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) StoreNotificationSettings(request *StoreNotificationSettings) (*StoreNotificationSettingsResponse, error) {
	response := new(StoreNotificationSettingsResponse)
	err := service.client.Call("urn:storeNotificationSettings", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) StoreCourseRequests(request *StoreCourseRequests) (*StoreCourseRequestsResponse, error) {
	response := new(StoreCourseRequestsResponse)
	err := service.client.Call("urn:storeCourseRequests", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) GetAllCourseRequests(request *GetAllCourseRequests) (*GetAllCourseRequestsResponse, error) {
	response := new(GetAllCourseRequestsResponse)
	err := service.client.Call("urn:getAllCourseRequests", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PublicPortalServiceJSONPortType) GetStartStopTimeForAllSections(request *GetStartStopTimeForAllSections) (*GetStartStopTimeForAllSectionsResponse, error) {
	response := new(GetStartStopTimeForAllSectionsResponse)
	err := service.client.Call("urn:getStartStopTimeForAllSections", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Header interface{}
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type DigestAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url  string
	tls  bool
	auth *DigestAuth
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls bool, auth *DigestAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
		//Header:        SoapHeader{},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}
	// AUTH
	req, err := http.NewRequest("POST", s.url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	digest := digestParts(resp)
	digest["uri"] = ""
	digest["method"] = "POST"
	digest["username"] = s.auth.Login
	digest["password"] = s.auth.Password
	// SOAP
	req, err = http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}
	req.Header.Set("Authorization", getDigestAuth(digest))
	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuth(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}
