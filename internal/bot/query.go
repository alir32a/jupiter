package bot

import "encoding/json"

const (
	QueryActionGetResource = "get_resource"
	QueryActionCancel      = "cancel"
)

type Query struct {
	Action   string `json:"action"`
	Resource string `json:"resource,omitempty"`
	Param    string `json:"param,omitempty"`
}

func (q *Query) SetResource(v string) *Query {
	q.Resource = v

	return q
}

func (q *Query) SetParam(v string) *Query {
	q.Param = v

	return q
}

func (q *Query) Marshal() (string, error) {
	data, err := json.Marshal(&q)
	if err != nil {
		return "", err
	}

	return string(data), err
}

func NewQuery(action string) *Query {
	return &Query{
		Action: action,
	}
}

func Unmarshal(v string) (*Query, error) {
	var q Query

	return &q, json.Unmarshal([]byte(v), &q)
}
