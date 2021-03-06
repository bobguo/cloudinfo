// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/banzaicloud/cloudinfo/pkg/cloudinfo-client/models"
)

// GetProductsReader is a Reader for the GetProducts structure.
type GetProductsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetProductsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetProductsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetProductsOK creates a GetProductsOK with default headers values
func NewGetProductsOK() *GetProductsOK {
	return &GetProductsOK{}
}

/*GetProductsOK handles this case with default header values.

ProductDetailsResponse
*/
type GetProductsOK struct {
	Payload *models.ProductDetailsResponse
}

func (o *GetProductsOK) Error() string {
	return fmt.Sprintf("[GET /providers/{provider}/services/{service}/regions/{region}/products][%d] getProductsOK  %+v", 200, o.Payload)
}

func (o *GetProductsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ProductDetailsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
