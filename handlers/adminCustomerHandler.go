package handlers

import (
	"net/http"
	"strconv"

	api "github.com/Ulbora/Six910API-Go"
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

//CusPage CusPage
type CusPage struct {
	Error    string
	Customer *sdbi.Customer
	User     *api.UserResponse
}

//StoreAdminEditCustomerPage StoreAdminEditCustomerPage
func (h *Six910Handler) StoreAdminEditCustomerPage(w http.ResponseWriter, r *http.Request) {
	ecps, suc := h.getSession(r)
	h.Log.Debug("session suc in customer edit view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(ecps) {
			hd := h.getHeader(ecps)
			epvars := mux.Vars(r)
			cidstr := epvars["id"]
			cID, _ := strconv.ParseInt(cidstr, 10, 64)
			h.Log.Debug("customer id in edit", cID)
			cust := h.API.GetCustomerID(cID, hd)
			h.Log.Debug("customer  in edit", cust)
			edErr := r.URL.Query().Get("error")
			var ceparm CusPage
			ceparm.Error = edErr
			ceparm.Customer = cust
			h.AdminTemplates.ExecuteTemplate(w, adminEditCustomerPage, &ceparm)
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminEditCustomer StoreAdminEditCustomer
func (h *Six910Handler) StoreAdminEditCustomer(w http.ResponseWriter, r *http.Request) {
	ecs, suc := h.getSession(r)
	h.Log.Debug("session suc in customer edit", suc)
	if suc {
		if h.isStoreAdminLoggedIn(ecs) {
			c := h.processCustomer(r)
			h.Log.Debug("customer edit", *c)
			hd := h.getHeader(ecs)
			ecres := h.API.UpdateCustomer(c, hd)
			h.Log.Debug("customer edit resp", *ecres)
			if ecres.Success {
				http.Redirect(w, r, adminCustomerListView, http.StatusFound)
			} else {
				http.Redirect(w, r, adminEditCustomerViewFail, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminEditCustomerUserPage StoreAdminEditCustomerUserPage
func (h *Six910Handler) StoreAdminEditCustomerUserPage(w http.ResponseWriter, r *http.Request) {
	ecups, suc := h.getSession(r)
	h.Log.Debug("session suc in customer edit view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(ecups) {
			hd := h.getHeader(ecups)
			ecuvars := mux.Vars(r)
			unstr := ecuvars["username"]
			eucid := ecuvars["cid"]
			cID, _ := strconv.ParseInt(eucid, 10, 64)
			h.Log.Debug("customer username in edit", unstr)
			var u api.User
			u.Username = unstr
			cusr := h.API.GetUser(&u, hd)
			if cusr.CustomerID == cID {
				h.Log.Debug("customer user in edit", cusr)
				edErr := r.URL.Query().Get("error")
				var ceparm CusPage
				ceparm.Error = edErr
				ceparm.User = cusr
				h.AdminTemplates.ExecuteTemplate(w, adminEditCustomerUserPage, &ceparm)
			} else {
				http.Redirect(w, r, adminCustomerListView, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminEditCustomerUser StoreAdminEditCustomerUser
func (h *Six910Handler) StoreAdminEditCustomerUser(w http.ResponseWriter, r *http.Request) {
	ecus, suc := h.getSession(r)
	h.Log.Debug("session suc in customer edit", suc)
	if suc {
		if h.isStoreAdminLoggedIn(ecus) {
			cu := h.processCustomerUser(r)
			h.Log.Debug("customer user edit", *cu)
			hd := h.getHeader(ecus)
			ecures := h.API.UpdateUser(cu, hd)
			h.Log.Debug("customer user edit resp", *ecures)
			if ecures.Success {
				http.Redirect(w, r, adminCustomerListView, http.StatusFound)
			} else {
				http.Redirect(w, r, adminEditCustomerUserViewFail, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

//StoreAdminViewCustomerList StoreAdminViewCustomerList
func (h *Six910Handler) StoreAdminViewCustomerList(w http.ResponseWriter, r *http.Request) {
	culs, suc := h.getSession(r)
	h.Log.Debug("session suc in customer user list view", suc)
	if suc {
		if h.isStoreAdminLoggedIn(culs) {
			hd := h.getHeader(culs)
			cul := h.API.GetCustomerList(hd)
			h.Log.Debug("customer  in list", cul)
			h.AdminTemplates.ExecuteTemplate(w, adminCustomerListPage, &cul)
		} else {
			http.Redirect(w, r, adminloginPage, http.StatusFound)
		}
	}
}

func (h *Six910Handler) processCustomer(r *http.Request) *sdbi.Customer {
	var c sdbi.Customer
	id := r.FormValue("id")
	c.ID, _ = strconv.ParseInt(id, 10, 64)
	c.Email = r.FormValue("email")
	resetPw := r.FormValue("resetPassword")
	c.ResetPassword, _ = strconv.ParseBool(resetPw)
	c.FirstName = r.FormValue("firstName")
	c.LastName = r.FormValue("lastName")
	c.Company = r.FormValue("company")
	c.City = r.FormValue("city")
	c.State = r.FormValue("state")
	c.Zip = r.FormValue("zip")
	c.Phone = r.FormValue("phone")
	storeID := r.FormValue("storeId")
	c.StoreID, _ = strconv.ParseInt(storeID, 10, 64)
	return &c
}

func (h *Six910Handler) processCustomerUser(r *http.Request) *api.User {
	var c api.User
	c.Username = r.FormValue("username")
	c.Password = r.FormValue("password")
	c.OldPassword = r.FormValue("oldPassword")
	c.Role = r.FormValue("role")
	cID := r.FormValue("cid")
	c.CustomerID, _ = strconv.ParseInt(cID, 10, 64)
	storeID := r.FormValue("storeId")
	c.StoreID, _ = strconv.ParseInt(storeID, 10, 64)
	enabled := r.FormValue("enabled")
	c.Enabled, _ = strconv.ParseBool(enabled)
	return &c
}
