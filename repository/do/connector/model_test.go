package connector

import "testing"

func Test_NextPage(t *testing.T) {
	cache := &TableMetaCache{
		ID: "test",
		RecordMap: map[string]*Record{
			"1": {
				Guid:    "1",
				PubDate: 1,
			},
			"2": {
				Guid:    "2",
				PubDate: 2,
			},
			"3": {
				Guid:    "3",
				PubDate: 3,
			},
			"4": {
				Guid:    "4",
				PubDate: 4,
			},
		},
	}

	cache.RecordPage = cache.SortByTimeASC()

	perPage, nextGuid := cache.RecordPage.NextPage("", 2)
	if len(perPage) != 2 {
		t.Errorf("len(perPage) != 2")
	}
	if perPage[0].Guid != "1" {
		t.Errorf("perPage[0].Guid != \"1\"")
	}
	if perPage[1].Guid != "2" {
		t.Errorf("perPage[1].Guid != \"2\"")
	}
	if nextGuid != "3" {
		t.Errorf("nextGuid != \"3\"")
	}

	perPage2, nextGuid2 := cache.RecordPage.NextPage(nextGuid, 2)
	if len(perPage2) != 2 {
		t.Errorf("len(perPage2) != 2")
	}
	if perPage2[0].Guid != "3" {
		t.Errorf("perPage2[0].Guid != \"3\"")
	}
	if perPage2[1].Guid != "4" {
		t.Errorf("perPage2[1].Guid != \"4\"")
	}
	if nextGuid2 != "" {
		t.Errorf("nextGuid2 != \"\"")
	}
}

func Test_LimitAndSave(t *testing.T) {
	cache := &TableMetaCache{
		ID: "test",
		RecordMap: map[string]*Record{
			"1": {
				Guid:    "1",
				PubDate: 1,
			},
			"2": {
				Guid:    "2",
				PubDate: 2,
			},
			"3": {
				Guid:    "3",
				PubDate: 3,
			},
			"4": {
				Guid:    "4",
				PubDate: 4,
			},
		},
	}

	cache.RecordPage = cache.SortByTimeASC()

	cache.LimitAndSave(2)
	if len(cache.RecordPage) != 4 {
		t.Errorf("len(cache.RecordPage) != 4")
	}
	if len(cache.RecordMap) != 2 {
		t.Errorf("len(cache.RecordMap) != 2")
	}
	if cache.RecordMap["3"].Guid != "3" {
		t.Errorf("cache.RecordMap[\"3\"].Guid != \"3\"")
	}
	if cache.RecordMap["4"].Guid != "4" {
		t.Errorf("cache.RecordMap[\"4\"].Guid != \"4\"")
	}
	if _, ok := cache.RecordMap["1"]; ok {
		t.Errorf("cache.RecordMap[\"1\"] should not exist")
	}
	if _, ok := cache.RecordMap["2"]; ok {
		t.Errorf("cache.RecordMap[\"2\"] should not exist")
	}
}
