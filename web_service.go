package restful

import (
	"log"
)

// WebService holds a collection of Route values that bind a Http Method + URL Path to a function.
type WebService struct {
	rootPath       string
	pathExpression *PathExpression // cached compilation of rootPath as RegExp
	routes         []Route
	produces       []string
	consumes       []string
	pathParameters []*Parameter
	filters        []FilterFunction
}

// Path specifies the root URL template path of the WebService.
// All Routes will be relative to this path.
func (self *WebService) Path(root string) *WebService {
	self.rootPath = root
	compiled, err := NewPathExpression(root)
	if err != nil {
		log.Fatalf("[restful] Invalid path:%s because:%v", root, err)
	}
	self.pathExpression = compiled
	return self
}

// RootExpression returns the compiled (RegExp) expression from the rootPath
func (self WebService) RootExpression() *PathExpression {
	return self.pathExpression
}

// AddParameter adds a PathParameter to document parameters used in the root path.
func (self *WebService) Param(parameter *Parameter) *WebService {
	if self.pathParameters == nil {
		self.pathParameters = []*Parameter{}
	}
	self.pathParameters = append(self.pathParameters, parameter)
	return self
}

// PathParameter creates a new Parameter of kind Path for documentation purposes.
func (self *WebService) PathParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true}}
	p.bePath()
	return p
}

// QueryParameter creates a new Parameter of kind Query for documentation purposes.
func (self *WebService) QueryParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: false}}
	p.beQuery()
	return p
}

// BodyParameter creates a new Parameter of kind Body for documentation purposes.
func (self *WebService) BodyParameter(name, description string) *Parameter {
	p := &Parameter{&ParameterData{Name: name, Description: description, Required: true}}
	p.beBody()
	return p
}

// Route creates a new Route using the RouteBuilder and add to the ordered list of Routes.
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.produces, self.consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}

// Method creates a new RouteBuilder and initialize its http method
func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method(httpMethod)
}

// Produces specifies that this WebService can produce one or more MIME types.
func (self *WebService) Produces(contentTypes ...string) *WebService {
	self.produces = contentTypes
	return self
}

// Produces specifies that this WebService can consume one or more MIME types.
func (self *WebService) Consumes(accepts ...string) *WebService {
	self.consumes = accepts
	return self
}

// Routes returns the Routes associated with this WebService
func (self WebService) Routes() []Route {
	return self.routes
}

// RootPath returns the RootPath associated with this WebService. Default "/"
func (self WebService) RootPath() string {
	return self.rootPath
}

// PathParameters return the path parameter names for (shared amoung its Routes)
func (self WebService) PathParameters() []*Parameter {
	return self.pathParameters
}

// Filters returns the list of FilterFunction
func (self WebService) Filters() []FilterFunction {
	return self.filters
}

// Filter adds a filter function to the chain of filters applicable to all its Routes
func (self *WebService) Filter(filter FilterFunction) *WebService {
	self.filters = append(self.filters, filter)
	return self
}

/*
	Convenience methods
*/

// GET is a shortcut for .Method("GET").Path(subPath)
func (self *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("GET").Path(subPath)
}

// POST is a shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("POST").Path(subPath)
}

// PUT is a shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("PUT").Path(subPath)
}

// DELETE is a shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("DELETE").Path(subPath)
}
