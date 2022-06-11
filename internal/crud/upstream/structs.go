package upstream

type Upstream struct {
	ID   string `json:"id" bson:"id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	URL  string `json:"url" bson:"url,omitempty"`
}
