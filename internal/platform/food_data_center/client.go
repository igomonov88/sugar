package food_data_center

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.opencensus.io/trace"
)

// Config is the required properties to use food data center search api
type Config struct {
	ConsumerKey string
	APIURL      string
}

// StatusCheck returns nil if it can successfully talk to the food data center api. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, cfg Config) error {
	ctx, span := trace.StartSpan(ctx, "platform.Search.StatusCheck")
	defer span.End()
	fs := FoodSearchRequest{
		GeneralSearchInput: "cheese",
	}
	b, err := json.Marshal(&fs)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	resp, err := http.Post(cfg.APIURL, "application/json", buf)
	if resp.StatusCode != http.StatusOK {
		return errors.New("food data center api respond with status not OK")
	}
	return nil
}
