// Code generated by go-swagger; DO NOT EDIT.

// Copyright (c) 2020 Samuel Mutel <12967891+smutel@users.noreply.github.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewUsersPermissionsBulkDeleteParams creates a new UsersPermissionsBulkDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUsersPermissionsBulkDeleteParams() *UsersPermissionsBulkDeleteParams {
	return &UsersPermissionsBulkDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUsersPermissionsBulkDeleteParamsWithTimeout creates a new UsersPermissionsBulkDeleteParams object
// with the ability to set a timeout on a request.
func NewUsersPermissionsBulkDeleteParamsWithTimeout(timeout time.Duration) *UsersPermissionsBulkDeleteParams {
	return &UsersPermissionsBulkDeleteParams{
		timeout: timeout,
	}
}

// NewUsersPermissionsBulkDeleteParamsWithContext creates a new UsersPermissionsBulkDeleteParams object
// with the ability to set a context for a request.
func NewUsersPermissionsBulkDeleteParamsWithContext(ctx context.Context) *UsersPermissionsBulkDeleteParams {
	return &UsersPermissionsBulkDeleteParams{
		Context: ctx,
	}
}

// NewUsersPermissionsBulkDeleteParamsWithHTTPClient creates a new UsersPermissionsBulkDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewUsersPermissionsBulkDeleteParamsWithHTTPClient(client *http.Client) *UsersPermissionsBulkDeleteParams {
	return &UsersPermissionsBulkDeleteParams{
		HTTPClient: client,
	}
}

/*
UsersPermissionsBulkDeleteParams contains all the parameters to send to the API endpoint

	for the users permissions bulk delete operation.

	Typically these are written to a http.Request.
*/
type UsersPermissionsBulkDeleteParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the users permissions bulk delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UsersPermissionsBulkDeleteParams) WithDefaults() *UsersPermissionsBulkDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the users permissions bulk delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UsersPermissionsBulkDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) WithTimeout(timeout time.Duration) *UsersPermissionsBulkDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) WithContext(ctx context.Context) *UsersPermissionsBulkDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) WithHTTPClient(client *http.Client) *UsersPermissionsBulkDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the users permissions bulk delete params
func (o *UsersPermissionsBulkDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *UsersPermissionsBulkDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}