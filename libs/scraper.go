package libs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

type FetchType2 int

const (
	FetchAll2 FetchType2 = iota
	FetchVendor2
	FetchCatalogs2
	FetchItems2
)

const key2 = "6953-oaypsHZN88GndHcNtyVBktnyy62VPAjK4qYAT9ga3XEhm6QytSMSjHqLvweL73yhGBerV9x8mEVvBt3A"

type Result2 struct {
	Response struct {
		Page struct {
			Current int `json:"current"`
			Next    int `json:"next"`
			Prev    int `json:"prev"`
			Pages   int `json:"pages"`
			Items   int `json:"items"`
		} `json:"page"`
		Vendors []struct {
			Name  string `json:"name"`
			Alias string `json:"alias"`
		} `json:"vendors"`
		Catalogs []struct {
			ID       string `json:"va_catalog_id"`
			ParentID string `json:"va_parent_id"`
			Name     string `json:"name"`
		} `json:"catalogs"`
		Items []struct {
			Images     []string `json:"images"`
			Code       string   `json:"p_code"`
			Mog        string   `json:"mog"`
			OEMNum     string   `json:"oem_num"`
			OEMBrand   string   `json:"oem_brand"`
			Name       string   `json:"name"`
			Shipment   int      `json:"shipment"`
			Delivery   int      `json:"delivery"`
			Department string   `json:"department"`
			Count      int      `json:"count"`
			CountChel  int      `json:"count_chel"`
			CountEkb   int      `json:"count_ekb"`
			UnitCode   int      `json:"unit_code"`
			Unit       string   `json:"unit"`
			Price      float32  `json:"price"`
			CatalogID  string   `json:"va_catalog_id"`
			ItemID     string   `json:"va_item_id"`
		} `json:"items"`
	} `json:"response"`
}

func Scrape(tp FetchType2, page int) Result {
	c := colly.NewCollector(colly.AllowURLRevisit())

	// Rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	// Print the response
	c.OnResponse(func(resp *colly.Response) {

		r := Result{}
		err = json.Unmarshal(resp.Body, &r)
		if err != nil {
			log.Println("Error")
		}

		switch tp {
		case FetchVendor2:
			Data.Response.Vendors = append(Data.Response.Vendors, r.Response.Vendors...)
		case FetchCatalogs2:
			Data.Response.Catalogs = append(Data.Response.Catalogs, r.Response.Catalogs...)
		case FetchItems2:
			Data.Response.Items = append(Data.Response.Items, r.Response.Items...)
		}

		log.Printf("V: %d, C: %d, I: %d", len(Data.Response.Vendors), len(Data.Response.Catalogs), len(Data.Response.Items))
	})

	if page == 0 {
		if tp == FetchItems2 || tp == FetchVendor2 {
			var wg sync.WaitGroup
			rs := get(tp, 0)
			log.Printf("Total pages: %d", rs.Response.Page.Pages)

			for i := 1; i < rs.Response.Page.Pages; i++ {
				wg.Add(1)

				go func(indx int, wgr *sync.WaitGroup) {
					//get(tp, indx)

					url := "https://api.v-avto.ru/v1/"
					switch tp {
					case FetchVendor2:
						url += "vendors"
					case FetchCatalogs2:
						url += "catalogs"
					case FetchItems2:
						url += "items"
					}
					reqUrl := fmt.Sprintf("%s?key=%s&page=%d", url, key, page)
					_ = c.Visit(reqUrl)

					defer func() {
						wgr.Done()
					}()
				}(i, &wg)
			}

			wg.Wait()
		}
	}

	return Data
}

func get(tp FetchType2, page int) Result {
	url := "https://api.v-avto.ru/v1/"
	switch tp {
	case FetchVendor2:
		url += "vendors"
	case FetchCatalogs2:
		url += "catalogs"
	case FetchItems2:
		url += "items"
	}

	reqUrl := fmt.Sprintf("%s?key=%s&page=%d", url, key, page)

	//log.Printf("Req: %s",reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%d", rand.Int()))

	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	r := Result{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		time.Sleep(500 * time.Microsecond)
		return get(tp, page)
	}

	switch tp {
	case FetchVendor2:
		Data.Response.Vendors = append(Data.Response.Vendors, r.Response.Vendors...)
	case FetchCatalogs2:
		Data.Response.Catalogs = append(Data.Response.Catalogs, r.Response.Catalogs...)
	case FetchItems2:
		Data.Response.Items = append(Data.Response.Items, r.Response.Items...)
	}

	log.Printf("V: %d, C: %d, I: %d", len(Data.Response.Vendors), len(Data.Response.Catalogs), len(Data.Response.Items))

	return r
}
