package gorm

import (
	"bytes"
	"log"
	"sort"
	"text/template"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

/*GenGormObjectTableNameOverrideFunc generates a function to override the GORM default table name
See examples over override at http://doc.gorm.io/models.html#conventions*/
func GenGormObjectTableNameOverrideFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormObjTblNameOverrideTemplate := "\nfunc ({{.Name}}) TableName() string {\n\treturn \"{{.SourceTableName}}\"\n}\n"
	tmpl := template.New("gormObjTblNameOverride")
	tmpl, err := tmpl.Parse(gormObjTblNameOverrideTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormCreateFunc generates functions for basic CRUD Create using gorm and writes it to a buffer*/
func GenObjectGormCreateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormCreateFuncTemplate := "\n\nfunc Create{{.Name}} (db *gorm.DB, new{{.Name}} {{.Name}}) []error {\n\treturn db.Create(&new{{.Name}}).GetErrors()\n}"
	gormCreateFuncTemplate += "\n\nfunc GetLast{{.Name}} (db *gorm.DB, new{{.Name}} *{{.Name}}) []error {\n\treturn db.Last(new{{.Name}}).GetErrors()\n}"
	tmpl := template.New("gormCreateFunc")
	tmpl, err := tmpl.Parse(gormCreateFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormReadFunc generates functions for basic CRUD Read/Get using gorm and writes it to a buffer*/
func GenObjectGormReadFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	//gormReadFuncTemplate := "\n\nfunc Get{{.Name}} (db *gorm.DB, query{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) []error {\n\treturn db.Where(&query{{.Name}}).First(result{{.Name}}).GetErrors()\n}"
	gormReadFuncTemplate := `

func Get{{.Name}} (db *gorm.DB, query{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) []error {
	return db.Where(&query{{.Name}}).First(result{{.Name}}).GetErrors()
}`

	tmpl := template.New("gormReadFunc")
	tmpl, err := tmpl.Parse(gormReadFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormReadAllFunc generates functions for basic CRUD Read/Get All using gorm and writes it to a buffer*/
func GenObjectGormReadAllFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormReadFuncTemplate := "\n\nfunc GetAll{{.Name}} (db *gorm.DB, query{{.Name}} {{.Name}}, result{{.Name}} *[]{{.Name}}) []error {\n\treturn db.Where(&query{{.Name}}).Find(result{{.Name}}).GetErrors()\n}"
	tmpl := template.New("gormReadFunc")
	tmpl, err := tmpl.Parse(gormReadFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormUpdateFunc generates functions for basic CRUD Update using gorm and writes it to a buffer*/
func GenObjectGormUpdateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	searchKeyTypes := []string{"p", "PRIMARY KEY", "u", "UNIQUE"}
	keyColumn := ""
	for _, searchKeyType := range searchKeyTypes {
		keyColumn = substancegen.SearchForKeyColumnByKeyType(gqlObjectType, searchKeyType)
		if keyColumn != "" {
			break
		}
	}

	if keyColumn == "" {
		log.Printf("No primary or unique key column found for object %s. Skipping Gorm update func.\nWriting useless update function\n", gqlObjectType.Name)
		gormUpdateFuncTemplate := "\n\nfunc Update{{.Name}} (db *gorm.DB, old{{.Name}} {{.Name}}, new{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) []error {\n\tvar err []error\n\t err=append(err,errors.New('Cannot update a view'))\n\treturn err}"

		var templateData = struct {
			Name string
			Key  string
		}{
			gqlObjectType.Name,
			keyColumn,
		}

		tmpl := template.New("gormUpdateFunc")
		tmpl, err := tmpl.Parse(gormUpdateFuncTemplate)

		if err != nil {
			log.Fatal("Parse: ", err)
			return
		}
		err1 := tmpl.Execute(buff, templateData)
		if err1 != nil {
			log.Fatal("Execute: ", err1)
			return
		}

		return

	}

	var templateData = struct {
		Name string
		Key  string
	}{
		gqlObjectType.Name,
		keyColumn,
	}

	gormUpdateFuncTemplate := "\n\nfunc Update{{.Name}} (db *gorm.DB, old{{.Name}} {{.Name}}, new{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) []error {\n\terr := db.Model(&old{{.Name}}).Updates(new{{.Name}}).GetErrors()\n\terr = append(err,Get{{.Name}}(db, new{{.Name}}, result{{.Name}})...)\n\treturn err\n}"
	tmpl := template.New("gormUpdateFunc")
	tmpl, err := tmpl.Parse(gormUpdateFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, templateData)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormDeleteFunc generates functions for basic CRUD Delete using gorm and writes it to a buffer*/
func GenObjectGormDeleteFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormDeleteFuncTemplate := "\n\nfunc Delete{{.Name}} (db *gorm.DB, old{{.Name}} {{.Name}}) []error {\n\treturn db.Delete(&old{{.Name}}).GetErrors()\n}"
	tmpl := template.New("gormReadFunc")
	tmpl, err := tmpl.Parse(gormDeleteFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormCrud generates functions for basic CRUD operations using gorm and writes it to a buffer*/
func GenObjectGormCrud(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	GenObjectGormCreateFunc(gqlObjectType, buff)

	GenObjectGormReadFunc(gqlObjectType, buff)

	GenObjectGormReadAllFunc(gqlObjectType, buff)

	GenObjectGormUpdateFunc(gqlObjectType, buff)

	GenObjectGormDeleteFunc(gqlObjectType, buff)
}

/*GenObjectsGormCrud processes gqlObjectTypes map in sorted key order and calls GenObjectGormCrud
This is done to produce predictable output*/
func GenObjectsGormCrud(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	keys := make([]string, 0)
	for key := range gqlObjectTypes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		GenObjectGormCrud(gqlObjectTypes[key], buff)
	}
}
