package hb

type Order struct {
	AmountSpent      float64      `json:"amount_spent"`
	Product          Product      `json:"product"`
	Gamekey          string       `json:"gamekey"`
	UID              string       `json:"uid"`
	AllCouponData    interface{}  `json:"all_coupon_data"`
	Created          string       `json:"created"`
	MissedCredit     interface{}  `json:"missed_credit"`
	Subproducts      []Subproduct `json:"subproducts"`
	TotalChoices     int64        `json:"total_choices"`
	TpkdDict         TpkdDict     `json:"tpkd_dict"`
	ChoicesRemaining int64        `json:"choices_remaining"`
	Currency         string       `json:"currency"`
	IsGiftee         bool         `json:"is_giftee"`
	HasWallet        bool         `json:"has_wallet"`
	Claimed          bool         `json:"claimed"`
	Total            float64      `json:"total"`
	WalletCredit     interface{}  `json:"wallet_credit"`
	PathIDS          []string     `json:"path_ids"`
}

type Product struct {
	Category           string     `json:"category"`
	MachineName        string     `json:"machine_name"`
	EmptyTpkds         EmptyTpkds `json:"empty_tpkds"`
	PostPurchaseText   string     `json:"post_purchase_text"`
	HumanName          string     `json:"human_name"`
	PartialGiftEnabled bool       `json:"partial_gift_enabled"`
}

type EmptyTpkds struct{}

type Subproduct struct {
	MachineName               string      `json:"machine_name"`
	URL                       string      `json:"url"`
	Downloads                 []Download  `json:"downloads"`
	LibraryFamilyName         *string     `json:"library_family_name"`
	Payee                     Payee       `json:"payee"`
	HumanName                 string      `json:"human_name"`
	CustomDownloadPageBoxCSS  interface{} `json:"custom_download_page_box_css"`
	CustomDownloadPageBoxHTML *string     `json:"custom_download_page_box_html"`
	Icon                      string      `json:"icon"`
}

type Download struct {
	MachineName           string           `json:"machine_name"`
	Platform              string           `json:"platform"`
	DownloadStruct        []DownloadStruct `json:"download_struct"`
	OptionsDict           EmptyTpkds       `json:"options_dict"`
	DownloadIdentifier    string           `json:"download_identifier"`
	AndroidAppOnly        bool             `json:"android_app_only"`
	DownloadVersionNumber interface{}      `json:"download_version_number"`
}

type DownloadStruct struct {
	Name             string  `json:"name"`
	URL              URL     `json:"url"`
	HumanSize        string  `json:"human_size"`
	UsesKindleSender *bool   `json:"uses_kindle_sender,omitempty"`
	FileSize         int64   `json:"file_size"`
	Small            int64   `json:"small"`
	MD5              string  `json:"md5"`
	SHA1             string  `json:"sha1,omitempty"`
	UploadedAt       *string `json:"uploaded_at,omitempty"`
	Timestamp        *int64  `json:"timestamp,omitempty"`
}

type URL struct {
	Web        string `json:"web"`
	Bittorrent string `json:"bittorrent"`
}

type Payee struct {
	HumanName   string `json:"human_name"`
	MachineName string `json:"machine_name"`
}

type TpkdDict struct {
	AllTpks []interface{} `json:"all_tpks"`
}
