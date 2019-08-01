package graphqlgo

var b = "`"

var JwtUtilities = `

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mitchellh/mapstructure"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte = []byte("thepolyglotdeveloper")


func checkHeaderJwtToken(header http.Header) error {

	if t, ok := header["Authorization"]; ok {
		token := t[0]
		result, validationErr := ValidateJWT(token)

		if validationErr != nil {
			return errors.New("Invalid Jwt token")
		}

		log.Println("Found valid jwt token :", result)
		return nil

	} else {
		return errors.New("Jwt token not found")

	}

}



type UserForToken struct {
	Username string ` + b + `json:"username"` + b + `
	Password string ` + b + `json:"password"` + b + `
}

func ValidateJWT(t string) (interface{}, error) {
	if t == "" || t == "null" {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("There was an error")
		}
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	} else {
		return nil, errors.New("Invalid authorization token")
	}
}


func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func CreateTokenEndpoint(response http.ResponseWriter, request *http.Request) {
	enableCors(&response)
	var user UserForToken
	log.Println(request.Body)
	errDecode := json.NewDecoder(request.Body).Decode(&user)

	if errDecode != nil {
		var tokenString string
		response.Header().Set("content-type", "application/json")
		response.Write([]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `))
		log.Println("User and pass not provided")
		return
	}
	
	if len(user.Password) < 1 || len(user.Username) < 1 {
		var tokenString string
		response.Header().Set("content-type", "application/json")
		response.Write([]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `))
		log.Println("User and pass not provided")
		return
	}

	log.Println(user)
	QueryUserObj := User{}
	QueryUserObj.Username = user.Username
	QueryUserObj.Password = user.Password

	var ResultUserObj User

	err := GetUser(DB, QueryUserObj, &ResultUserObj)

	log.Println("Error :", err)
	if len(err) > 0 {
		log.Println(err)
		var tokenString string
		response.Header().Set("content-type", "application/json")
		response.Write([]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `))
		log.Println("User and pass dosent match")
		return

	}

	log.Println(ResultUserObj)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	tokenString, error := token.SignedString(jwtSecret)
	if error != nil {
		response.Header().Set("content-type", "application/json")
		response.Write([]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `))
		log.Println("Error signin token")
		log.Println(error)
	}
	

	result, validationErr := ValidateJWT(tokenString)
	if validationErr != nil {

                tokenString = ""
				byteToken :=[]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `)
				response.Header().Set("content-type", "application/json")
				response.Write(byteToken)


		log.Println(validationErr)
	} else {
		
			log.Println("Validation ok in request login")
			log.Println(result)
			UpdateUserObj := User{}
			//copier.Copy(ResultUserObj, UpdateUserObj)
	        if ResultUserObj.Token != tokenString {

			UpdateUserObj = ResultUserObj
			UpdateUserObj.Token = tokenString
			var Result User
			log.Println(ResultUserObj, UpdateUserObj)
	
			err = UpdateUser(DB, ResultUserObj, UpdateUserObj, &Result)
			log.Println(Result)
			if len(err) > 0 {
	
				tokenString = ""
				byteToken :=[]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `)
				response.Header().Set("content-type", "application/json")
				response.Write(byteToken)
				log.Println("Error updating token")
				log.Println(error)
			} else {
				byteToken :=[]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `)
				response.Header().Set("content-type", "application/json")
				response.Write(byteToken)
			}
	
			} else {
				byteToken :=[]byte(` + b + `{ "token": "` + b + `+ tokenString +` + b + `"}` + b + `)
				response.Header().Set("content-type", "application/json")
				response.Write(byteToken)
			} // token is different

	}
}

`

/*GraphqlGoExecuteQueryFunc boilerplate string for graphql-go function to execute a graphql query*/
var GraphqlGoExecuteQueryFunc = `
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
`

var graphqlGoMainFunc = `
var DB *gorm.DB

func main() {

	DB, _ = gorm.Open("{{.DbType}}","{{.ConnectionString}}")
	defer DB.Close()


	fmt.Println("Test with Get	:	curl -g 'http://localhost:8080/graphql?query={ {{.SampleQuery}} }'")

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: QueryFields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: MutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// apply middleware
	h := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// some kind of authentication here
			enableCors(&w)

			// do normal graphql
			key := "header"
			innerCtx := context.WithValue(r.Context(), key, r.Header)
			inner.ServeHTTP(w, r.WithContext(innerCtx))
		})
	}(handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}))

	http.Handle("/graphql", h)
	http.HandleFunc("/login", CreateTokenEndpoint)

	fmt.Println("Now server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}
`

var graphqlTypesTemplate = `{{range $key, $value := . }}
var {{.LowerName}}Type = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "{{.Name}}",
		Fields: graphql.Fields{ {{range .Properties}}
			"{{.ScalarNameUpper}}":&graphql.Field{
				Type: {{index .AltScalarType "graphql-go"}},
			},{{end}}
		},
	},
)
{{end}}`

var graphqlGoFieldsQueryTemplate = `
var QueryFields graphql.Fields

