package fields

//***** Metadata text field:strings/texts ******************
type FieldBool struct {
	Field
}
func (f *FieldBool) GetDataType() FieldDataType {
	return FIELD_TYPE_BOOL
}

