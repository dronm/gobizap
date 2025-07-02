package md

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//File contains metadata structure

type MdBoolean bool

func (b *MdBoolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := ""
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}
	*b = (strings.ToUpper(s) == "TRUE")
	return nil
}

func (b *MdBoolean) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := ""
	if bool(*b) {
		v = "TRUE"
	} else {
		v = "FALSE"
	}
	return e.EncodeElement(v, start)
}

type MdDateTime time.Time

func (dt *MdDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := ""
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}
	dat, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*dt = MdDateTime(dat)
	return nil
}

type Metadata struct {
	XMLName       xml.Name       `xml:"metadata"`
	Debug         MdBoolean      `xml:"debug,attr"`
	Owner         string         `xml:"owner,attr"`
	DataSchema    string         `xml:"dataSchema,attr"`
	Versions      Versions       `xml:"versions"`
	GlobalFilters GlobalFilters  `xml:"globalFilters"`
	Enums         []Enum         `xml:"enums>enum"`
	Models        Models         `xml:"models"`
	Constants     Constants      `xml:"constants"`
	Controllers   Controllers    `xml:"controllers"`
	Permissions   []Permission   `xml:"permissions>permission"`
	Views         []View         `xml:"views>view"`
	ServerTempl   ServerTemplate `xml:"serverTemplates>serverTemplate"`
	JSTemplates   JSTemplate     `xml:"jsTemplates>jsTemplate"`
	Packages      []Package      `xml:"packages>package"`
	SQLScripts    []SQLScript    `xml:"sqlScripts>sqlScript"`
	JSScripts     []JSScript     `xml:"jsScripts>jsScript"`
	CSSScripts    []CSSScript    `xml:"cssScripts>cssScript"`
}

type EnumLang struct {
	XMLName xml.Name
	Descr   string `xml:"descr,attr"`
}

type EnumLangs map[string]string //key - lang ID, value - descr

func (m *EnumLangs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if *m == nil {
		*m = EnumLangs{}
	}

	e := EnumLang{}
	err := d.DecodeElement(&e, &start)
	if err != nil {
		return err
	}
	(*m)[e.XMLName.Local] = e.Descr

	return nil
}

func (m *EnumLangs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var v strings.Builder
	for lang_id, lang_descr := range *m {
		v.WriteString(fmt.Sprintf(`<%s descr="%s"></%s>`, lang_id, lang_descr, lang_id))
	}
	return e.EncodeElement(v.String(), start)
}

type Versions struct {
	Version Version `xml:"version"`
}

type Version struct {
	DateOpen  string `xml:"dateOpen,attr"`
	Value     string `xml:",chardata"`
	LastBuild string `xml:"lastBuild"`
}

type GlobalFilters struct {
	Field Field `xml:"field"`
}

type Field struct {
	ID            string    `xml:"id,attr,omitempty"`
	Cmd           *string   `xml:"cmd,attr,omitempty" json:"cmd,omitempty"`
	OldID         string    `xml:"oldId,attr,omitempty"`
	SortDirect    string    `xml:"sortDirect,attr,omitempty"`
	DataType      string    `xml:"dataType,attr,omitempty"`
	OldDataType   string    `xml:"oldDataType,attr,omitempty"`
	EnumID        string    `xml:"enumId,attr,omitempty"`
	Alias         string    `xml:"alias,attr,omitempty"`
	SysCol        MdBoolean `xml:"sysCol,attr,omitempty"`
	RefTable      string    `xml:"refTable,attr,omitempty"`
	RefField      string    `xml:"refField,attr,omitempty"`
	Length        string    `xml:"length,attr,omitempty"`
	Precision     string    `xml:"precision,attr,omitempty"`
	Required      MdBoolean `xml:"required,attr,omitempty"`
	AutoInc       MdBoolean `xml:"autoInc,attr,omitempty"`
	PrimaryKey    MdBoolean `xml:"primaryKey,attr,omitempty"`
	NoValueOnCopy MdBoolean `xml:"noValueOnCopy,attr,omitempty"`
}

type Enum struct {
	ID     string  `xml:"id,attr,omitempty" json:"id"`
	Cmd    *string `xml:"cmd,attr,omitempty" json:"cmd,omitempty"`
	Descr  string  `xml:"descr,attr,omitempty" json:"descr"`
	Values []Value `xml:"value,omitempty" json:"values,omitempty"`
}

type Value struct {
	ID               string    `xml:"id,attr"`
	LangDescriptions EnumLangs `xml:",any"`
}

type Models struct {
	Model []Model `xml:"model"`
}

type Model struct {
	ID              string         `xml:"id,attr,omitempty"`
	Cmd             *string        `xml:"cmd,attr,omitempty"`
	Parent          string         `xml:"parent,attr,omitempty"`
	DataSchema      string         `xml:"dataSchema,attr,omitempty"`
	DataTable       string         `xml:"dataTable,attr,omitempty"`
	BaseModelID     string         `xml:"baseModelId,attr,omitempty"`
	Virtual         *MdBoolean     `xml:"virtual,attr,omitempty"`
	DefaultOrder    []DefaultOrder `xml:"defaultOrder,omitempty"`
	Field           []Field        `xml:"field,omitempty"`
	Index           []Index        `xml:"index,omitempty"`
	AggFunctions    []AggFunction  `xml:"aggFunction,omitempty"`
	LimitCount      string         `xml:"limitCount,attr,omitempty"`
	DocPerPageCount string         `xml:"docPerPageCount,attr,omitempty"`
}

type DefaultOrder struct {
	Field Field `xml:"field"`
}

