package firemodel

import "fmt"

func (s *Schema) ParentModel(model *SchemaModel) *SchemaModel {
	parts := model.FirestorePath.CollectionParts
	if len(parts) < 2 {
		return nil
	}

	parentCollectionName := parts[len(parts)-2].CollectionName
	for _, maybeParentModel := range s.Models {
		maybeParentFirestorePathCollectionParts := maybeParentModel.FirestorePath.CollectionParts
		if parent := maybeParentFirestorePathCollectionParts[len(maybeParentFirestorePathCollectionParts)-1]; parent.CollectionName == parentCollectionName {
			return maybeParentModel
		}
	}
	panic(fmt.Sprintf("no parent collection exists for model %s", model.FirestorePath))
}

func (s *Schema) RootModels() []*SchemaModel {
	var ret []*SchemaModel
	for _, model := range s.Models {
		if len(model.FirestorePath.CollectionParts) == 1 {
			ret = append(ret, model)
		}
	}
	return ret
}

func (s *Schema) DirectSubcollectionsOfModel(model *SchemaModel) []*SchemaModel {
	var ret []*SchemaModel
	for _, schemaModel := range s.Models {
		// check if this schema model is a direct child of this model
		if len(schemaModel.FirestorePath.CollectionParts) != len(model.FirestorePath.CollectionParts)+1 {
			// Not a direct descendant.
			goto continueNextSchemaModel
		}
		for modelPathTemplatePartIdx, modelPathTemplatePart := range model.FirestorePath.CollectionParts {
			if modelPathTemplatePart.CollectionName != schemaModel.FirestorePath.CollectionParts[modelPathTemplatePartIdx].CollectionName {
				// Collections don't match up.
				goto continueNextSchemaModel
			}
		}
		// match!
		ret = append(ret, schemaModel)
	continueNextSchemaModel:
	}
	return ret
}
