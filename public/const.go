package public

const (
	ValidatorKey  = "ValidatorKey"
	TranslatorKey = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfoKey"		// redis save


	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0			// 前缀URL匹配
	HTTPRuleTypeDomain    = 1			// 域名匹配

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal          = "flow_total"
	FlowServicePrefix  = "flow_service_"
	FlowAppPrefix = "flow_app_"

	JwtSignKey = "my_sign_key"
	JwtExpires = 60*60
)
