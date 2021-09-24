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

package extras

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// ExtrasCustomLinksBulkDeleteReader is a Reader for the ExtrasCustomLinksBulkDelete structure.
type ExtrasCustomLinksBulkDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ExtrasCustomLinksBulkDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewExtrasCustomLinksBulkDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewExtrasCustomLinksBulkDeleteNoContent creates a ExtrasCustomLinksBulkDeleteNoContent with default headers values
func NewExtrasCustomLinksBulkDeleteNoContent() *ExtrasCustomLinksBulkDeleteNoContent {
	return &ExtrasCustomLinksBulkDeleteNoContent{}
}

/* ExtrasCustomLinksBulkDeleteNoContent describes a response with status code 204, with default header values.

ExtrasCustomLinksBulkDeleteNoContent extras custom links bulk delete no content
*/
type ExtrasCustomLinksBulkDeleteNoContent struct {
}

func (o *ExtrasCustomLinksBulkDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /extras/custom-links/][%d] extrasCustomLinksBulkDeleteNoContent ", 204)
}

func (o *ExtrasCustomLinksBulkDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}