func init() {
	QueryFields = make(graphql.Fields,1)
	{{template "graphqlFieldsGet" .}}
	{{template "graphqlFieldsGetAll" .}}
}
`

var graphqlQueryTemplate = `{{range $key, $value := . }}Get{{.Name}} { {{range .Properties}}{{if not .IsObjectType}}{{.ScalarName}},{{end}} {{end}}},{{end}}`

var graphqlGoQueryFieldsGetTemplate = `{{define "graphqlFieldsGet"}}{{range $key, $value := . }}
	QueryFields["Get{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		       header := p.Context.Value("header").(http.Header)
			validationError := checkHeaderJwtToken(header)
			if validationError != nil {

				return nil, validationError

			}

			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName | Title}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}{{$name := .Name}}
			var Result{{$name}}Obj {{.Name}}
			err := Get{{.Name}}(DB,Query{{.Name}}Obj,&Result{{$name}}Obj){{range .Properties}}{{if .IsObjectType}}
			{{.ScalarName}}Obj := {{if .IsList}}[]{{end}}{{.ScalarType}}{}
			err = append(err,DB.Model(&Result{{$name}}Obj).Association("{{.ScalarName}}").Find(&{{.ScalarName}}Obj).Error)
			Result{{$name}}Obj.{{.ScalarName}} = {{if .IsList}}append(Result{{$name}}Obj.{{.ScalarName}}, {{.ScalarName}}Obj...){{else}}{{.ScalarName}}Obj{{end}}{{end}}{{end}}
			if len(err) > 0 {
				return Result{{$name}}Obj, err[len(err)-1]
			}
			return Result{{$name}}Obj, nil
		},
	}
{{end}}{{end}}
`

var graphqlGoFieldsMutationTemplate = `
var MutationFields graphql.Fields

func init() {
	MutationFields = make(graphql.Fields,1)
	{{template "graphqlFieldsCreate" .}}
	{{template "graphqlFieldsDelete" .}}
	{{template "graphqlFieldsUpdate" .}}
}
`

var graphqlGoMutationCreateTemplate = `{{define "graphqlFieldsCreate"}}{{range $key, $value := . }}
	MutationFields["Create{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if and (not .IsObjectType)  (not (eq .ScalarName "id"))}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		       
		       header := p.Context.Value("header").(http.Header)
		       validationError := checkHeaderJwtToken(header)
		       if validationError != nil {

				return nil, validationError

			}

			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if and (not .IsObjectType)  (not (eq .ScalarName "id"))}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName | Title}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}
			err := Create{{.Name}}(DB,Query{{.Name}}Obj)
			var Result{{.Name}}Obj {{.Name}}
			err = append(err,GetLast{{.Name}}(DB,&Result{{.Name}}Obj)...)
			if len(err) > 0 {
				return Result{{.Name}}Obj, err[len(err)-1]
			} 
			return Result{{.Name}}Obj, nil
		},
	}
{{end}}{{end}}
`

var graphqlGoMutationDeleteTemplate = `{{define "graphqlFieldsDelete"}}{{range $key, $value := . }}
	MutationFields["Delete{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		header := p.Context.Value("header").(http.Header)
			validationError := checkHeaderJwtToken(header)
			if validationError != nil {

				return nil, validationError

			}

			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName | Title}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}
			err := Delete{{.Name}}(DB,Query{{.Name}}Obj)
			var Result{{.Name}}Obj {{.Name}}
			//err = append(err,Get{{.Name}}(DB,Query{{.Name}}Obj,&Result{{.Name}}Obj)...)
			if len(err) > 0 {
				return Result{{.Name}}Obj, err[len(err)-1]
			}
			return Result{{.Name}}Obj, nil
		},
	}
{{end}}{{end}}
`

var graphqlGoMutationUpdateTemplate = `{{define "graphqlFieldsUpdate"}}{{range $key, $value := . }}{{$primaryKeyCol := getPkeyCol . "p"}}{{$primaryKeyColAlt := getPkeyCol . "PRIMARY KEY"}}
	MutationFields["Update{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		header := p.Context.Value("header").(http.Header)
			validationError := checkHeaderJwtToken(header)
			if validationError != nil {

				return nil, validationError

			}

			Old{{.Name}}Obj := {{.Name}}{}
			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				{{if or (eq $primaryKeyCol .ScalarName) (eq $primaryKeyColAlt .ScalarName)}}Old{{$value.Name}}Obj.{{.ScalarName}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}{{end}}
				Query{{$value.Name}}Obj.{{.ScalarName | Title}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}
			var Result{{.Name}}Obj {{.Name}}
			err := Update{{.Name}}(DB,Old{{.Name}}Obj,Query{{.Name}}Obj,&Result{{.Name}}Obj)
			if len(err) > 0 {
				return Result{{.Name}}Obj, err[len(err)-1]
			}
			return Result{{.Name}}Obj, nil
		},
	}
{{end}}{{end}}
`

var graphqlGoQueryFieldsGetAllTemplate = `{{define "graphqlFieldsGetAll"}}{{range $key, $value := . }}
	QueryFields["GetAll{{.Name}}"] = &graphql.Field{
		Type: graphql.NewList({{.LowerName}}Type),
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		
		//Jwt token checker
		header := p.Context.Value("header").(http.Header)
			validationError := checkHeaderJwtToken(header)
			if validationError != nil {

				return nil, validationError

			}

			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName | Title}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}{{$name := .Name}}
			var Result{{$name}}Obj []{{.Name}}
			err := GetAll{{.Name}}(DB,Query{{.Name}}Obj,&Result{{$name}}Obj)
			if len(err) > 0 {
				return Result{{$name}}Obj, err[len(err)-1]
			}
			return Result{{$name}}Obj, nil
		},
	}
{{end}}{{end}}
`
