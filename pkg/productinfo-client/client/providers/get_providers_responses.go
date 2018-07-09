// Code generated by go-swagger; DO NOT EDIT.

package providers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/banzaicloud/productinfo/pkg/productinfo-client/models"
)

// GetProvidersReader is a Reader for the GetProviders structure.
type GetProvidersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetProvidersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetProvidersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetProvidersOK creates a GetProvidersOK with default headers values
func NewGetProvidersOK() *GetProvidersOK {
	return &GetProvidersOK{}
}

/*GetProvidersOK handles this case with default header values.

ProviderResponse
*/
type GetProvidersOK struct {
	Payload models.ProviderResponse
}

func (o *GetProvidersOK) Error() string {
	return fmt.Sprintf("[GET /providers][%d] getProvidersOK  %+v", 200, o.Payload)
}

func (o *GetProvidersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