type AggFunction struct {
	Alias string `xml:"alias"`
	Expr  string `xml:"expr"`
}

type IndexField struct {
	ID    string `xml:"id,attr,omitempty" json:"id"`
	Order string `xml:"order,attr,omitempty" json:"order,omitempty"`
	Nulls string `xml:"nulls,attr,omitempty" json:"nulls"`
}

type Index struct {
	ID     string       `xml:"id,attr"`
	Cmd    *string      `xml:"cmd,attr,omitempty"`
	Unique MdBoolean    `xml:"unique,attr,omitempty"`
	Expr   string       `xml:"expression,omitempty"`
	Type   string       `xml:"type,omitempty"`
	Field  []IndexField `xml:"field,omitempty"`
}

type Constants struct {
	Constant []Constant `xml:"constant"`
}

type Constant struct {
	ID           string  `xml:"id,attr,omitempty"`
	Cmd          *string `xml:"cmd,attr,omitempty"`
	Name         string  `xml:"name,attr,omitempty"`
	Descr        string  `xml:"descr,attr,omitempty"`
	DataType     string  `xml:"dataType,attr,omitempty"`
	DefaultValue string  `xml:"defaultValue,attr,omitempty"`
	Autoload     string  `xml:"autoload,attr,omitempty"`
}

type Controllers struct {
	Controller []Controller `xml:"controller"`
}

type Controller struct {
	ID           string         `xml:"id,attr,omitempty"`
	ParentID     string         `xml:"parentId,attr,omitempty"`
	Cmd          *string        `xml:"cmd,attr,omitempty"`
	Client       *MdBoolean     `xml:"client,attr,omitempty"`
	Server       *MdBoolean     `xml:"server,attr,omitempty"`
	PublicMethod []PublicMethod `xml:"publicMethod,omitempty"`
}

type PublicMethod struct {
	ID    string `xml:"id,attr"`
	Field Field  `xml:"field"`
}

type Permission struct {
	Type         string  `xml:"type,attr,omitempty"`
	Cmd          *string `xml:"cmd,attr,omitempty"`
	ControllerID string  `xml:"controllerId,attr,omitempty"`
	MethodID     string  `xml:"methodId,attr,omitempty"`
	RoleID       string  `xml:"roleId,attr,omitempty"`
}

type View struct {
	ID      string `xml:"id,attr,omitempty"`
	C       string `xml:"c,attr,omitempty"`
	T       string `xml:"t,attr,omitempty"`
	F       string `xml:"f,attr,omitempty"`
	Section string `xml:"section,attr,omitempty"`
	Descr   string `xml:"descr,attr,omitempty"`
}

type ServerTemplate struct {
	ID    string `xml:"id,attr,omitempty"`
	Class string `xml:"class,attr,omitempty"`
}

type JSTemplate struct {
	ID   string `xml:"id,attr,omitempty"`
	File string `xml:"file,attr,omitempty"`
}

type Package struct {
}

type SQLScript struct {
	File string `xml:"file,attr"`
}

type JSScript struct {
	File       string    `xml:"file,attr,omitempty"`
	Compressed MdBoolean `xml:"compressed,attr,omitempty"`
	Standalone MdBoolean `xml:"standalone,attr,omitempty"`
	Cmd        *string   `xml:"cmd,attr,omitempty"`
}

func (s JSScript) GetFile() string {
	return s.File
}
func (s JSScript) GetCompressed() bool {
	return bool(s.Compressed)
}
func (s JSScript) GetStandalone() bool {
	return bool(s.Standalone)
}

type CSSScript struct {
	File       string    `xml:"file,attr,omitempty"`
	Compressed MdBoolean `xml:"compressed,attr,omitempty"`
	Standalone MdBoolean `xml:"standalone,attr,omitempty"`
	Cmd        *string   `xml:"cmd,attr,omitempty"`
}

func (s CSSScript) GetFile() string {
	return s.File
}
func (s CSSScript) GetCompressed() bool {
	return bool(s.Compressed)
}
func (s CSSScript) GetStandalone() bool {
	return bool(s.Standalone)
}

func NewMetadata() (*Metadata, error) {
	projDir, err := GetProjectDir()
	if err != nil {
		return nil, fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	mdFile := filepath.Join(projDir, BUILD_DIR, MD_FILE_NAME)
	md_file_cont, err := os.ReadFile(mdFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile() failed: %v on file: %s", err, mdFile)
	}
	md := Metadata{}
	if err := xml.Unmarshal(md_file_cont, &md); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal() failed: %v on file: %s", err, mdFile)
	}

	return &md, nil
}

func (md *Metadata) AddNewJsScripts(process string, jsScripts []string) bool {
	modified := false
	for _, scr := range jsScripts {
		//check if not exists
		rel_scr := filepath.Join(BUILD_DIR, JS_DIR, scr)
		scr_exitst := false
		for _, md_scr := range md.JSScripts {
			if md_scr.File == rel_scr {
				scr_exitst = true
				break
			}
		}
		if !scr_exitst {
			//add
			cmd := "add"
			md.JSScripts = append(md.JSScripts, JSScript{File: rel_scr, Cmd: &cmd})
			// LogInfo(process, "added new java script '%s'", rel_scr)
			if !modified {
				modified = true
			}
		}
	}
	return modified
}

// WriteMd saves metadata file
func (md *Metadata) WriteMd(projDir string) error {
	md_file_name := filepath.Join(BUILD_DIR, MD_FILE_NAME)
	output, err := xml.MarshalIndent(md, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(md_file_name, output, FILE_PERMISSION); err != nil {
		return err
	}
	return nil
}
