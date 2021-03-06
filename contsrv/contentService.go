package contsrv

/*
 Six910 is a shopping cart and E-commerce system.
 Copyright (C) 2020 Ulbora Labs LLC. (www.ulboralabs.com)
 All rights reserved.
 Copyright (C) 2020 Ken Williamson
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	b64 "encoding/base64"
	"encoding/json"
	"html/template"
	"time"
)

const (
	contentExistsCode   = 1001
	contentNotFoundCode = 1002
)

//Content content
type Content struct {
	Name              string    `json:"name"`
	Title             string    `json:"title"`
	Subject           string    `json:"subject"`
	Author            string    `json:"author"`
	CreateDate        time.Time `json:"createDate"`
	ModifiedDate      time.Time `json:"modifiedDate"`
	Hits              int64     `json:"hits"`
	MetaAuthorName    string    `json:"metaAuthorName"`
	MetaDesc          string    `json:"metaDesc"`
	MetaKeyWords      string    `json:"metaKeyWords"`
	MetaRobotKeyWords string    `json:"metaRobotKeyWords"`
	Text              string    `json:"text"`
	TextHTML          template.HTML
	Archived          bool `json:"archived"`
	Visible           bool `json:"visible"`
	UseModifiedDate   bool
	BlogPost          bool `json:"blogPost"`
}

// PageHead used for page head
type PageHead struct {
	Title        string
	MetaAuthor   string
	MetaDesc     string
	MetaKeyWords string
}

//Response res
type Response struct {
	Success  bool   `json:"success"`
	Name     string `json:"name"`
	FailCode int    `json:"failCode"`
}

//AddContent add content
func (c *CmsService) AddContent(content *Content) *Response {
	var rtn = new(Response)
	content.CreateDate = time.Now()
	content.Text = b64.StdEncoding.EncodeToString([]byte(content.Text))
	c.Log.Debug("content in add: ", *content)
	ec := c.Store.Read(content.Name)
	//c.Log.Debug("found content in add: ", *ec)
	if *ec == nil {
		suc := c.Store.Save(content.Name, content)
		rtn.Success = suc
		rtn.Name = content.Name
	} else {
		rtn.FailCode = contentExistsCode
	}
	return rtn
}

//UpdateContent add content
func (c *CmsService) UpdateContent(content *Content) *Response {
	var rtn = new(Response)
	content.Text = b64.StdEncoding.EncodeToString([]byte(content.Text))
	c.Log.Debug("content in update: ", *content)
	ec := c.Store.Read(content.Name)
	if *ec != nil {
		var cd Content
		err := json.Unmarshal(*ec, &cd)
		c.Log.Debug("found content in update: ", cd)
		if err == nil {
			cd.Archived = content.Archived
			cd.Hits = content.Hits
			cd.MetaAuthorName = content.MetaAuthorName
			cd.MetaDesc = content.MetaDesc
			cd.MetaKeyWords = content.MetaKeyWords
			cd.MetaRobotKeyWords = content.MetaRobotKeyWords
			cd.ModifiedDate = time.Now()
			cd.Text = content.Text
			cd.Title = content.Title
			cd.Subject = content.Subject
			cd.Author = content.Author
			cd.Visible = content.Visible
			cd.BlogPost = content.BlogPost
			suc := c.Store.Save(content.Name, cd)
			rtn.Success = suc
			rtn.Name = content.Name
		}
	} else {
		rtn.FailCode = contentNotFoundCode
	}
	return rtn
}

// GetContent get content
func (c *CmsService) GetContent(name string) (bool, *Content) {
	var rtn Content
	var suc bool
	ec := c.Store.Read(name)
	if *ec != nil {
		var cd Content
		err := json.Unmarshal(*ec, &cd)
		if err == nil {
			txt, err2 := b64.StdEncoding.DecodeString(cd.Text)
			if err2 == nil {
				cd.Text = string(txt)
				cd.TextHTML = template.HTML(cd.Text)
				if cd.ModifiedDate.Year() != 1 {
					cd.UseModifiedDate = true
				}
				rtn = cd
				c.HitTotal++
				c.ContentHits[name]++
				// if c.HitTotal >= c.HitLimit {
				// 	c.SaveHits()
				// }
				suc = true
			}
		}
	}
	return suc, &rtn
}

// GetContentList get content list by client
func (c *CmsService) GetContentList(published bool) *[]Content {
	var rtn []Content
	res := c.Store.ReadAll()
	//c.Log.Debug("found content bytes in list: ", *res)
	for r := range *res {
		var ct Content
		err := json.Unmarshal((*res)[r], &ct)
		c.Log.Debug("found content item in list: ", ct)
		if err == nil {
			if published && !ct.Visible {
				continue
			} else {
				txt, err2 := b64.StdEncoding.DecodeString((ct.Text))
				c.Log.Debug("found content item in list err2: ", err2)
				if err2 == nil {
					ct.Text = string(txt)
					ct.TextHTML = template.HTML(ct.Text)
					if ct.ModifiedDate.Year() != 1 {
						ct.UseModifiedDate = true
					}
					c.Log.Debug("found content item in list before append: ", ct)
					rtn = append(rtn, ct)
				}
			}
		}
	}
	return &rtn
}

// DeleteContent delete content
func (c *CmsService) DeleteContent(name string) *Response {
	var rtn = new(Response)
	suc := c.Store.Delete(name)
	if suc {
		rtn.Success = true
		rtn.Name = name
	}
	return rtn
}

//SaveHits SaveHits
func (c *CmsService) SaveHits() {
	c.hitmu.Lock()
	defer c.hitmu.Unlock()
	for n, h := range c.ContentHits {
		c.Log.Debug("found content name in content hits loop: ", n)
		suc, ct := c.GetContent(n)
		c.Log.Debug("found content suc in content hits loop: ", suc)
		c.Log.Debug("found content suc in content hits loop: ", *ct)
		if suc {
			ct.Hits += h
			res := c.UpdateContent(ct)
			c.Log.Debug("update content in content hits loop: ", *res)
			c.ContentHits[n] = 0
		}
	}
	c.HitTotal = 0
}

//HitCheck HitCheck
func (c *CmsService) HitCheck() {
	c.Log.Debug("in hitCheck c.HitLimit: ", c.HitLimit)
	c.Log.Debug("in hitCheck c.HitTotal: ", c.HitTotal)
	if c.HitTotal >= c.HitLimit {
		c.SaveHits()
	}
}
