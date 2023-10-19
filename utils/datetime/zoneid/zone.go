package zoneid

import "github.com/qianwj/typed/utils/datetime"

var shortIDs map[string]string

func init() {
	shortIDs = make(map[string]string)
	shortIDs["EST"] = "-05:00"
	shortIDs["HST"] = "-10:00"
	shortIDs["MST"] = "-07:00"
	shortIDs["ACT"] = "Australia/Darwin"
	shortIDs["AET"] = "Australia/Sydney"
	shortIDs["AGT"] = "America/Argentina/Buenos_Aires"
	shortIDs["ART"] = "Africa/Cairo"
	shortIDs["AST"] = "America/Anchorage"
	shortIDs["BET"] = "America/Sao_Paulo"
	shortIDs["BST"] = "Asia/Dhaka"
	shortIDs["CAT"] = "Africa/Harare"
	shortIDs["CNT"] = "America/St_Johns"
	shortIDs["CST"] = "America/Chicago"
	shortIDs["CTT"] = "Asia/Shanghai"
	shortIDs["EAT"] = "Africa/Addis_Ababa"
	shortIDs["ECT"] = "Europe/Paris"
	shortIDs["IET"] = "America/Indiana/Indianapolis"
	shortIDs["IST"] = "Asia/Kolkata"
	shortIDs["JST"] = "Asia/Tokyo"
	shortIDs["MIT"] = "Pacific/Apia"
	shortIDs["NET"] = "Asia/Yerevan"
	shortIDs["NST"] = "Pacific/Auckland"
	shortIDs["PLT"] = "Asia/Karachi"
	shortIDs["PNT"] = "America/Phoenix"
	shortIDs["PRT"] = "America/Puerto_Rico"
	shortIDs["PST"] = "America/Los_Angeles"
	shortIDs["SST"] = "Pacific/Guadalcanal"
	shortIDs["VST"] = "Asia/Ho_Chi_Minh"
}

type ZoneId struct {
}

func From(accessor datetime.TemporalAccessor) *ZoneId {
	return nil
}

func SystemDefault() *ZoneId {
	return nil
}
