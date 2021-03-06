package main

type OneLogin struct {
	AccountID                       int    `json:"account_id,omitempty"`
	ActorSystem                     string `json:"actor_system,omitempty"`
	ActorUserID                     int    `json:"actor_user_id,omitempty"`
	ActorUserName                   string `json:"actor_user_name,omitempty"`
	AdcID                           int    `json:"adc_id,omitempty"`
	AdcName                         string `json:"adc_name,omitempty"`
	APICredentialName               string `json:"api_credential_name,omitempty"`
	AppID                           int    `json:"app_id,omitempty"`
	AppName                         string `json:"app_name,omitempty"`
	AssumedBySuperadminOrReseller   string `json:"assumed_by_superadmin_or_reseller,omitempty"`
	AssumingActingUserID            int    `json:"assuming_acting_user_id,omitempty"`
	AuthenticationFactorDescription string `json:"authentication_factor_description,omitempty"`
	AuthenticationFactorID          int    `json:"authentication_factor_id,omitempty"`
	AuthenticationFactorType        string `json:"authentication_factor_type,omitempty"`
	BrowserFingerprint              int    `json:"browser_fingerprint,omitempty"`
	CertificateID                   int    `json:"certificate_id,omitempty"`
	CertificateName                 string `json:"certificate_name,omitempty"`
	ClientID                        int    `json:"client_id,omitempty"`
	Create                          struct {
		ID string `json:"ID,omitempty"`
	} `json:"create,omitempty"`
	CustomMessage      string `json:"custom_message,omitempty"`
	DirectoryID        int    `json:"directory_id,omitempty"`
	DirectoryName      string `json:"directory_name,omitempty"`
	DirectorySyncRunID int    `json:"directory_sync_run_id,omitempty"`
	Entity             string `json:"entity,omitempty"`
	ErrorDescription   string `json:"error_description,omitempty"`
	EventTimestamp     string `json:"event_timestamp,omitempty"`
	EventTypeID        int    `json:"event_type_id,omitempty"`
	GroupID            int    `json:"group_id,omitempty"`
	GroupName          string `json:"group_name,omitempty"`
	ImportedUserID     int    `json:"imported_user_id,omitempty"`
	ImportedUserName   string `json:"imported_user_name,omitempty"`
	Ipaddr             string `json:"ipaddr,omitempty"`
	LoginID            int    `json:"login_id,omitempty"`
	LoginName          string `json:"login_name,omitempty"`
	MappingID          int    `json:"mapping_id,omitempty"`
	MappingName        string `json:"mapping_name,omitempty"`
	NoteID             int    `json:"note_id,omitempty"`
	NoteTitle          string `json:"note_title,omitempty"`
	Notes              string `json:"notes,omitempty"`
	ObjectID           int    `json:"object_id,omitempty"`
	OtpDeviceID        int    `json:"otp_device_id,omitempty"`
	OtpDeviceName      string `json:"otp_device_name,omitempty"`
	Param              string `json:"param,omitempty"`
	PolicyID           int    `json:"policy_id,omitempty"`
	PolicyName         string `json:"policy_name,omitempty"`
	PolicyType         string `json:"policy_type,omitempty"`
	PrivilegeID        int    `json:"privilege_id,omitempty"`
	PrivilegeName      string `json:"privilege_name,omitempty"`
	ProxyAgentID       int    `json:"proxy_agent_id,omitempty"`
	ProxyAgentName     string `json:"proxy_agent_name,omitempty"`
	ProxyIP            string `json:"proxy_ip,omitempty"`
	RadiusConfigID     int    `json:"radius_config_id,omitempty"`
	RadiusConfigName   string `json:"radius_config_name,omitempty"`
	Resolution         string `json:"resolution,omitempty"`
	ResolvedAt         string `json:"resolved_at,omitempty"`
	ResolvedByUserID   int    `json:"resolved_by_user_id,omitempty"`
	ResourceTypeID     int    `json:"resource_type_id,omitempty"`
	RiskCookieID       int    `json:"risk_cookie_id,omitempty"`
	RiskReasons        string `json:"risk_reasons,omitempty"`
	RiskScore          int    `json:"risk_score,omitempty"`
	RoleID             string `json:"role_id,omitempty"`
	RoleName           string `json:"role_name,omitempty"`
	ServiceDirectoryID string `json:"service_directory_id,omitempty"`
	Solved             string `json:"solved,omitempty"`
	TaskName           string `json:"task_name,omitempty"`
	TrustedIdpID       string `json:"trusted_idp_id,omitempty"`
	TrustedIdpName     string `json:"trusted_idp_name,omitempty"`
	UserAgent          string `json:"user_agent,omitempty"`
	UserFieldID        string `json:"user_field_id,omitempty"`
	UserFieldName      string `json:"user_field_name,omitempty"`
	UserID             int    `json:"user_id,omitempty"`
	UserName           string `json:"user_name,omitempty"`
	UUID               string `json:"uuid,omitempty"`
}
