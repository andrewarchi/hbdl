package hb

import (
	"encoding/json"
	"fmt"

	"github.com/andrewarchi/unixtime"
)

// GetOrder fetches information on an order.
func (c *Client) GetOrder(gamekey string) (*Order, error) {
	url := fmt.Sprintf("https://www.humblebundle.com/api/v1/order/%s?all_tpkds=true&wallet_data=true", gamekey)
	resp, err := c.c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

// TODO: Type definitions are not complete. PRs welcome.

type Order struct {
	AmountSpent      float64       `json:"amount_spent"`
	Product          Product       `json:"product"`
	Gamekey          string        `json:"gamekey"`
	UID              string        `json:"uid"`
	AllCouponData    []Coupon      `json:"all_coupon_data"`
	Created          string        `json:"created"`
	MissedCredit     interface{}   `json:"missed_credit"`
	Subproducts      []Subproduct  `json:"subproducts"`
	TotalChoices     int           `json:"total_choices"`
	TpkdDict         TpkdDict      `json:"tpkd_dict"`
	ChoicesRemaining int           `json:"choices_remaining"`
	Currency         string        `json:"currency"`
	IsGiftee         bool          `json:"is_giftee"`
	HasWallet        bool          `json:"has_wallet"`
	Claimed          bool          `json:"claimed"`
	Total            float64       `json:"total"`
	WalletCredit     *WalletCredit `json:"wallet_credit"`
	PathIDs          []string      `json:"path_ids"`
}

type Coupon struct {
	CouponMinCount       interface{} `json:"coupon_min_count"`
	CouponValidProducts  []string    `json:"coupon_valid_products"`
	CouponType           string      `json:"coupon_type"`
	CouponDiscount       interface{} `json:"coupon_discount"`
	CouponMachineName    string      `json:"coupon_machine_name"`
	CouponCredit         interface{} `json:"coupon_credit"`
	CouponMaxCount       interface{} `json:"coupon_max_count"`
	Subscriptions        []string    `json:"subscriptions"`
	CouponExcludeMonthly bool        `json:"coupon_exclude_monthly"`
	CouponExpiration     string      `json:"coupon_expiration"`
	CouponPrice          interface{} `json:"coupon_price"`
	CouponStack          bool        `json:"coupon_stack"`
	CouponKey            int64       `json:"coupon_key"`
	CouponStorefrontLink string      `json:"coupon_storefront_link"`
	StorefrontProduct    interface{} `json:"storefront_product"`
	Strings              interface{} `json:"strings"`
	CouponStatus         string      `json:"coupon_status"`
	CouponHumanName      string      `json:"coupon_human_name"`
}

type Product struct {
	Category           Category    `json:"category"`
	MachineName        string      `json:"machine_name"`
	EmptyTpkds         interface{} `json:"empty_tpkds"`
	PostPurchaseText   string      `json:"post_purchase_text"`
	HumanName          string      `json:"human_name"`
	PartialGiftEnabled bool        `json:"partial_gift_enabled"`
}

type Subproduct struct {
	MachineName               string             `json:"machine_name"`
	URL                       string             `json:"url"`
	Downloads                 []PlatformDownload `json:"downloads"`
	LibraryFamilyName         string             `json:"library_family_name"`
	Payee                     Payee              `json:"payee"`
	HumanName                 string             `json:"human_name"`
	CustomDownloadPageBoxCSS  string             `json:"custom_download_page_box_css"`
	CustomDownloadPageBoxHTML string             `json:"custom_download_page_box_html"`
	Icon                      string             `json:"icon"`
}

type PlatformDownload struct {
	MachineName           string      `json:"machine_name"`
	Platform              Platform    `json:"platform"`
	Downloads             []Download  `json:"download_struct"`
	OptionsDict           interface{} `json:"options_dict"`
	DownloadIdentifier    string      `json:"download_identifier"`
	AndroidAppOnly        bool        `json:"android_app_only"`
	DownloadVersionNumber interface{} `json:"download_version_number"`
}

type Download struct {
	Name             string        `json:"name"`
	URL              URL           `json:"url"`
	HumanSize        string        `json:"human_size"`
	FileSize         int64         `json:"file_size"`
	Small            int           `json:"small"`
	MD5              string        `json:"md5,omitempty"`
	SHA1             string        `json:"sha1,omitempty"`
	UploadedAt       string        `json:"uploaded_at,omitempty"` // TODO: custom format timestamp
	Timestamp        unixtime.Time `json:"timestamp,omitempty"`
	UsesKindleSender bool          `json:"uses_kindle_sender,omitempty"`
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
	AllTpks []AllTpk `json:"all_tpks"`
}

type AllTpk struct {
	IsGift                                bool             `json:"is_gift"`
	ExclusiveCountries                    []interface{}    `json:"exclusive_countries"`
	MachineName                           string           `json:"machine_name"`
	Gamekey                               string           `json:"gamekey"`
	CustomInstructionsHTML                string           `json:"custom_instructions_html,omitempty"`
	DisallowedCountries                   []string         `json:"disallowed_countries"`
	ShowCustomInstructionsInUserLibraries bool             `json:"show_custom_instructions_in_user_libraries"`
	KeyType                               KeyType          `json:"key_type"`
	KeyTypeHumanName                      KeyTypeHumanName `json:"key_type_human_name"`
	Visible                               bool             `json:"visible"`
	DisplaySeparately                     bool             `json:"display_separately"`
	RedeemedKeyVal                        string           `json:"redeemed_key_val"`
	Keyindex                              int64            `json:"keyindex"`
	HumanName                             string           `json:"human_name"`
	AutoExpand                            bool             `json:"auto_expand"`
	IsExpired                             bool             `json:"is_expired"`
	Class                                 Class            `json:"class"`
	NumDaysUntilExpired                   int64            `json:"num_days_until_expired"`
	InstructionsHTML                      string           `json:"instructions_html,omitempty"`
	SteamAppID                            int64            `json:"steam_app_id,omitempty"`
	PreinstructionText                    string           `json:"preinstruction_text,omitempty"`
	Disclaimer                            string           `json:"disclaimer,omitempty"`
}

type WalletCredit struct {
	Gamekey         string      `json:"gamekey"`
	ExpirableCredit bool        `json:"expirable_credit"`
	Expiry          interface{} `json:"expiry"`
	Currency        string      `json:"currency"`
	Amount          float64     `json:"amount"`
	Settled         bool        `json:"settled"`
}

type Category string

const (
	Bundle     Category = "bundle"
	Storefront Category = "storefront"
)

type Platform string

const (
	Audio   Platform = "audio"
	Ebook   Platform = "ebook"
	Windows Platform = "windows"
)

type Class string

const (
	Genericbutton Class = "genericbutton"
	Steambutton   Class = "steambutton"
)

type KeyType string

const (
	Generic KeyType = "generic"
	Steam   KeyType = "steam"
)

type KeyTypeHumanName string

const (
	KeyTypeHumanNameSteam KeyTypeHumanName = "Steam"
	OtherKey              KeyTypeHumanName = "other-key"
	Paizo                 KeyTypeHumanName = "Paizo"
	ProFantasy            KeyTypeHumanName = "ProFantasy"
)
