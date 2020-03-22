package utils

import (
  "fmt"
)

type Envi struct {
  Name *string
  initialTime *int32
}

func (this *Envi) SetName (name string) error {
  if this.Name != nil {
    return fmt.Errorf("name already defined "+ *this.Name)
  }
  this.Name = &name
  return nil
}
