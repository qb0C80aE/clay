package model

// TemplateGenerationType is the model class what represents template generation
type TemplateGenerationType struct {
	TemplateGeneration
}

// NewTemplateGenerationType creates a template generation model instance
func NewTemplateGenerationType() *TemplateGenerationType {
	return &TemplateGenerationType{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateGenerationType) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}
