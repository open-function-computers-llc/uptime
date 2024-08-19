package webhook

type Webhook struct {
	ID        int    `db:"id"`
	WebsiteID int    `db:"website_id"`
	Name      string `db:"hook_name"`
	URL       string `db:"hook_url"`
	Verb      string `db:"hook_verb"`
	Type      int    `db:"hook_type"`
}

type Hooktype int

var standard Hooktype = 1
var emergency Hooktype = 2
