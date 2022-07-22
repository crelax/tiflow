// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// delete binlog operator
	// (DELETE /binlog/tasks/{task-name})
	DMAPIDeleteBinlogOperator(c *gin.Context, taskName string, params DMAPIDeleteBinlogOperatorParams)
	// get binlog operator
	// (GET /binlog/tasks/{task-name})
	DMAPIGetBinlogOperator(c *gin.Context, taskName string, params DMAPIGetBinlogOperatorParams)
	// set binlog operator
	// (POST /binlog/tasks/{task-name})
	DMAPISetBinlogOperator(c *gin.Context, taskName string)
	// get job config
	// (GET /config)
	DMAPIGetJobConfig(c *gin.Context)
	// update job config
	// (PUT /config)
	DMAPIUpdateJobConfig(c *gin.Context)
	// get schema
	// (GET /schema/tasks/{task-name})
	DMAPIGetSchema(c *gin.Context, taskName string, params DMAPIGetSchemaParams)
	// set schema
	// (PUT /schema/tasks/{task-name})
	DMAPISetSchema(c *gin.Context, taskName string)
	// get the current status of the job
	// (GET /status)
	DMAPIGetJobStatus(c *gin.Context, params DMAPIGetJobStatusParams)
	// operate the stage of the job
	// (PUT /status)
	DMAPIOperateJob(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(c *gin.Context)

// DMAPIDeleteBinlogOperator operation middleware
func (siw *ServerInterfaceWrapper) DMAPIDeleteBinlogOperator(c *gin.Context) {
	var err error

	// ------------- Path parameter "task-name" -------------
	var taskName string

	err = runtime.BindStyledParameter("simple", false, "task-name", c.Param("task-name"), &taskName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter task-name: %s", err)})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DMAPIDeleteBinlogOperatorParams

	// ------------- Optional query parameter "binlog_pos" -------------
	if paramValue := c.Query("binlog_pos"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "binlog_pos", c.Request.URL.Query(), &params.BinlogPos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter binlog_pos: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIDeleteBinlogOperator(c, taskName, params)
}

// DMAPIGetBinlogOperator operation middleware
func (siw *ServerInterfaceWrapper) DMAPIGetBinlogOperator(c *gin.Context) {
	var err error

	// ------------- Path parameter "task-name" -------------
	var taskName string

	err = runtime.BindStyledParameter("simple", false, "task-name", c.Param("task-name"), &taskName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter task-name: %s", err)})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DMAPIGetBinlogOperatorParams

	// ------------- Optional query parameter "binlog_pos" -------------
	if paramValue := c.Query("binlog_pos"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "binlog_pos", c.Request.URL.Query(), &params.BinlogPos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter binlog_pos: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIGetBinlogOperator(c, taskName, params)
}

// DMAPISetBinlogOperator operation middleware
func (siw *ServerInterfaceWrapper) DMAPISetBinlogOperator(c *gin.Context) {
	var err error

	// ------------- Path parameter "task-name" -------------
	var taskName string

	err = runtime.BindStyledParameter("simple", false, "task-name", c.Param("task-name"), &taskName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter task-name: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPISetBinlogOperator(c, taskName)
}

// DMAPIGetJobConfig operation middleware
func (siw *ServerInterfaceWrapper) DMAPIGetJobConfig(c *gin.Context) {
	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIGetJobConfig(c)
}

// DMAPIUpdateJobConfig operation middleware
func (siw *ServerInterfaceWrapper) DMAPIUpdateJobConfig(c *gin.Context) {
	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIUpdateJobConfig(c)
}

// DMAPIGetSchema operation middleware
func (siw *ServerInterfaceWrapper) DMAPIGetSchema(c *gin.Context) {
	var err error

	// ------------- Path parameter "task-name" -------------
	var taskName string

	err = runtime.BindStyledParameter("simple", false, "task-name", c.Param("task-name"), &taskName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter task-name: %s", err)})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DMAPIGetSchemaParams

	// ------------- Optional query parameter "database" -------------
	if paramValue := c.Query("database"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "database", c.Request.URL.Query(), &params.Database)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter database: %s", err)})
		return
	}

	// ------------- Optional query parameter "table" -------------
	if paramValue := c.Query("table"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "table", c.Request.URL.Query(), &params.Table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter table: %s", err)})
		return
	}

	// ------------- Optional query parameter "target" -------------
	if paramValue := c.Query("target"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "target", c.Request.URL.Query(), &params.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter target: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIGetSchema(c, taskName, params)
}

// DMAPISetSchema operation middleware
func (siw *ServerInterfaceWrapper) DMAPISetSchema(c *gin.Context) {
	var err error

	// ------------- Path parameter "task-name" -------------
	var taskName string

	err = runtime.BindStyledParameter("simple", false, "task-name", c.Param("task-name"), &taskName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter task-name: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPISetSchema(c, taskName)
}

// DMAPIGetJobStatus operation middleware
func (siw *ServerInterfaceWrapper) DMAPIGetJobStatus(c *gin.Context) {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params DMAPIGetJobStatusParams

	// ------------- Optional query parameter "tasks" -------------
	if paramValue := c.Query("tasks"); paramValue != "" {
	}

	err = runtime.BindQueryParameter("form", true, false, "tasks", c.Request.URL.Query(), &params.Tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter tasks: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIGetJobStatus(c, params)
}

// DMAPIOperateJob operation middleware
func (siw *ServerInterfaceWrapper) DMAPIOperateJob(c *gin.Context) {
	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DMAPIOperateJob(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.DELETE(options.BaseURL+"/binlog/tasks/:task-name", wrapper.DMAPIDeleteBinlogOperator)

	router.GET(options.BaseURL+"/binlog/tasks/:task-name", wrapper.DMAPIGetBinlogOperator)

	router.POST(options.BaseURL+"/binlog/tasks/:task-name", wrapper.DMAPISetBinlogOperator)

	router.GET(options.BaseURL+"/config", wrapper.DMAPIGetJobConfig)

	router.PUT(options.BaseURL+"/config", wrapper.DMAPIUpdateJobConfig)

	router.GET(options.BaseURL+"/schema/tasks/:task-name", wrapper.DMAPIGetSchema)

	router.PUT(options.BaseURL+"/schema/tasks/:task-name", wrapper.DMAPISetSchema)

	router.GET(options.BaseURL+"/status", wrapper.DMAPIGetJobStatus)

	router.PUT(options.BaseURL+"/status", wrapper.DMAPIOperateJob)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{
	"H4sIAAAAAAAC/+RYS4/bNhD+K8S0RyU22px0a5qiSIGgRRdFD4ERUNJIS4ciaT4OhqH/XpCUZMumtDZQ",
	"J4vdS+Ilh9+8vnnYByhlq6RAYQ3kBzDlI7Y0fPxNa6n/ZfbxExpDG/RnFZpSM2WZFJCDVKip/0zQy0IG",
	"SvszyzAglLJKvAqyJNxlYPcKIQcmLDaoocugQksZD++ZxTZ86KWM1Uw0Xqg/oFrTvf+7nTMxKhuus3Ok",
	"LgONO8c0VpB/ht6oQXwzystii6X1mv4MPuMfsvgbdw6N9Tqnbkvl/0XhWo+p0bjWoyrqzCnmiTfUfL3J",
	"4TOzpUqa+oD2PRNcNtFmqWctLoLYFyXTyqcOma9MgTdAcVp6x5gIGlOemR2/p2MPga2zblXU0oIaTCqu",
	"tWy/GOl02dOmpo5byGvKDY5KCyk5UjE+sFQ3aK97YHY87TIteMqmM+dH64cXETEVkH9UFTn5qxQ1a2YD",
	"Uobrp1X3cpeqvCATtYRcOM49M1BQxSCHn9+u364Dy+1j0LWKpFoFcq8O/r83grbYxeBxtBiLpe8hHyvI",
	"4cOnX/76+CFcTqkbkDVt0aI2kH8+r3MPTzx84CPkwQ7IIBzF6zf99dFRqx1mfctLBuUQsXYO9f4IdlIt",
	"S683XpVRUpgY/J/W7y77k3Flicb4JL5br/scWRQheVQpzsoQndXWePnDib4fNdaQww+rYwdf9e17ddG7",
	"Q+amqmvKOFYh98a1LdV7yPvEkOgikcfoW9r4uIdIwqbLoC+DRP5+P+87l8l7tim6LQU9XrhIRPi5JbdB",
	"e01mlTRzqX14OrXfoC43URiNfS+r/f8W19l5mYhvr58U3oDuddHIXEOjLoPVcdwsNotxbsF947iVRYxh",
	"og9HU4mQltTSiVTlbGVBysHOwVuPGWrGzTl4NpvhPtyd2QC66WT3hTVD1udNORfcW0qB51tUk944FhkY",
	"98jvtWNMtQxr30TT2Xg7WQ1vQA5b5BLsuGbegun3YTK8TKOGlTkBO+7Kr3EQm4Fzl/PXLYzf70jVe4/d",
	"6be561vXS56zSZaYod9Zap25Zr4+RMknWNNwWVDO98QJtnNIfJsh8SvyctvwJp0y5+qv+6+w7u0jktJp",
	"jcKSmEAi63DqB9ktm8Xxp6g7LRWXv3W9lH0iRhND2I2lDS7loOvGk0NipRzrIO6XqVY8rRXoNt1/AQAA",
	"//9aG6NCeBUAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
