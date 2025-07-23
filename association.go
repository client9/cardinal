package sexpr

import (
	"github.com/client9/sexpr/stdlib"
)

// Re-export association types from stdlib for backward compatibility
type AssociationEntry = stdlib.AssociationEntry
type AssociationValue = stdlib.AssociationValue

// Re-export association functions from stdlib for backward compatibility
var (
	NewAssociationValue     = stdlib.NewAssociationValue
	NewAssociationFromPairs = stdlib.NewAssociationFromPairs
	NewAssociation          = stdlib.NewAssociation
)
