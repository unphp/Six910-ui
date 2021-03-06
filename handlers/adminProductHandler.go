package handlers

import (
	"net/http"
	"strconv"

	sdbi "github.com/Ulbora/six910-database-interface"
	"github.com/gorilla/mux"
)

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

//ProdError ProdError
type ProdError struct {
	Error    string
	Product  *sdbi.Product
	Products *[]sdbi.Product
}

//StoreAdminAddProductPage StoreAdminAddProductPage
func (h *Six910Handler) StoreAdminAddProductPage(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prod add view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			loginErr := r.URL.Query().Get("error")
			var lge ProcError
			lge.Error = loginErr
			h.AdminTemplates.ExecuteTemplate(w, adminAddProductPage, &lge)
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminAddProduct StoreAdminAddProduct
func (h *Six910Handler) StoreAdminAddProduct(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prod add", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			p := h.processProduct(r)
			h.Log.Debug("prod add", *p)
			hd := h.getHeader(s)
			prres := h.API.AddProduct(p, hd)
			h.Log.Debug("prod add resp", *prres)
			if prres.Success {
				http.Redirect(w, r, adminAddProdView, http.StatusFound)
			} else {
				http.Redirect(w, r, adminAddProdViewFail, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminEditProductPage StoreAdminEditProductPage
func (h *Six910Handler) StoreAdminEditProductPage(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prod add view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			hd := h.getHeader(s)
			epvars := mux.Vars(r)
			idstr := epvars["id"]
			prodID, _ := strconv.ParseInt(idstr, 10, 64)
			h.Log.Debug("prod id in edit", prodID)
			prod := h.API.GetProductByID(prodID, hd)
			h.Log.Debug("prod  in edit", prod)
			edErr := r.URL.Query().Get("error")
			var epparm ProdError
			epparm.Error = edErr
			epparm.Product = prod
			h.AdminTemplates.ExecuteTemplate(w, adminEditProductPage, &epparm)
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminEditProduct StoreAdminEditProduct
func (h *Six910Handler) StoreAdminEditProduct(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prod edit", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			epp := h.processProduct(r)
			h.Log.Debug("prod update", *epp)
			hd := h.getHeader(s)
			res := h.API.UpdateProduct(epp, hd)
			h.Log.Debug("prod update resp", *res)
			if res.Success {
				http.Redirect(w, r, adminProductListView, http.StatusFound)
			} else {
				http.Redirect(w, r, adminEditProdViewFail, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminViewProductList StoreAdminViewProductList
func (h *Six910Handler) StoreAdminViewProductList(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prods view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			hd := h.getHeader(s)
			vpvars := mux.Vars(r)
			ststr := vpvars["start"]
			endstr := vpvars["end"]
			vpstart, _ := strconv.ParseInt(ststr, 10, 64)
			vpend, _ := strconv.ParseInt(endstr, 10, 64)
			prods := h.API.GetProductList(vpstart, vpend, hd)
			plErr := r.URL.Query().Get("error")
			var plparm ProdError
			plparm.Error = plErr
			plparm.Products = prods
			h.Log.Debug("prods  in edit", prods)
			h.AdminTemplates.ExecuteTemplate(w, adminProductListPage, &plparm)
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminDeleteProduct StoreAdminDeleteProduct
func (h *Six910Handler) StoreAdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	s, suc := h.getSession(r)
	h.Log.Debug("session suc in prod list delete", suc)
	if suc {
		if h.isStoreAdminLoggedIn(s) {
			hd := h.getHeader(s)
			dpvars := mux.Vars(r)
			idstrd := dpvars["id"]
			idd, _ := strconv.ParseInt(idstrd, 10, 64)
			res := h.API.DeleteProduct(idd, hd)
			h.Log.Debug("prod delete resp", *res)
			if res.Success {
				http.Redirect(w, r, adminProductListView, http.StatusFound)
			} else {
				http.Redirect(w, r, adminProductListViewFail, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

func (h *Six910Handler) processProduct(r *http.Request) *sdbi.Product {
	var p sdbi.Product
	id := r.FormValue("id")
	p.ID, _ = strconv.ParseInt(id, 10, 64)
	sku := r.FormValue("sku")
	p.Sku = sku
	p.Gtin = r.FormValue("gtin")
	p.Name = r.FormValue("name")
	p.ShortDesc = r.FormValue("shortDest")
	p.Desc = r.FormValue("desc")
	cost := r.FormValue("cost")
	p.Cost, _ = strconv.ParseFloat(cost, 64)
	msrp := r.FormValue("msrp")
	p.Msrp, _ = strconv.ParseFloat(msrp, 64)
	mapp := r.FormValue("map")
	p.Map, _ = strconv.ParseFloat(mapp, 64)
	price := r.FormValue("price")
	p.Price, _ = strconv.ParseFloat(price, 64)
	salePrice := r.FormValue("salePrice")
	p.SalePrice, _ = strconv.ParseFloat(salePrice, 64)
	p.Currency = r.FormValue("currency")
	p.ManufacturerID = r.FormValue("manfId")
	p.Manufacturer = r.FormValue("manf")
	stock := r.FormValue("stock")
	p.Stock, _ = strconv.ParseInt(stock, 10, 64)
	stockAlrt := r.FormValue("stockAlrt")
	p.StockAlert, _ = strconv.ParseInt(stockAlrt, 10, 64)
	weight := r.FormValue("weight")
	p.Weight, _ = strconv.ParseFloat(weight, 64)
	width := r.FormValue("width")
	p.Width, _ = strconv.ParseFloat(width, 64)
	height := r.FormValue("height")
	p.Height, _ = strconv.ParseFloat(height, 64)
	depth := r.FormValue("depth")
	p.Depth, _ = strconv.ParseFloat(depth, 64)
	shipMarkup := r.FormValue("shipMarkup")
	p.ShippingMarkup, _ = strconv.ParseFloat(shipMarkup, 64)
	visible := r.FormValue("visible")
	p.Visible, _ = strconv.ParseBool(visible)
	searchable := r.FormValue("searchable")
	p.Searchable, _ = strconv.ParseBool(searchable)
	multibox := r.FormValue("multibox")
	p.MultiBox, _ = strconv.ParseBool(multibox)
	shipSep := r.FormValue("shipSep")
	p.ShipSeparately, _ = strconv.ParseBool(shipSep)
	freeShipping := r.FormValue("freeShipping")
	p.FreeShipping, _ = strconv.ParseBool(freeShipping)
	promoted := r.FormValue("promoted")
	p.Promoted, _ = strconv.ParseBool(promoted)
	dropship := r.FormValue("dropship")
	p.Dropship, _ = strconv.ParseBool(dropship)
	specialProc := r.FormValue("specialProc")
	p.SpecialProcessing, _ = strconv.ParseBool(specialProc)
	specialProcType := r.FormValue("specialProcType")
	p.SpecialProcessingType = specialProcType
	p.Size = r.FormValue("size")
	p.Color = r.FormValue("color")
	p.Thumbnail = r.FormValue("thumbnail")
	p.Image1 = r.FormValue("image1")
	p.Image2 = r.FormValue("image2")
	p.Image3 = r.FormValue("image3")
	p.Image4 = r.FormValue("image4")
	distributorID := r.FormValue("distributorId")
	p.DistributorID, _ = strconv.ParseInt(distributorID, 10, 64)
	storeID := r.FormValue("storeId")
	p.StoreID, _ = strconv.ParseInt(storeID, 10, 64)
	parentProductID := r.FormValue("parentProductId")
	p.ParentProductID, _ = strconv.ParseInt(parentProductID, 10, 64)

	return &p
}
