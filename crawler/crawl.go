package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/phuwn/go-crawler-example/db"
)

const squadDataEndPoint = `https://nft.pancakeswap.com/api/v1/collections/0x0a8901b0E25DEb55A87524f0cC164E9644020EBA/tokens/%d`

type SquadResp struct {
	Data *PancakeSquad `json:"data"`
}

func fetchPancakeSquad(client *http.Client, id int) (*PancakeSquad, error) {
	resp, err := client.Get(fmt.Sprintf(squadDataEndPoint, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed to get data with status: %d, and response %s", resp.StatusCode, b)
	}

	var respData SquadResp
	if err = json.Unmarshal(b, &respData); err != nil {
		return nil, err
	}

	return respData.Data, nil
}

func CrawlPancakeSquad(client *http.Client, id int) error {
	squad, err := fetchPancakeSquad(client, id)
	if err != nil {
		return err
	}

	tx := db.Get()
	if err := tx.Create(squad).Error; err != nil {
		return err
	}

	return nil
}
