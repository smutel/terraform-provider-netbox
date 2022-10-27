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

package ipam

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

// NewIpamRouteTargetsBulkDeleteParams creates a new IpamRouteTargetsBulkDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewIpamRouteTargetsBulkDeleteParams() *IpamRouteTargetsBulkDeleteParams {
	return &IpamRouteTargetsBulkDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewIpamRouteTargetsBulkDeleteParamsWithTimeout creates a new IpamRouteTargetsBulkDeleteParams object
// with the ability to set a timeout on a request.
func NewIpamRouteTargetsBulkDeleteParamsWithTimeout(timeout time.Duration) *IpamRouteTargetsBulkDeleteParams {
	return &IpamRouteTargetsBulkDeleteParams{
		timeout: timeout,
	}
}

// NewIpamRouteTargetsBulkDeleteParamsWithContext creates a new IpamRouteTargetsBulkDeleteParams object
// with the ability to set a context for a request.
func NewIpamRouteTargetsBulkDeleteParamsWithContext(ctx context.Context) *IpamRouteTargetsBulkDeleteParams {
	return &IpamRouteTargetsBulkDeleteParams{
		Context: ctx,
	}
}

// NewIpamRouteTargetsBulkDeleteParamsWithHTTPClient creates a new IpamRouteTargetsBulkDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewIpamRouteTargetsBulkDeleteParamsWithHTTPClient(client *http.Client) *IpamRouteTargetsBulkDeleteParams {
	return &IpamRouteTargetsBulkDeleteParams{
		HTTPClient: client,
	}
}

/*
IpamRouteTargetsBulkDeleteParams contains all the parameters to send to the API endpoint

	for the ipam route targets bulk delete operation.

	Typically these are written to a http.Request.
*/
type IpamRouteTargetsBulkDeleteParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the ipam route targets bulk delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IpamRouteTargetsBulkDeleteParams) WithDefaults() *IpamRouteTargetsBulkDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the ipam route targets bulk delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *IpamRouteTargetsBulkDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) WithTimeout(timeout time.Duration) *IpamRouteTargetsBulkDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) WithContext(ctx context.Context) *IpamRouteTargetsBulkDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) WithHTTPClient(client *http.Client) *IpamRouteTargetsBulkDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ipam route targets bulk delete params
func (o *IpamRouteTargetsBulkDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *IpamRouteTargetsBulkDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}